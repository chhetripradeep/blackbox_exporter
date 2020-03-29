// Copyright 2016 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build !windows

package prober

import (
	"net"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/blackbox_exporter/config"
)

func setupDialer(dialProtocol, dialTarget string, conn syscall.RawConn, module config.Module, logger log.Logger) *net.Dialer {
	return &net.Dialer{
		Control: func(dialProtocol, dialTarget string, conn syscall.RawConn) error {
			return conn.Control(func(fd uintptr) {
				if module.TCP.TOS != 0 {
					level.Info(logger).Log("msg", "Setting TOS", "TOS", module.TCP.TOS)
					err := syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_TOS, module.TCP.TOS)
					if err != nil {
						level.Error(logger).Log("msg", "Could not set TOS", "err", err)
						return
					}
				}
			})
		},
	}
}
