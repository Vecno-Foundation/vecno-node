// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package ethashb3 implements the ethashb3 proof-of-work consensus engine.
package ethashb3

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/edsrzf/mmap-go"
	"github.com/ethereum/go-ethereum/common/hexutil"
	lrupkg "github.com/ethereum/go-ethereum/common/lru"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/rpc"
)

var ErrInvalidDumpMagic = errors.New("invalid dump magic")

var (
	// two256 is a big integer representing 2^256
	two256 = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0))

	// sharedEthash is a full instance that can be shared between multiple users.
	sharedEthash *EthashB3

	// algorithmRevision is the data structure version used for file naming.
	algorithmRevision = 23

	// dumpMagic is a dataset dump header to sanity check a data dump.
	dumpMagic = []uint32{0xbaddcafe, 0xfee1dead}
)

func init() {
	sharedConfig := Config{
		PowMode:       ModeNormal,
		CachesInMem:   3,
		DatasetsInMem: 1,
	}
	sharedEthash = New(sharedConfig, nil, false)
}

// isLittleEndian returns whether the local system is running in little or big
// endian byte order.
func isLittleEndian() bool {
	n := uint32(0x01020304)
	return *(*byte)(unsafe.Pointer(&n)) == 0x04
}

// uint32Array2ByteArray returns the bytes represented by uint32 array c
// nolint:unused
func uint32Array2ByteArray(c []uint32) []byte {
	buf := make([]byte, len(c)*4)
	if isLittleEndian() {
		for i, v := range c {
			binary.LittleEndian.PutUint32(buf[i*4:], v)
		}
	} else {
		for i, v := range c {
			binary.BigEndian.PutUint32(buf[i*4:], v)
		}
	}
	return buf
}

// bytes2Keccak256 returns the keccak256 hash as a hex string (0x prefixed)
// for a given uint32 array (cache/dataset)
// nolint:unused
func uint32Array2Keccak256(data []uint32) string {
	// convert to bytes
	bytes := uint32Array2ByteArray(data)
	// hash with keccak256
	digest := crypto.Keccak256(bytes)
	// return hex string
	return hexutil.Encode(digest)
}

// memoryMap tries to memory map a file of uint32s for read only access.
func memoryMap(path string, lock bool) (*os.File, mmap.MMap, []uint32, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, nil, nil, err
	}
	mem, buffer, err := memoryMapFile(file, false)
	if err != nil {
		file.Close()
		return nil, nil, nil, err
	}
	for i, magic := range dumpMagic {
		if buffer[i] != magic {
			mem.Unmap()
			file.Close()
			return nil, nil, nil, ErrInvalidDumpMagic
		}
	}
	if lock {
		if err := mem.Lock(); err != nil {
			mem.Unmap()
			file.Close()
			return nil, nil, nil, err
		}
	}
	return file, mem, buffer[len(dumpMagic):], err
}

// memoryMapFile tries to memory map an already opened file descriptor.
func memoryMapFile(file *os.File, write bool) (mmap.MMap, []uint32, error) {
	// Try to memory map the file
	flag := mmap.RDONLY
	if write {
		flag = mmap.RDWR
	}
	mem, err := mmap.Map(file, flag, 0)
	if err != nil {
		return nil, nil, err
	}
	// The file is now memory-mapped. Create a []uint32 view of the file.
	var view []uint32
	header := (*reflect.SliceHeader)(unsafe.Pointer(&view))
	header.Data = (*reflect.SliceHeader)(unsafe.Pointer(&mem)).Data
	header.Cap = len(mem) / 4
	header.Len = header.Cap
	return mem, view, nil
}

