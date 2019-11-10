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
// but WITHOUT ANY WARRANTY; witho
// ut even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with the go-demeton library.  If not, see <http://www.gnu.org/licenses/>.
//

package rpc

import (
	"github.com/cyber-demeton/go-demeton/metrics"
)

// Metrics for rpc
var (
	metricsRPCCounter = metrics.NewMeter("deb.rpc.request")

	metricsAccountStateSuccess = metrics.NewMeter("deb.rpc.account.success")
	metricsAccountStateFailed  = metrics.NewMeter("deb.rpc.account.failed")

	metricsSendTxSuccess = metrics.NewMeter("deb.rpc.sendTx.success")
	metricsSendTxFailed  = metrics.NewMeter("deb.rpc.sendTx.failed")

	metricsSignTxSuccess = metrics.NewMeter("deb.rpc.signTx.success")
	metricsSignTxFailed  = metrics.NewMeter("deb.rpc.signTx.failed")

	metricsUnlockSuccess = metrics.NewMeter("deb.rpc.unlock.success")
	metricsUnlockFailed  = metrics.NewMeter("deb.rpc.unlock.failed")
)
