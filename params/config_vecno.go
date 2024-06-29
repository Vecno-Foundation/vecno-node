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
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/vars"
)

const VecnoChainId = 65357

var (
	// VecnoChainConfig is the chain parameters to run a node on the Vecno main network.
	VecnoChainConfig = &coregeth.CoreGethChainConfig{
		NetworkID:                 VecnoChainId,
		EthashB3:                  new(ctypes.EthashB3Config),
		ChainID:                   big.NewInt(VecnoChainId),
		SupportedProtocolVersions: vars.DefaultProtocolVersions,

		EIP2FBlock:   big.NewInt(0),
		EIP7FBlock:   big.NewInt(0),
		EIP150Block:  big.NewInt(0),
		EIP155Block:  big.NewInt(0),
		EIP160FBlock: big.NewInt(0),
		EIP161FBlock: big.NewInt(0),
		EIP170FBlock: big.NewInt(0),

		// Byzantium eq -- Enables Smart contracts
		EIP100FBlock: big.NewInt(70),
		EIP140FBlock: big.NewInt(70),
		EIP198FBlock: big.NewInt(70),
		EIP211FBlock: big.NewInt(70),
		EIP212FBlock: big.NewInt(70),
		EIP213FBlock: big.NewInt(70),
		EIP214FBlock: big.NewInt(70),
		//EIP649FBlock: big.NewInt(1001),
		EIP658FBlock: big.NewInt(70),

		// Constantinople eq
		EIP145FBlock:    big.NewInt(80),
		EIP1014FBlock:   big.NewInt(80),
		EIP1052FBlock:   big.NewInt(80),
		EIP1283FBlock:   big.NewInt(80),
		PetersburgBlock: big.NewInt(80),

		// Istanbul eq
		EIP152FBlock:  big.NewInt(90),
		EIP1108FBlock: big.NewInt(90),
		EIP1344FBlock: big.NewInt(90),
		EIP1884FBlock: big.NewInt(90),
		EIP2028FBlock: big.NewInt(90),
		EIP2200FBlock: big.NewInt(90),
		EIP2384FBlock: big.NewInt(90),

		// Berlin
		EIP2565FBlock: big.NewInt(100), // ModExp Gas Cost
		EIP2718FBlock: big.NewInt(100), // Typed Transaction Envelope
		EIP2929FBlock: big.NewInt(100), // Gas cost increases for state access opcodes
		EIP2930FBlock: big.NewInt(100), // Optional access lists

		// Veldin fork was used to enable rewards to miners for including uncle blocks on Vecno network.
		// Previously overlooked and unrewarded.
		HIPVeldinFBlock: big.NewInt(101),

		//Gaspar fork was used to upgrade the EVM to include new opcodes and features.
		HIPGasparFBlock: big.NewInt(101),

		// London + shanghai chain upgrades, aka Planned Eudora
		// TODO: move block numbers closer once testing has concluded
		//HIPEudoraFBlock: big.NewInt(13_524_557), // Vecno planned TX rewards change
		//EIP1559FBlock: big.NewInt(13_524_557), // EIP-1559 transactions`
		//EIP3541FBlock: big.NewInt(13_524_557), // EIP-3541 Reject code starting with 0xEF
		//EIP3855FBlock: big.NewInt(13_524_557), // PUSH0 instruction
		//EIP3860FBlock: big.NewInt(13_524_557), // Limit and meter initcode
		//EIP3198FBlock: big.NewInt(13_524_557), // BASEFEE Opcode
		//EIP3529FBlock: big.NewInt(13_524_557), // Reduction in refunds

		// Unplanned Upgrade, aka Olantis
		// EIP3651FBlock: big.NewInt(13_524_557), // Warm COINBASE (gas reprice)
		// EIP6049FBlock: big.NewInt(13_524_557), // Deprecate SELFDESTRUCT
		// EIP3541FBlock: big.NewInt(13_524_557), // Reject new contract code starting with the 0xEF byte

		// Spiral, aka Shanghai (partially)
		// EIP4399FBlock: nil, // Supplant DIFFICULTY with PREVRANDAO. Vecno  does not spec 4399 because it's still PoW, and 4399 is only applicable for the PoS system.
		// EIP4895FBlock: nil, // Beacon chain push withdrawals as operations

		// Dummy EIPs, unused by ethashb3 but used by forkid
		EIP3554FBlock: big.NewInt(13_524_557),
		//EIP4345FBlock: big.NewInt(27_200_177),
		//EIP5133FBlock: big.NewInt(40_725_107),

		// Define the planned 3 year decreasing rewards.
		BlockRewardSchedule: map[uint64]*big.Int{
			0:           big.NewInt(9 * vars.Ether),
			3_944_700:   big.NewInt(4.371893735 * vars.Ether),
			5_259_600:   big.NewInt(4.247434407 * vars.Ether),
			6_574_500:   big.NewInt(4.126518194 * vars.Ether),
			7_889_400:   big.NewInt(4.009044232 * vars.Ether),
			9_204_300:   big.NewInt(3.894914525 * vars.Ether),
			10_519_200:  big.NewInt(3.784033869 * vars.Ether),
			11_834_100:  big.NewInt(3.67630977 * vars.Ether),
			13_149_000:  big.NewInt(3.571652367 * vars.Ether),
			14_463_900:  big.NewInt(3.469974357 * vars.Ether),
			15_778_800:  big.NewInt(3.371190923 * vars.Ether),
			17_093_700:  big.NewInt(3.275219661 * vars.Ether),
			18_408_600:  big.NewInt(3.181980515 * vars.Ether),
			19_723_500:  big.NewInt(3.091395707 * vars.Ether),
			21_038_400:  big.NewInt(3.003389672 * vars.Ether),
			22_353_300:  big.NewInt(2.917888998 * vars.Ether),
			23_668_200:  big.NewInt(2.834822362 * vars.Ether),
			24_983_100:  big.NewInt(2.754120472 * vars.Ether),
			26_298_000:  big.NewInt(2.675716009 * vars.Ether),
			27_612_900:  big.NewInt(2.599543568 * vars.Ether),
			28_927_800:  big.NewInt(2.525539609 * vars.Ether),
			30_242_700:  big.NewInt(2.453642398 * vars.Ether),
			31_557_600:  big.NewInt(2.383791962 * vars.Ether),
			32_872_500:  big.NewInt(2.315930032 * vars.Ether),
			34_187_400:  big.NewInt(2.25 * vars.Ether),
			35_502_300:  big.NewInt(2.185946868 * vars.Ether),
			36_817_200:  big.NewInt(2.123717204 * vars.Ether),
			38_132_100:  big.NewInt(2.063259097 * vars.Ether),
			39_447_000:  big.NewInt(2.004522116 * vars.Ether),
			40_761_900:  big.NewInt(1.947457262 * vars.Ether),
			42_076_800:  big.NewInt(1.892016934 * vars.Ether),
			43_391_700:  big.NewInt(1.838154885 * vars.Ether),
			44_706_600:  big.NewInt(1.785826183 * vars.Ether),
			46_021_500:  big.NewInt(1.734987179 * vars.Ether),
			47_336_400:  big.NewInt(1.685595461 * vars.Ether),
			48_651_300:  big.NewInt(1.637609831 * vars.Ether),
			49_966_200:  big.NewInt(1.590990258 * vars.Ether),
			51_281_100:  big.NewInt(1.545697853 * vars.Ether),
			52_596_000:  big.NewInt(1.501694836 * vars.Ether),
			53_910_900:  big.NewInt(1.458944499 * vars.Ether),
			55_225_800:  big.NewInt(1.417411181 * vars.Ether),
			56_540_700:  big.NewInt(1.377060236 * vars.Ether),
			57_855_600:  big.NewInt(1.337858004 * vars.Ether),
			59_170_500:  big.NewInt(1.299771784 * vars.Ether),
			60_485_400:  big.NewInt(1.262769804 * vars.Ether),
			61_800_300:  big.NewInt(1.226821199 * vars.Ether),
			63_115_200:  big.NewInt(1.191895981 * vars.Ether),
			64_430_100:  big.NewInt(1.157965016 * vars.Ether),
			65_745_000:  big.NewInt(1.125 * vars.Ether),
			67_059_900:  big.NewInt(1.092973434 * vars.Ether),
			68_374_800:  big.NewInt(1.061858602 * vars.Ether),
			69_689_700:  big.NewInt(1.031629549 * vars.Ether),
			71_004_600:  big.NewInt(1.002261058 * vars.Ether),
			72_319_500:  big.NewInt(0.973728631 * vars.Ether),
			73_634_400:  big.NewInt(0.946008467 * vars.Ether),
			74_949_300:  big.NewInt(0.919077442 * vars.Ether),
			76_264_200:  big.NewInt(0.892913092 * vars.Ether),
			77_579_100:  big.NewInt(0.867493589 * vars.Ether),
			78_894_000:  big.NewInt(0.842797731 * vars.Ether),
			80_208_900:  big.NewInt(0.818804915 * vars.Ether),
			81_523_800:  big.NewInt(0.795495129 * vars.Ether),
			82_838_700:  big.NewInt(0.772848927 * vars.Ether),
			84_153_600:  big.NewInt(0.750847418 * vars.Ether),
			85_468_500:  big.NewInt(0.729472249 * vars.Ether),
			86_783_400:  big.NewInt(0.708705591 * vars.Ether),
			88_098_300:  big.NewInt(0.688530118 * vars.Ether),
			89_413_200:  big.NewInt(0.668929002 * vars.Ether),
			90_728_100:  big.NewInt(0.649885892 * vars.Ether),
			92_043_000:  big.NewInt(0.631384902 * vars.Ether),
			93_357_900:  big.NewInt(0.6134106 * vars.Ether),
			94_672_800:  big.NewInt(0.595947991 * vars.Ether),
			95_987_700:  big.NewInt(0.578982508 * vars.Ether),
			97_302_600:  big.NewInt(0.5625 * vars.Ether),
			98_617_500:  big.NewInt(0.546486717 * vars.Ether),
			99_932_400:  big.NewInt(0.530929301 * vars.Ether),
			101_247_300: big.NewInt(0.515814774 * vars.Ether),
			102_562_200: big.NewInt(0.501130529 * vars.Ether),
			103_877_100: big.NewInt(0.486864316 * vars.Ether),
			105_192_000: big.NewInt(0.473004234 * vars.Ether),
			106_506_900: big.NewInt(0.459538721 * vars.Ether),
			107_821_800: big.NewInt(0.446456546 * vars.Ether),
			109_136_700: big.NewInt(0.433746795 * vars.Ether),
			110_451_600: big.NewInt(0.421398865 * vars.Ether),
			111_766_500: big.NewInt(0.409402458 * vars.Ether),
			113_081_400: big.NewInt(0.397747564 * vars.Ether),
			114_396_300: big.NewInt(0.386424463 * vars.Ether),
			115_711_200: big.NewInt(0.375423709 * vars.Ether),
			117_026_100: big.NewInt(0.364736125 * vars.Ether),
			118_341_000: big.NewInt(0.354352795 * vars.Ether),
			119_655_900: big.NewInt(0.344265059 * vars.Ether),
			120_970_800: big.NewInt(0.334464501 * vars.Ether),
			122_285_700: big.NewInt(0.324942946 * vars.Ether),
			123_600_600: big.NewInt(0.315692451 * vars.Ether),
			124_915_500: big.NewInt(0.3067053 * vars.Ether),
			126_230_400: big.NewInt(0.297973995 * vars.Ether),
			127_545_300: big.NewInt(0.289491254 * vars.Ether),
			128_860_200: big.NewInt(0.28125 * vars.Ether),
			130_175_100: big.NewInt(0.273243358 * vars.Ether),
			131_490_000: big.NewInt(0.26546465 * vars.Ether),
			132_804_900: big.NewInt(0.257907387 * vars.Ether),
			134_119_800: big.NewInt(0.250565264 * vars.Ether),
			135_434_700: big.NewInt(0.243432158 * vars.Ether),
			136_749_600: big.NewInt(0.236502117 * vars.Ether),
			138_064_500: big.NewInt(0.229769361 * vars.Ether),
			139_379_400: big.NewInt(0.223228273 * vars.Ether),
			140_694_300: big.NewInt(0.216873397 * vars.Ether),
			142_009_200: big.NewInt(0.210699433 * vars.Ether),
			143_324_100: big.NewInt(0.204701229 * vars.Ether),
			144_639_000: big.NewInt(0.98873782 * vars.Ether),
			145_953_900: big.NewInt(0.193212232 * vars.Ether),
			147_268_800: big.NewInt(0.187711854 * vars.Ether),
			148_583_700: big.NewInt(0.182368062 * vars.Ether),
			149_898_600: big.NewInt(0.177176398 * vars.Ether),
			151_213_500: big.NewInt(0.17213253 * vars.Ether),
			152_528_400: big.NewInt(0.167232251 * vars.Ether),
			153_843_300: big.NewInt(0.162471473 * vars.Ether),
			155_158_200: big.NewInt(0.157846226 * vars.Ether),
			156_473_100: big.NewInt(0.15335265 * vars.Ether),
			157_788_000: big.NewInt(0.148986998 * vars.Ether),
			159_102_900: big.NewInt(0.144745627 * vars.Ether),
			160_417_800: big.NewInt(0.140625 * vars.Ether),
			161_732_700: big.NewInt(0.136621679 * vars.Ether),
			163_047_600: big.NewInt(0.132732325 * vars.Ether),
			164_362_500: big.NewInt(0.128953694 * vars.Ether),
			165_677_400: big.NewInt(0.125282632 * vars.Ether),
			166_992_300: big.NewInt(0.121716079 * vars.Ether),
			168_307_200: big.NewInt(0.118251058 * vars.Ether),
			169_622_100: big.NewInt(0.11488468 * vars.Ether),
			170_937_000: big.NewInt(0.111614136 * vars.Ether),
			172_251_900: big.NewInt(0.108436699 * vars.Ether),
			173_566_800: big.NewInt(0.105349716 * vars.Ether),
			174_881_700: big.NewInt(0.102350614 * vars.Ether),
			176_196_600: big.NewInt(0.099436891 * vars.Ether),
			177_511_500: big.NewInt(0.096606116 * vars.Ether),
			178_826_400: big.NewInt(0.093855927 * vars.Ether),
			180_141_300: big.NewInt(0.091184031 * vars.Ether),
			181_456_200: big.NewInt(0.088588199 * vars.Ether),
			182_771_100: big.NewInt(0.086066265 * vars.Ether),
			184_086_000: big.NewInt(0.083616125 * vars.Ether),
			185_400_900: big.NewInt(0.081235736 * vars.Ether),
			186_715_800: big.NewInt(0.078923113 * vars.Ether),
			188_030_700: big.NewInt(0.076676325 * vars.Ether),
			189_345_600: big.NewInt(0.074493499 * vars.Ether),
			190_660_500: big.NewInt(0.072372814 * vars.Ether),
			191_975_400: big.NewInt(0.0703125 * vars.Ether),
			193_290_300: big.NewInt(0.06831084 * vars.Ether),
			194_605_200: big.NewInt(0.066366163 * vars.Ether),
			195_920_100: big.NewInt(0.064476847 * vars.Ether),
			197_235_000: big.NewInt(0.062641316 * vars.Ether),
			198_549_900: big.NewInt(0.060858039 * vars.Ether),
			199_864_800: big.NewInt(0.059125529 * vars.Ether),
			201_179_700: big.NewInt(0.05744234 * vars.Ether),
			202_494_600: big.NewInt(0.055807068 * vars.Ether),
			203_809_500: big.NewInt(0.054218349 * vars.Ether),
			205_124_400: big.NewInt(0.052674858 * vars.Ether),
			206_439_300: big.NewInt(0.051175307 * vars.Ether),
			207_754_200: big.NewInt(0.049718446 * vars.Ether),
			209_069_100: big.NewInt(0.048303058 * vars.Ether),
			210_384_000: big.NewInt(0.046927964 * vars.Ether),
			211_698_900: big.NewInt(0.045592016 * vars.Ether),
			213_013_800: big.NewInt(0.044294099 * vars.Ether),
			214_328_700: big.NewInt(0.043033132 * vars.Ether),
			215_643_600: big.NewInt(0.041808063 * vars.Ether),
			216_958_500: big.NewInt(0.040617868 * vars.Ether),
			218_273_400: big.NewInt(0.039461556 * vars.Ether),
			219_588_300: big.NewInt(0.038338162 * vars.Ether),
			220_903_200: big.NewInt(0.037246749 * vars.Ether),
			222_218_100: big.NewInt(0.036186407 * vars.Ether),
			223_533_000: big.NewInt(0.03515625 * vars.Ether),
			224_847_900: big.NewInt(0.03415542 * vars.Ether),
			226_162_800: big.NewInt(0.033183081 * vars.Ether),
			227_477_700: big.NewInt(0.032238423 * vars.Ether),
			228_792_600: big.NewInt(0.031320658 * vars.Ether),
			230_107_500: big.NewInt(0.03042902 * vars.Ether),
			231_422_400: big.NewInt(0.029562765 * vars.Ether),
			232_737_300: big.NewInt(0.02872117 * vars.Ether),
			234_052_200: big.NewInt(0.027903534 * vars.Ether),
			235_367_100: big.NewInt(0.027109175 * vars.Ether),
			236_682_000: big.NewInt(0.026337429 * vars.Ether),
			237_996_900: big.NewInt(0.025587654 * vars.Ether),
			239_311_800: big.NewInt(0.024859223 * vars.Ether),
			240_626_700: big.NewInt(0.024151529 * vars.Ether),
			241_941_600: big.NewInt(0.023463982 * vars.Ether),
			243_256_500: big.NewInt(0.022796008 * vars.Ether),
			244_571_400: big.NewInt(0.02214705 * vars.Ether),
			245_886_300: big.NewInt(0.021516566 * vars.Ether),
			247_201_200: big.NewInt(0.020904031 * vars.Ether),
			248_516_100: big.NewInt(0.020308934 * vars.Ether),
			249_831_000: big.NewInt(0.019730778 * vars.Ether),
			251_145_900: big.NewInt(0.019169081 * vars.Ether),
			252_460_800: big.NewInt(0.018623375 * vars.Ether),
			253_775_700: big.NewInt(0.018093203 * vars.Ether),
			255_090_600: big.NewInt(0.017578125 * vars.Ether),
			256_405_500: big.NewInt(0.01707771 * vars.Ether),
			257_720_400: big.NewInt(0.016591541 * vars.Ether),
			259_035_300: big.NewInt(0.016119212 * vars.Ether),
			260_350_200: big.NewInt(0.015660329 * vars.Ether),
			261_665_100: big.NewInt(0.01521451 * vars.Ether),
			262_980_000: big.NewInt(0.014781382 * vars.Ether),
			264_294_900: big.NewInt(0.014360585 * vars.Ether),
			265_609_800: big.NewInt(0.013951767 * vars.Ether),
			266_924_700: big.NewInt(0.013554587 * vars.Ether),
			268_239_600: big.NewInt(0.013168715 * vars.Ether),
			269_554_500: big.NewInt(0.012793827 * vars.Ether),
			270_869_400: big.NewInt(0.012429611 * vars.Ether),
			272_184_300: big.NewInt(0.011398004 * vars.Ether),
			273_499_200: big.NewInt(0.01 * vars.Ether),
		},

		TrustedCheckpoint: &ctypes.TrustedCheckpoint{
			BloomRoot:    common.HexToHash(""),
			CHTRoot:      common.HexToHash(""),
			SectionHead:  common.HexToHash(""),
			SectionIndex: 0,
		},

		RequireBlockHashes: map[uint64]common.Hash{
			0: common.HexToHash("0x529e9a367fd6ad8fa08dfbd8c1b26b766c46e287584cf714d6c786c5a52e6a4d"),
		},

		// No go
		EIP4399FBlock:           nil,
		EIP4895FBlock:           nil,
		TerminalTotalDifficulty: nil,
	}
)

func init() {
	// override the default bomb schedule
	eip2384 := VecnoChainConfig.EIP2384FBlock.Uint64()
	VecnoChainConfig.SetEthashEIP2384Transition(&eip2384)
}