// memoryMapAndGenerate tries to memory map a temporary file of uint32s for write
// access, fill it with the data from a generator and then move it into the final
// path requested.
func memoryMapAndGenerate(path string, size uint64, lock bool, generator func(buffer []uint32)) (*os.File, mmap.MMap, []uint32, error) {
	// Ensure the data folder exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, nil, nil, err
	}
	// Create a huge temporary empty file to fill with data
	temp := path + "." + strconv.Itoa(rand.Int())

	dump, err := os.Create(temp)
	if err != nil {
		return nil, nil, nil, err
	}
	if err = ensureSize(dump, int64(len(dumpMagic))*4+int64(size)); err != nil {
		dump.Close()
		os.Remove(temp)
		return nil, nil, nil, err
	}
	// Memory map the file for writing and fill it with the generator
	mem, buffer, err := memoryMapFile(dump, true)
	if err != nil {
		dump.Close()
		os.Remove(temp)
		return nil, nil, nil, err
	}
	copy(buffer, dumpMagic)

	data := buffer[len(dumpMagic):]
	generator(data)

	if err := mem.Unmap(); err != nil {
		return nil, nil, nil, err
	}
	if err := dump.Close(); err != nil {
		return nil, nil, nil, err
	}
	if err := os.Rename(temp, path); err != nil {
		return nil, nil, nil, err
	}
	return memoryMap(path, lock)
}

type cacheOrDataset interface {
	*cache | *dataset
}

// lru tracks caches or datasets by their last use time, keeping at most N of them.
type lru[T cacheOrDataset] struct {
	what string
	new  func(epoch uint64, epochLength uint64) T
	mu   sync.Mutex
	// Items are kept in a LRU cache, but there is a special case:
	// We always keep an item for (highest seen epoch) + 1 as the 'future item'.
	cache      lrupkg.BasicLRU[uint64, T]
	future     uint64
	futureItem T
}

// newlru create a new least-recently-used cache for either the verification caches
// or the mining datasets.
func newlru[T cacheOrDataset](maxItems int, new func(epoch uint64, epochLength uint64) T) *lru[T] {
	var what string
	switch any(T(nil)).(type) {
	case *cache:
		what = "cache"
	case *dataset:
		what = "dataset"
	default:
		panic("unknown type")
	}
	return &lru[T]{
		what:  what,
		new:   new,
		cache: lrupkg.NewBasicLRU[uint64, T](maxItems),
	}
}

// get retrieves or creates an item for the given epoch. The first return value is always
// non-nil. The second return value is non-nil if lru thinks that an item will be useful in
// the near future.
func (lru *lru[T]) get(epoch uint64, epochLength uint64) (item, future T) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	// Use the sum of epoch and epochLength as the cache key.
	// This is not perfectly safe, but it's good enough (at least for the first 30000 epochs, or the first 427 years).
	cacheKey := epochLength + epoch

	// Get or create the item for the requested epoch.
	item, ok := lru.cache.Get(cacheKey)
	if !ok {
		if lru.future > 0 && lru.future == epoch {
			item = lru.futureItem
		} else {
			log.Trace("Requiring new ethashb3 "+lru.what, "epoch", epoch)
			item = lru.new(epoch, epochLength)
		}
		lru.cache.Add(cacheKey, item)
	}

	// Ensure pre-generation handles ecip-1099 changeover correctly
	var nextEpoch = epoch + 1
	var nextEpochLength = epochLength

	// Update the 'future item' if epoch is larger than previously seen.
	// Last conditional clause ('lru.future > nextEpoch') handles the ECIP1099 case where
	// the next epoch is expected to be LESSER THAN that of the previous state's future epoch number.
	if epoch < maxEpoch-1 && lru.future != nextEpoch {
		log.Trace("Requiring new future ethashb3 "+lru.what, "epoch", nextEpoch)
		future = lru.new(nextEpoch, nextEpochLength)
		lru.future = nextEpoch
		lru.futureItem = future
	}
	return item, future
}

// cache wraps an ethashb3 cache with some metadata to allow easier concurrent use.
type cache struct {
	epoch       uint64    // Epoch for which this cache is relevant
	epochLength uint64    // Epoch length (ECIP-1099)
	dump        *os.File  // File descriptor of the memory mapped cache
	mmap        mmap.MMap // Memory map itself to unmap before releasing
	cache       []uint32  // The actual cache data content (may be memory mapped)
	once        sync.Once // Ensures the cache is generated only once
}

// newCache creates a new ethashb3 verification cache.
func newCache(epoch uint64, epochLength uint64) *cache {
	return &cache{epoch: epoch, epochLength: epochLength}
}

