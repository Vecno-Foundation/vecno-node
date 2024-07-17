// Copyright 2019 The multi-geth Authors
// This file is part of the multi-geth library.
//
// The multi-geth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The multi-geth library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the multi-geth library. If not, see <http://www.gnu.org/licenses/>.

package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
)

var VecnoGenesisHash = common.HexToHash("0x0bf5cdeea9c53d84d071bac6d30570baef6ae155ea92f20f2f90833e7981092a")

// DefaultVecnoGenesisBlock returns the Vecno Network genesis block.
func DefaultVecnoGenesisBlock() *genesisT.Genesis {
	return &genesisT.Genesis{
		Config:     VecnoChainConfig,
		Nonce:      0x0,
		ExtraData:  hexutil.MustDecode("0x4a756c69616e20417373616e6765206c616e647320696e204175737472616c69612061667465722077616c6b696e672066726565"),
		GasLimit:   0x1E8480,
		Difficulty: big.NewInt(0x20000),
		Timestamp:  1719482400,
		Alloc:      genesisT.DecodePreAlloc(vecnoAllocData),
	}
}
