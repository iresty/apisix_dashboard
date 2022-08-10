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
package filter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/apache/apisix-dashboard/api/internal/log"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func RecoverHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("err；", err)
				//uuid := c.Writer.Header().Get("X-Request-Id")
				stack := stack(3)
				fmt.Printf("[Recovery] %s panic recovered:\n\n%s\n%s", timeFormat(time.Now()), err, stack)

				//log.With(zap.String("uuid", uuid))
				log.Errorf("[Recovery] %s panic recovered:\n\n%s\n%s", timeFormat(time.Now()), err, stack)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

func WrapGo(f func(...interface{}), args ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			stack := stack(3)
			log.Errorf("[Recovery] %s panic recovered:\n\n%s\n%s", timeFormat(time.Now()), err, stack)
		}
	}()
	f(args...)
}

func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

func timeFormat(t time.Time) string {
	var timeString = t.Format("2006/01/02 - 15:04:05")
	return timeString
}