// generate ensures that the cache content is generated before use.
func (c *cache) generate(dir string, limit int, lock bool, test bool) {
	c.once.Do(func() {
		size := cacheSize(c.epoch)
		seed := seedHash(c.epoch, c.epochLength)
		if test {
			size = 1024
		}
		// If we don't store anything on disk, generate and return.
		if dir == "" {
			c.cache = make([]uint32, size/4)
			generateCache(c.cache, c.epoch, c.epochLength, seed)
			return
		}
		// Disk storage is needed, this will get fancy
		var endian string
		if !isLittleEndian() {
			endian = ".be"
		}
		// The file path naming scheme was changed to include epoch values in the filename,
		// which enables a filepath glob with scan to identify out-of-bounds caches and remove them.
		// The legacy path declaration is provided below as a comment for reference.
		//
		// path := filepath.Join(dir, fmt.Sprintf("cache-R%d-%x%s", algorithmRevision, seed[:8], endian))                 // LEGACY
		path := filepath.Join(dir, fmt.Sprintf("cache-R%d-%d-%x%s", algorithmRevision, c.epoch, seed[:8], endian)) // CURRENT
		logger := log.New("epoch", c.epoch, "epochLength", c.epochLength)

		// We're about to mmap the file, ensure that the mapping is cleaned up when the
		// cache becomes unused.
		runtime.SetFinalizer(c, (*cache).finalizer)

		// Try to load the file from disk and memory map it
		var err error
		c.dump, c.mmap, c.cache, err = memoryMap(path, lock)
		if err == nil {
			logger.Debug("Loaded old ethashb3 cache from disk")
			return
		}
		logger.Debug("Failed to load old ethashb3 cache", "err", err)

		// No usable previous cache available, create a new cache file to fill
		c.dump, c.mmap, c.cache, err = memoryMapAndGenerate(path, size, lock, func(buffer []uint32) { generateCache(buffer, c.epoch, c.epochLength, seed) })
		if err != nil {
			logger.Error("Failed to generate mapped ethashb3 cache", "err", err)

			c.cache = make([]uint32, size/4)
			generateCache(c.cache, c.epoch, c.epochLength, seed)
		}

		// Iterate over all cache file instances, deleting any out of bounds (where epoch is below lower limit, or above upper limit).
		matches, _ := filepath.Glob(filepath.Join(dir, fmt.Sprintf("cache-R%d*", algorithmRevision)))
		for _, file := range matches {
			var ar int   // algorithm revision
			var e uint64 // epoch
			var s string // seed
			if _, err := fmt.Sscanf(filepath.Base(file), "cache-R%d-%d-%s"+endian, &ar, &e, &s); err != nil {
				// There is an unrecognized file in this directory.
				// See if the name matches the expected pattern of the legacy naming scheme.
				if _, err := fmt.Sscanf(filepath.Base(file), "cache-R%d-%s"+endian, &ar, &s); err == nil {
					// This file matches the previous generation naming pattern (sans epoch).
					if err := os.Remove(file); err != nil {
						logger.Error("Failed to remove legacy ethashb3 cache file", "file", file, "err", err)
					} else {
						logger.Warn("Deleted legacy ethashb3 cache file", "path", file)
					}
				}
				// Else the file is unrecognized (unknown name format), leave it alone.
				continue
			}
			if e <= c.epoch-uint64(limit) || e > c.epoch+2 {
				if err := os.Remove(file); err == nil {
					logger.Debug("Deleted ethashb3 cache file", "target.epoch", e, "file", file)
				} else {
					logger.Error("Failed to delete ethashb3 cache file", "target.epoch", e, "file", file, "err", err)
				}
			}
		}
	})
}

// finalizer unmaps the memory and closes the file.
func (c *cache) finalizer() {
	if c.mmap != nil {
		c.mmap.Unmap()
		c.dump.Close()
		c.mmap, c.dump = nil, nil
	}
}

// dataset wraps an ethashb3 dataset with some metadata to allow easier concurrent use.
type dataset struct {
	epoch       uint64      // Epoch for which this cache is relevant
	epochLength uint64      // Epoch length (ECIP-1099)
	dump        *os.File    // File descriptor of the memory mapped cache
	mmap        mmap.MMap   // Memory map itself to unmap before releasing
	dataset     []uint32    // The actual cache data content
	once        sync.Once   // Ensures the cache is generated only once
	done        atomic.Bool // Atomic flag to determine generation status
}

