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
package mutations

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
)

var (
	UncleBlockReward    = big.NewInt(0)
	DevRewardPercentage = 10                                           // Set the developer reward percentage (10%)
	DevWalletAddress    = "0x44A9161474953cEe9CEd3C2a518B52B22CdAbeC5" // Replace with the actual developer wallet address
)

// Function to calculate the developer reward based on the block reward
func calculateDevReward(blockReward *big.Int) *big.Int {
	devReward := new(big.Int).Mul(blockReward, big.NewInt(int64(DevRewardPercentage)))
	devReward = new(big.Int).Div(devReward, big.NewInt(100))
	return devReward
}

func GetRewardsVecno(config ctypes.ChainConfigurator, header *types.Header, uncles []*types.Header, txs []*types.Transaction) (*big.Int, *big.Int, []*big.Int) {
	// Select the correct block minerReward based on chain progression
	blockReward := ctypes.EthashBlockReward(config, header.Number)
	minerReward := new(big.Int).Set(blockReward)
	devReward := calculateDevReward(blockReward)
	uncleReward := new(big.Int).Set(UncleBlockReward)
	uncleCount := new(big.Int).SetUint64(uint64(len(uncles)))
	blockFeeReward := new(big.Int)

	// Collect the fee for all transactions.
	for _, tx := range txs {
		gas := new(big.Int).SetUint64(tx.Gas())
		gasPrice := tx.GasPrice()
		blockFeeReward.Add(blockFeeReward, new(big.Int).Mul(gas, gasPrice))
	}

	if len(uncles) == 0 { // If no uncles, the miner gets the entire block fee.
		minerReward.Add(minerReward, blockFeeReward)
	} else if config.IsEnabled(config.GetEIP2929Transition, header.Number) { // During Berlin block, each miner and uncles are rewarded the block fee.
		uncleReward.Add(uncleReward, blockFeeReward)
		minerReward.Add(minerReward, blockFeeReward)
	} else if config.IsEnabled(config.GetEthashHomesteadTransition, header.Number) { // Until Berlin block, Miners and Uncles are rewarded for the amount of uncles generated.
		uncleReward.Add(uncleReward, blockFeeReward)
		uncleReward.Mul(uncleReward, uncleCount)
		minerReward.Add(minerReward, uncleReward)
	}

	uncleRewards := make([]*big.Int, len(uncles))
	for i := range uncles {
		uncleRewards[i] = uncleReward

		if config.IsEnabled(config.GetHIPVeldinTransition, header.Number) {
			minerReward.Add(minerReward, uncleReward)
		}
	}

	return minerReward, devReward, uncleRewards
}
