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
package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/sony/sonyflake"
	"github.com/yuin/gopher-lua/parse"
)

var _sf *sonyflake.Sonyflake

func init() {
	saltStr, ok := os.LookupEnv("FLAKE_SALT")
	var salt uint16
	if ok {
		i, err := strconv.ParseUint(saltStr, 10, 16)
		if err != nil {
			panic(err)
		}
		salt = uint16(i)
	}
	ips, err := getLocalIPs()
	if err != nil {
		panic(err)
	}
	_sf = sonyflake.NewSonyflake(sonyflake.Settings{
		MachineID: func() (u uint16, e error) {
			return sumIPs(ips) + salt, nil
		},
	})
	if _sf == nil {
		panic("sonyflake init failed")
	}
}

func sumIPs(ips []net.IP) uint16 {
	total := 0
	for _, ip := range ips {
		for i := range ip {
			total += int(ip[i])
		}
	}
	return uint16(total)
}

func getLocalIPs() ([]net.IP, error) {
	var ips []net.IP
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}
	for _, a := range addrs {
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			ips = append(ips, ipNet.IP)
		}
	}
	return ips, nil
}

func GetFlakeUid() uint64 {
	uid, err := _sf.NextID()
	if err != nil {
		panic("get sony flake uid failed:" + err.Error())
	}
	return uid
}

func GetFlakeUidStr() string {
	return strconv.FormatUint(GetFlakeUid(), 10)
}

func InterfaceToString(val interface{}) string {
	if val == nil {
		return ""
	}
	str := fmt.Sprintf("%v", val)
	return str
}

// Note: json.Marshal and json.Unmarshal may cause the precision loss
func ObjectClone(origin, copy interface{}) error {
	byt, err := json.Marshal(origin)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byt, copy)
	return err
}

func GenLabelMap(label string) (map[string]struct{}, error) {
	var err = errors.New("malformed label")
	mp := make(map[string]struct{})

	if label == "" {
		return mp, nil
	}

	labels := strings.Split(label, ",")
	for _, l := range labels {
		kv := strings.Split(l, ":")
		if len(kv) == 2 {
			if kv[0] == "" || kv[1] == "" {
				return nil, err
			}

			// Because the labels may contain the same key, like this: label=version:v1,version:v2
			// we need to combine them as a map's key
			mp[l] = struct{}{}
		} else if len(kv) == 1 {
			if kv[0] == "" {
				return nil, err
			}

			mp[kv[0]] = struct{}{}
		} else {
			return nil, err
		}
	}

	return mp, nil
}

func LabelContains(labels map[string]string, reqLabels map[string]struct{}) bool {
	if len(reqLabels) == 0 {
		return true
	}

	for k, v := range labels {
		// first check the key
		if _, exist := reqLabels[k]; exist {
			return true
		}

		// second check the key:value
		if _, exist := reqLabels[k+":"+v]; exist {
			return true
		}
	}

	return false
}

// ValidateLuaCode validates lua syntax for input code, return nil
// if passed, otherwise a non-nil error will be returned
func ValidateLuaCode(code string) error {
	_, err := parse.Parse(strings.NewReader(code), "<string>")
	return err
}

func StringSliceContains(a, b []string) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	for i := range a {
		for j := range b {
			if a[i] == b[j] {
				return true
			}
		}
	}

	return false
}

//
func StringSliceEqual(a, b []string) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	sort.Strings(a)
	sort.Strings(b)

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// value compare
func ValueEqual(a interface{}, b interface{}) bool {
	aBytes, err := json.Marshal(a)
	if err != nil {
		return false
	}
	bBytes, err := json.Marshal(b)
	if err != nil {
		return false
	}
	return bytes.Equal(aBytes, bBytes)
}

func PathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
