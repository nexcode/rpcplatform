/*
 * Copyright 2025 RPCPlatform Authors
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

package rpcplatform

import (
	"github.com/nexcode/rpcplatform/internal/config"
)

// PlatformOption is used to pass [PlatformOptions] when creating a new [RPCPlatform]
type PlatformOption = func(*config.Platform)

// ClientOption is used to pass [ClientOptions] when creating a new [Client]
type ClientOption = func(*config.Client)

// ServerOption is used to pass [ServerOptions] when creating a new [Server]
type ServerOption = func(*config.Server)
