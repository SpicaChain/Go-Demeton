// Copyright (C) 2018 go-demeton authors
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

package net

import (
	"fmt"

	metrics "github.com/cyber-demeton/go-demeton/metrics"
)

// Metrics map for different in/out network msg types
var (
	metricsPacketsIn = metrics.NewMeter("deb.net.packets.in")
	metricsBytesIn   = metrics.NewMeter("deb.net.bytes.in")

	metricsPacketsOut = metrics.NewMeter("deb.net.packets.out")
	metricsBytesOut   = metrics.NewMeter("deb.net.bytes.out")
)

func metricsPacketsInByMessageName(messageName string, size uint64) {
	meter := metrics.NewMeter(fmt.Sprintf("deb.net.packets.in.%s", messageName))
	meter.Mark(1)

	meter = metrics.NewMeter(fmt.Sprintf("deb.net.bytes.in.%s", messageName))
	meter.Mark(int64(size))
}

func metricsPacketsOutByMessageName(messageName string, size uint64) {
	meter := metrics.NewMeter(fmt.Sprintf("deb.net.packets.out.%s", messageName))
	meter.Mark(1)

	meter = metrics.NewMeter(fmt.Sprintf("deb.net.bytes.out.%s", messageName))
	meter.Mark(int64(size))
}