// newDataset creates a new ethashb3 mining dataset and returns it as a plain Go
// interface to be usable in an LRU cache.
func newDataset(epoch uint64, epochLength uint64) *dataset {
	return &dataset{epoch: epoch, epochLength: epochLength}
}

// generate ensures that the dataset content is generated before use.
func (d *dataset) generate(dir string, limit int, lock bool, test bool) {
	d.once.Do(func() {
		// Mark the dataset generated after we're done. This is needed for remote
		defer d.done.Store(true)

		csize := cacheSize(d.epoch)
		dsize := datasetSize(d.epoch)
		seed := seedHash(d.epoch, d.epochLength)
		if test {
			csize = 1024
			dsize = 32 * 1024
		}
		// If we don't store anything on disk, generate and return
		if dir == "" {
			cache := make([]uint32, csize/4)
			generateCache(cache, d.epoch, d.epochLength, seed)

			d.dataset = make([]uint32, dsize/4)
			generateDataset(d.dataset, d.epoch, d.epochLength, cache)

			return
		}
		// Disk storage is needed, this will get fancy
		var endian string
		if !isLittleEndian() {
			endian = ".be"
		}
		path := filepath.Join(dir, fmt.Sprintf("full-R%d-%d-%x%s", algorithmRevision, d.epoch, seed[:8], endian))
		logger := log.New("epoch", d.epoch)

		// We're about to mmap the file, ensure that the mapping is cleaned up when the
		// cache becomes unused.
		runtime.SetFinalizer(d, (*dataset).finalizer)

		// Try to load the file from disk and memory map it
		var err error
		d.dump, d.mmap, d.dataset, err = memoryMap(path, lock)
		if err == nil {
			logger.Debug("Loaded old ethashb3 dataset from disk", "path", path)
			return
		}
		logger.Debug("Failed to load old ethashb3 dataset", "err", err)

		// No usable previous dataset available, create a new dataset file to fill
		cache := make([]uint32, csize/4)
		generateCache(cache, d.epoch, d.epochLength, seed)

		d.dump, d.mmap, d.dataset, err = memoryMapAndGenerate(path, dsize, lock, func(buffer []uint32) { generateDataset(buffer, d.epoch, d.epochLength, cache) })
		if err != nil {
			logger.Error("Failed to generate mapped ethashb3 dataset", "err", err)

			d.dataset = make([]uint32, dsize/4)
			generateDataset(d.dataset, d.epoch, d.epochLength, cache)
		}

		// Iterate over all full file instances, deleting any out of bounds (where epoch is below lower limit, or above upper limit).
		matches, _ := filepath.Glob(filepath.Join(dir, fmt.Sprintf("full-R%d*", algorithmRevision)))
		for _, file := range matches {
			var ar int   // algorithm revision
			var e uint64 // epoch
			var s string // seed
			if _, err := fmt.Sscanf(filepath.Base(file), "full-R%d-%d-%s"+endian, &ar, &e, &s); err != nil {
				// There is an unrecognized file in this directory.
				// See if the name matches the expected pattern of the legacy naming scheme.
				if _, err := fmt.Sscanf(filepath.Base(file), "full-R%d-%s"+endian, &ar, &s); err == nil {
					// This file matches the previous generation naming pattern (sans epoch).
					if err := os.Remove(file); err != nil {
						logger.Error("Failed to remove legacy ethashb3 full file", "file", file, "err", err)
					} else {
						logger.Warn("Deleted legacy ethashb3 full file", "path", file)
					}
				}
				// Else the file is unrecognized (unknown name format), leave it alone.
				continue
			}
			if e <= d.epoch-uint64(limit) || e > d.epoch+2 {
				if err := os.Remove(file); err == nil {
					logger.Debug("Deleted ethashb3 full file", "target.epoch", e, "file", file)
				} else {
					logger.Error("Failed to delete ethashb3 full file", "target.epoch", e, "file", file, "err", err)
				}
			}
		}
	})
}

// generated returns whether this particular dataset finished generating already
// or not (it may not have been started at all). This is useful for remote miners
// to default to verification caches instead of blocking on DAG generations.
func (d *dataset) generated() bool {
	return d.done.Load()
}

