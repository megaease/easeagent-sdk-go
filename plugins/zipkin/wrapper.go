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

package zipkin

import (
	"net/http"

	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
)

type (
	// HTTPHandlerWrapper is the wrapper of http.Handler.
	HTTPHandlerWrapper struct {
		handlerFunc http.HandlerFunc
	}

	// HTTPClientWrapper is the wrapper of http.Client.
	HTTPClientWrapper struct {
		client *zipkinhttp.Client
	}
)

func (h *HTTPHandlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handlerFunc(w, r)
}

// Do implements plugins.HTTPDoer.
func (c *HTTPClientWrapper) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}
