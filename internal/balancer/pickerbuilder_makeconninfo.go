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
	"cmp"
	"math"
	"slices"

	"github.com/nexcode/rpcplatform/internal/config"
	"github.com/nexcode/rpcplatform/internal/grpcattrs"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

type connInfo struct {
	subConn  balancer.SubConn
	priority int
	weight   int
	factor   int
	count    int
}

func (*pickerBuilder) makeConnInfo(pickerInfo base.PickerBuildInfo) ([]*connInfo, int) {
	connInfoArr := make([]*connInfo, 0, len(pickerInfo.ReadySCs))

	var totalWeight int
	var config *config.Client

	for subConn, subConnInfo := range pickerInfo.ReadySCs {
		if config == nil {
			config = grpcattrs.GetClientConfig(subConnInfo.Address.Attributes)
		}

		attributes := grpcattrs.GetAttributes(subConnInfo.Address.Attributes)

		if attributes.BalancerWeight <= 0 {
			continue
		}

		connInfoArr = append(connInfoArr, &connInfo{
			subConn:  subConn,
			priority: attributes.BalancerPriority,
			weight:   attributes.BalancerWeight,
		})
	}

	slices.SortFunc(connInfoArr, func(a, b *connInfo) int {
		return cmp.Compare(b.priority, a.priority)
	})

	if config.MaxActiveServers > 0 && config.MaxActiveServers < len(connInfoArr) {
		connInfoArr = connInfoArr[:config.MaxActiveServers]
	}

	for _, connInfo := range connInfoArr {
		connInfo.factor = int(math.Ceil(float64(connInfo.weight) / float64(len(connInfoArr))))
		totalWeight += connInfo.weight
	}

	return connInfoArr, totalWeight
}
