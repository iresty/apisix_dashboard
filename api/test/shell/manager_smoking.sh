#!/usr/bin/env bash

#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

set -ex

helpFunction()
{
   echo ""
   echo "Usage: $0 -s true"
   echo -e "\t-s whether skip docker, true or false"
   echo -e "\t-h helper info"
   exit 1
}

while getopts "s:h:" opt
do
   case "$opt" in
      s ) skip="$OPTARG" ;;
      ? ) helpFunction ;;
   esac
done

if [ -z "$skip" ]
then
   echo "Some parameters are empty";
   helpFunction;
fi

if "$skip"
then
   echo "skip docker check"
else
   # Version output
	verline=$(docker logs docker-deploy_managerapi_1 | grep -E "^Version : [A-Za-z0-9\-\_\.]+")
	if [ -z "$verline" ];then
	    echo "no Version output"
	    exit 1
	fi
	# Version output
	hashline=$(docker logs docker-deploy_managerapi_1 | grep -E "^GitHash : [A-Za-z0-9\-\_\.]+")
	if [ -z "$hashline" ];then
	    echo "no GitHash output"
	    exit 1
	fi
fi


# web page
curl http://127.0.0.1:9000
code=$(curl -k -i -m 20 -o /dev/null -s -w %{http_code} http://127.0.0.1:9000)
if [ ! $code -eq 200 ]; then
    echo "failed: failed to custom port"
    exit 1
fi

# login
resp=$(curl http://127.0.0.1:9000/apisix/admin/user/login -X POST -H "Content-Type:application/json" -d '{"username":"admin", "password": "admin"}')
token=$(echo "${resp}" | sed 's/{/\n/g' | sed 's/,/\n/g' | grep "token" | sed 's/:/\n/g' | sed '1d' | sed 's/}//g'  | sed 's/"//g')
if [ -z "${token}" ]; then
    echo "login failed"
    exit 1
fi

# plugin orchestration
echo $token
# need apisix server now, temporarily disable this test
# code=$(curl -k -i -m 20 -o /dev/null -s -w %{http_code} http://127.0.0.1:9000/apisix/admin/routes/1 -X PUT -i -H "Authorization: $token" -d '{"id":"1","name":"route1","uri":"/index.html","upstream":{"type":"roundrobin","nodes":[{"host":"www.test.com","port":80,"weight":1}]},"script":{"rule":{"root":"451106f8-560c-43a4-acf2-2a6ed0ea57b8","451106f8-560c-43a4-acf2-2a6ed0ea57b8":[["code==403","b93d622c-92ef-48b4-b6bb-57e1ce893ee3"],["","988ef5c2-c896-4606-a666-3d4cbe24a731"]]},"conf":{"451106f8-560c-43a4-acf2-2a6ed0ea57b8":{"name":"uri-blocker","conf":{"block_rules":["root.exe","root.m+"],"rejected_code":403}},"988ef5c2-c896-4606-a666-3d4cbe24a731":{"name":"kafka-logger","conf":{"batch_max_size":1000,"broker_list":{},"buffer_duration":60,"inactive_timeout":5,"include_req_body":false,"kafka_topic":"1","key":"2","max_retry_count":0,"name":"kafkalogger","retry_delay":1,"timeout":3}},"b93d622c-92ef-48b4-b6bb-57e1ce893ee3":{"name":"fault-injection","conf":{"abort":{"body":"200","http_status":300},"delay":{"duration":500}}}},"chart":{}}}')
# if [ ! $code -eq 200 ]; then
#     echo "failed to create route"
#     exit 1
# fi