// finalizer closes any file handlers and memory maps open.
func (d *dataset) finalizer() {
	if d.mmap != nil {
		d.mmap.Unmap()
		d.dump.Close()
		d.mmap, d.dump = nil, nil
	}
}

// MakeCache generates a new ethashb3 cache and optionally stores it to disk.
func MakeCache(block uint64, epochLength uint64, dir string) {
	epoch := calcEpoch(block, epochLength)
	c := cache{epoch: epoch, epochLength: epochLength}
	c.generate(dir, math.MaxInt32, false, false)
}

// MakeDataset generates a new ethashb3 dataset and optionally stores it to disk.
func MakeDataset(block uint64, epochLength uint64, dir string) {
	epoch := calcEpoch(block, epochLength)
	d := dataset{epoch: epoch, epochLength: epochLength}
	d.generate(dir, math.MaxInt32, false, false)
}

// Mode defines the type and amount of PoW verification an ethashb3 engine makes.
type Mode uint

const (
	ModeNormal Mode = iota
	ModeShared
	ModeTest
	ModeFake
	ModePoissonFake
	ModeFullFake
)

func (m Mode) String() string {
	switch m {
	case ModeNormal:
		return "Normal"
	case ModeShared:
		return "Shared"
	case ModeTest:
		return "Test"
	case ModeFake:
		return "Fake"
	case ModePoissonFake:
		return "PoissonFake"
	case ModeFullFake:
		return "FullFake"
	}
	return "unknown"
}

// Config are the configuration parameters of the ethashb3.
type Config struct {
	CacheDir         string
	CachesInMem      int
	CachesOnDisk     int
	CachesLockMmap   bool
	DatasetDir       string
	DatasetsInMem    int
	DatasetsOnDisk   int
	DatasetsLockMmap bool
	PowMode          Mode

	// When set, notifications sent by the remote sealer will
	// be block header JSON objects instead of work package arrays.
	NotifyFull bool

	Log log.Logger `toml:"-"`
}

// EthashB3 is a consensus engine based on proof-of-work implementing the ethashb3
// algorithm.
type EthashB3 struct {
	config Config

	caches   *lru[*cache]   // In memory caches to avoid regenerating too often
	datasets *lru[*dataset] // In memory datasets to avoid regenerating too often

	// Mining related fields
	rand     *rand.Rand    // Properly seeded random source for nonces
	threads  int           // Number of threads to mine on if mining
	update   chan struct{} // Notification channel to update mining parameters
	hashrate metrics.Meter // Meter tracking the average hashrate
	remote   *remoteSealer

	// The fields below are hooks for testing
	shared    *EthashB3     // Shared PoW verifier to avoid cache regeneration
	fakeFail  uint64        // Block number which fails PoW check even in fake mode
	fakeDelay time.Duration // Time delay to sleep for before returning from verify

	lock      sync.Mutex // Ensures thread safety for the in-memory caches and mining fields
	closeOnce sync.Once  // Ensures exit channel will not be closed twice.
}

// New creates a full sized ethashb3 PoW scheme and starts a background thread for
// remote mining, also optionally notifying a batch of remote services of new work
// packages.
func New(config Config, notify []string, noverify bool) *EthashB3 {
	if config.Log == nil {
		config.Log = log.Root()
	}
	if config.CachesInMem <= 0 {
		config.Log.Warn("One ethashb3 cache must always be in memory", "requested", config.CachesInMem)
		config.CachesInMem = 1
	}
	if config.CacheDir != "" && config.CachesOnDisk > 0 {
		config.Log.Info("Disk storage enabled for ethashb3 caches", "dir", config.CacheDir, "count", config.CachesOnDisk)
	}
	if config.DatasetDir != "" && config.DatasetsOnDisk > 0 {
		config.Log.Info("Disk storage enabled for ethashb3 DAGs", "dir", config.DatasetDir, "count", config.DatasetsOnDisk)
	}
	ethash := &EthashB3{
		config:   config,
		caches:   newlru(config.CachesInMem, newCache),
		datasets: newlru(config.DatasetsInMem, newDataset),
		update:   make(chan struct{}),
		hashrate: metrics.NewMeter(),
	}
	if config.PowMode == ModeShared {
		ethash.shared = sharedEthash
	}
	ethash.remote = startRemoteSealer(ethash, notify, noverify)
	return ethash
}

