/*
 * Copyright 2022 RPCPlatform Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package balancer

import (
	"github.com/nexcode/rpcplatform/internal/gears"
	"google.golang.org/grpc/balancer"
)

func (*pickerBuilder) makePicker(connInfoArr []*connInfo, totalWeight int) *picker {
	picker := picker{
		subConns: make([]balancer.SubConn, 0, totalWeight),
	}

	for {
		prevLen := len(picker.subConns)

		for _, connInfo := range connInfoArr {
			if connInfo.count < connInfo.factor && connInfo.weight > 0 {
				picker.subConns = append(picker.subConns, connInfo.subConn)
				connInfo.weight--
				connInfo.count++
			}
		}

		if totalWeight == len(picker.subConns) {
			break
		}

		if prevLen == len(picker.subConns) {
			for _, connInfo := range connInfoArr {
				connInfo.count = 0
			}
		}
	}

	picker.next = gears.Intn(totalWeight)
	return &picker
}
