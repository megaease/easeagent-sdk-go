/**
 * Copyright 2022 MegaEase
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
package zero

import (
	"net/http"

	"github.com/megaease/easeagent-sdk-go/stdlib"
)

// ServeDefault is the same with stdlib ServeDefault.
// The caller must call it to activate default agent.
var ServeDefault = stdlib.ServeDefaultAgent

// EaseMeshHandler wraps handler of go-zero as middleware.
func EaseMeshHandler(next http.HandlerFunc) http.HandlerFunc {
	return stdlib.DefaultAgent.WrapHandleFunc(next)
}