// NewTester creates a small sized ethashb3 PoW scheme useful only for testing
// purposes.
func NewTester(notify []string, noverify bool) *EthashB3 {
	return New(Config{PowMode: ModeTest}, notify, noverify)
}

// NewFaker creates a ethashb3 consensus engine with a fake PoW scheme that accepts
// all blocks' seal as valid, though they still have to conform to the Ethereum
// consensus rules.
func NewFaker() *EthashB3 {
	return &EthashB3{
		config: Config{
			PowMode: ModeFake,
			Log:     log.Root(),
		},
	}
}

// NewFakeFailer creates a ethashb3 consensus engine with a fake PoW scheme that
// accepts all blocks as valid apart from the single one specified, though they
// still have to conform to the Ethereum consensus rules.
func NewFakeFailer(fail uint64) *EthashB3 {
	return &EthashB3{
		config: Config{
			PowMode: ModeFake,
			Log:     log.Root(),
		},
		fakeFail: fail,
	}
}

// NewFakeDelayer creates a ethashb3 consensus engine with a fake PoW scheme that
// accepts all blocks as valid, but delays verifications by some time, though
// they still have to conform to the Ethereum consensus rules.
func NewFakeDelayer(delay time.Duration) *EthashB3 {
	return &EthashB3{
		config: Config{
			PowMode: ModeFake,
			Log:     log.Root(),
		},
		fakeDelay: delay,
	}
}

// NewPoissonFaker creates a ethashb3 consensus engine with a fake PoW scheme that
// accepts all blocks as valid, but delays mining by some time based on miner.threads, though
// they still have to conform to the Ethereum consensus rules.
func NewPoissonFaker() *EthashB3 {
	return &EthashB3{
		config: Config{
			PowMode: ModePoissonFake,
			Log:     log.Root(),
		},
	}
}

// NewFullFaker creates an ethashb3 consensus engine with a full fake scheme that
// accepts all blocks as valid, without checking any consensus rules whatsoever.
func NewFullFaker() *EthashB3 {
	return &EthashB3{
		config: Config{
			PowMode: ModeFullFake,
			Log:     log.Root(),
		},
	}
}

// NewShared creates a full sized ethashb3 PoW shared between all requesters running
// in the same process.
func NewShared() *EthashB3 {
	return &EthashB3{shared: sharedEthash}
}

// Close closes the exit channel to notify all backend threads exiting.
func (ethashb3 *EthashB3) Close() error {
	return ethashb3.StopRemoteSealer()
}

// StopRemoteSealer stops the remote sealer
func (ethashb3 *EthashB3) StopRemoteSealer() error {
	ethashb3.closeOnce.Do(func() {
		// Short circuit if the exit channel is not allocated.
		if ethashb3.remote == nil {
			return
		}
		close(ethashb3.remote.requestExit)
		<-ethashb3.remote.exitCh
	})
	return nil
}

// cache tries to retrieve a verification cache for the specified block number
// by first checking against a list of in-memory caches, then against caches
// stored on disk, and finally generating one if none can be found.
func (ethashb3 *EthashB3) cache(block uint64) *cache {
	epochLength := calcEpochLength(block)
	epoch := calcEpoch(block, epochLength)
	current, future := ethashb3.caches.get(epoch, epochLength)

	// Wait for generation finish.
	current.generate(ethashb3.config.CacheDir, ethashb3.config.CachesOnDisk, ethashb3.config.CachesLockMmap, ethashb3.config.PowMode == ModeTest)

	// If we need a new future cache, now's a good time to regenerate it.
	if future != nil {
		go future.generate(ethashb3.config.CacheDir, ethashb3.config.CachesOnDisk, ethashb3.config.CachesLockMmap, ethashb3.config.PowMode == ModeTest)
	}
	return current
}

