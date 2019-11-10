// Copyright (C) 2017 go-demeton authors
//
// This file is part of the go-demeton library.
//
// the go-demeton library is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// the go-demeton library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with the go-demeton library.  If not, see <http://www.gnu.org/licenses/>.
//

package core

import (
	metrics "github.com/cyber-demeton/go-demeton/metrics"
)

// Metrics for core
var (
	// block metrics
	metricsBlockHeightGauge      = metrics.NewGauge("deb.block.height")
	metricsBlocktailHashGauge    = metrics.NewGauge("deb.block.tailhash")
	metricsBlockRevertTimesGauge = metrics.NewGauge("deb.block.revertcount")
	metricsBlockRevertMeter      = metrics.NewMeter("deb.block.revert")
	metricsBlockOnchainTimer     = metrics.NewTimer("deb.block.onchain")
	metricsTxOnchainTimer        = metrics.NewTimer("deb.transaction.onchain")
	metricsBlockPackTxTime       = metrics.NewGauge("deb.block.packtx")

	// block_pool metrics
	metricsCachedNewBlock      = metrics.NewGauge("deb.block.new.cached")
	metricsCachedDownloadBlock = metrics.NewGauge("deb.block.download.cached")
	metricsLruPoolCacheBlock   = metrics.NewGauge("deb.block.lru.poolcached")
	metricsLruCacheBlock       = metrics.NewGauge("deb.block.lru.blocks")
	metricsLruTailBlock        = metrics.NewGauge("deb.block.lru.tailblock")

	metricsDuplicatedBlock   = metrics.NewCounter("deb.block.duplicated")
	metricsInvalidBlock      = metrics.NewCounter("deb.block.invalid")
	metricsTxsInBlock        = metrics.NewGauge("deb.block.txs")
	metricsBlockVerifiedTime = metrics.NewGauge("deb.block.executed")
	metricsTxVerifiedTime    = metrics.NewGauge("deb.tx.executed")
	metricsTxPackedCount     = metrics.NewGauge("deb.tx.packed")
	metricsTxUnpackedCount   = metrics.NewGauge("deb.tx.unpacked")
	metricsTxGivebackCount   = metrics.NewGauge("deb.tx.giveback")

	// txpool metrics
	metricsReceivedTx                      = metrics.NewGauge("deb.txpool.received")
	metricsCachedTx                        = metrics.NewGauge("deb.txpool.cached")
	metricsBucketTx                        = metrics.NewGauge("deb.txpool.bucket")
	metricsCandidates                      = metrics.NewGauge("deb.txpool.candidates")
	metricsInvalidTx                       = metrics.NewCounter("deb.txpool.invalid")
	metricsDuplicateTx                     = metrics.NewCounter("deb.txpool.duplicate")
	metricsTxPoolBelowGasPrice             = metrics.NewCounter("deb.txpool.below_gas_price")
	metricsTxPoolOutOfGasLimit             = metrics.NewCounter("deb.txpool.out_of_gas_limit")
	metricsTxPoolGasLimitLessOrEqualToZero = metrics.NewCounter("deb.txpool.gas_limit_less_equal_zero")

	// transaction metrics
	metricsTxSubmit     = metrics.NewMeter("deb.transaction.submit")
	metricsTxExecute    = metrics.NewMeter("deb.transaction.execute")
	metricsTxExeSuccess = metrics.NewMeter("deb.transaction.execute.success")
	metricsTxExeFailed  = metrics.NewMeter("deb.transaction.execute.failed")

	// event metrics
	metricsCachedEvent = metrics.NewGauge("deb.event.cached")

	// unexpect behavior
	metricsUnexpectedBehavior = metrics.NewGauge("deb.unexpected")
)
