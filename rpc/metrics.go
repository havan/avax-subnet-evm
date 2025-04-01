// (c) 2019-2020, Ava Labs, Inc.
//
// This file is a derived work, based on the go-ethereum library whose original
// notices appear below.
//
// It is distributed under a license compatible with the licensing terms of the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********
// Copyright 2020 The go-ethereum Authors
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

package rpc

import (
	"fmt"
	"time"

	"github.com/ava-labs/libevm/metrics"

	// Force libevm metrics of the same name to be registered first.
	_ "github.com/ava-labs/libevm/rpc"
)

// ====== If resolving merge conflicts ======
//
// All calls to metrics.NewRegistered*() for metrics also defined in libevm/rpc have
// been replaced with metrics.GetOrRegister*() to get metrics already registered in
// libevm/rpc or register them here otherwise. These replacements ensure the same
// metrics are shared between the two packages.
var (
	rpcRequestGauge        = metrics.GetOrRegisterGauge("rpc/requests", nil)
	successfulRequestGauge = metrics.GetOrRegisterGauge("rpc/success", nil)
	failedRequestGauge     = metrics.GetOrRegisterGauge("rpc/failure", nil)

	// serveTimeHistName is the prefix of the per-request serving time histograms.
	serveTimeHistName = "rpc/duration"

	rpcServingTimer = metrics.GetOrRegisterTimer("rpc/duration/all", nil)
)

// updateServeTimeHistogram tracks the serving time of a remote RPC call.
func updateServeTimeHistogram(method string, success bool, elapsed time.Duration) {
	note := "success"
	if !success {
		note = "failure"
	}
	h := fmt.Sprintf("%s/%s/%s", serveTimeHistName, method, note)
	sampler := func() metrics.Sample {
		return metrics.ResettingSample(
			metrics.NewExpDecaySample(1028, 0.015),
		)
	}
	metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(elapsed.Nanoseconds())
}
