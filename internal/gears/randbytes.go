/*
 * Copyright 2023 RPCPlatform Authors
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

package gears

import (
	"crypto/rand"
	"io"
)

func randBytes() [8]byte {
	var b [8]byte

	if _, err := io.ReadFull(rand.Reader, b[:]); err != nil {
		panic(err)
	}

	return b
}