// dataset tries to retrieve a mining dataset for the specified block number
// by first checking against a list of in-memory datasets, then against DAGs
// stored on disk, and finally generating one if none can be found.
//
// If async is specified, not only the future but the current DAG is also
// generates on a background thread.
func (ethashb3 *EthashB3) dataset(block uint64, async bool) *dataset {
	// Retrieve the requested ethashb3 dataset
	epochLength := calcEpochLength(block)
	epoch := calcEpoch(block, epochLength)
	current, future := ethashb3.datasets.get(epoch, epochLength)

	// If async is specified, generate everything in a background thread
	if async && !current.generated() {
		go func() {
			current.generate(ethashb3.config.DatasetDir, ethashb3.config.DatasetsOnDisk, ethashb3.config.DatasetsLockMmap, ethashb3.config.PowMode == ModeTest)
			if future != nil {
				future.generate(ethashb3.config.DatasetDir, ethashb3.config.DatasetsOnDisk, ethashb3.config.DatasetsLockMmap, ethashb3.config.PowMode == ModeTest)
			}
		}()
	} else {
		// Either blocking generation was requested, or already done
		current.generate(ethashb3.config.DatasetDir, ethashb3.config.DatasetsOnDisk, ethashb3.config.DatasetsLockMmap, ethashb3.config.PowMode == ModeTest)
		if future != nil {
			go future.generate(ethashb3.config.DatasetDir, ethashb3.config.DatasetsOnDisk, ethashb3.config.DatasetsLockMmap, ethashb3.config.PowMode == ModeTest)
		}
	}
	return current
}

// Threads returns the number of mining threads currently enabled. This doesn't
// necessarily mean that mining is running!
func (ethashb3 *EthashB3) Threads() int {
	ethashb3.lock.Lock()
	defer ethashb3.lock.Unlock()

	return ethashb3.threads
}

// SetThreads updates the number of mining threads currently enabled. Calling
// this method does not start mining, only sets the thread count. If zero is
// specified, the miner will use all cores of the machine. Setting a thread
// count below zero is allowed and will cause the miner to idle, without any
// work being done.
func (ethashb3 *EthashB3) SetThreads(threads int) {
	ethashb3.lock.Lock()
	defer ethashb3.lock.Unlock()

	// If we're running a shared PoW, set the thread count on that instead
	if ethashb3.shared != nil {
		ethashb3.shared.SetThreads(threads)
		return
	}
	// Update the threads and ping any running seal to pull in any changes
	ethashb3.threads = threads
	select {
	case ethashb3.update <- struct{}{}:
	default:
	}
}

// Hashrate implements PoW, returning the measured rate of the search invocations
// per second over the last minute.
// Note the returned hashrate includes local hashrate, but also includes the total
// hashrate of all remote miner.
func (ethashb3 *EthashB3) Hashrate() float64 {
	// Short circuit if we are run the ethashb3 in normal/test mode.
	if ethashb3.config.PowMode != ModeNormal && ethashb3.config.PowMode != ModeTest {
		ms := ethashb3.hashrate.Snapshot()
		return ms.Rate1()
	}
	var res = make(chan uint64, 1)

	select {
	case ethashb3.remote.fetchRateCh <- res:
	case <-ethashb3.remote.exitCh:
		// Return local hashrate only if ethashb3 is stopped.
		ms := ethashb3.hashrate.Snapshot()
		return ms.Rate1()

	}

	// Gather total submitted hash rate of remote sealers.
	ms := ethashb3.hashrate.Snapshot()
	return ms.Rate1() + float64(<-res)
}

// APIs implements consensus.Engine, returning the user facing RPC APIs.
func (ethashb3 *EthashB3) APIs(chain consensus.ChainHeaderReader) []rpc.API {
	// In order to ensure backward compatibility, we exposes ethashb3 RPC APIs
	// to both eth and ethashb3 namespaces.
	return []rpc.API{
		{
			Namespace: "eth",
			Service:   &API{ethashb3},
		},
		{
			Namespace: "ethashb3",
			Service:   &API{ethashb3},
		},
	}
}

// SeedHash is the seed to use for generating a verification cache and the mining
// dataset.
func SeedHash(epoch uint64, epochLength uint64) []byte {
	return seedHash(epoch, epochLength)
}

// CalcEpochLength returns the epoch length for a given block number (ECIP-1099)
func CalcEpochLength(block uint64) uint64 {
	return calcEpochLength(block)
}

// CalcEpoch returns the epoch for a given block number (ECIP-1099)
func CalcEpoch(block uint64, epochLength uint64) uint64 {
	return calcEpoch(block, epochLength)
}
