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

package storage

import (
	metrics "github.com/cyber-demeton/go-demeton/metrics"
)

// Metrics for storage
var (
	// rocksdb metrics
	metricsRocksdbFlushTime = metrics.NewGauge("deb.rocksdb.flushtime")
	metricsRocksdbFlushLen  = metrics.NewGauge("deb.rocksdb.flushlen")

	metricsBlocksdbCacheSize       = metrics.NewGauge("deb.rocksdb.cache.size")       //block_cache->GetUsage()
	metricsBlocksdbCachePinnedSize = metrics.NewGauge("deb.rocksdb.cachepinned.size") //block_cache->GetPinnedUsage()
	metricsBlocksdbTableReaderMem  = metrics.NewGauge("deb.rocksdb.tablereader.mem")  //estimate-table-readers-mem
	metricsBlocksdbAllMemTables    = metrics.NewGauge("deb.rocksdb.alltables.mem")    //cur-size-all-mem-tables
)
