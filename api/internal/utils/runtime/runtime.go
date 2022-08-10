/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package runtime

import (
	"net/http"
	"runtime"

	"github.com/apache/apisix-dashboard/api/internal/log"
)

var (
	ActuallyPanic = true
)

var PanicHandlers = []func(interface{}){logPanic}

func HandlePanic(additionalHandlers ...func(interface{})) {
	if err := recover(); err != nil {
		for _, fn := range PanicHandlers {
			fn(err)
		}
		for _, fn := range additionalHandlers {
			fn(err)
		}
		if ActuallyPanic {
			panic(err)
		}
	}
}

func logPanic(r interface{}) {
	if r == http.ErrAbortHandler {
		return
	}

	const size = 32 << 10
	stacktrace := make([]byte, size)
	stacktrace = stacktrace[:runtime.Stack(stacktrace, false)]
	if _, ok := r.(string); ok {
		log.Errorf("observed a panic: %s\n%s", r, stacktrace)
	} else {
		log.Errorf("observed a panic: %#v (%v)\n%s", r, r, stacktrace)
	}
}
