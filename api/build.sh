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
set -e

VERSION=$(cat ./api/VERSION)
GITHASH=$(cat ./.githash 2> /dev/null || HASH="ref: HEAD"; while [[ $HASH == ref\:* ]]; do HASH="$(cat ".git/$(echo $HASH | cut -d \  -f 2)")"; done; echo ${HASH:0:7})

GOLDFLAGS="-X github.com/apache/apisix-dashboard/api/internal/utils.version=${VERSION} -X github.com/apache/apisix-dashboard/api/internal/utils.gitHash=${GITHASH}"

# Enter dry-run mode
if [ "$1" == "--dry-run" ]; then
    cd ./api && go run -ldflags "${GOLDFLAGS}" ./main.go
    exit 0
fi

set -x
export ENV=local
pwd=$(pwd)

rm -rf output && mkdir -p output/conf && mkdir -p output/dag-to-lua

# get dag-to-lua lib
if [[ ! -f "dag-to-lua-1.1/lib/dag-to-lua.lua" ]]; then
    curl -Lso /tmp/v1.1.tar.gz https://github.com/api7/dag-to-lua/archive/v1.1.tar.gz
    tar -zxvf /tmp/v1.1.tar.gz -C /tmp
    cp -r /tmp/dag-to-lua-1.1/lib/* ./output/dag-to-lua
fi

# build
cd ./api && go build -o ../output/manager-api -ldflags "${GOLDFLAGS}" ./main.go && cd ..

cp ./api/conf/schema.json ./output/conf/schema.json
cp ./api/conf/customize_schema.json ./output/conf/customize_schema.json
cp ./api/conf/conf*.yaml ./output/conf/

echo "Build the Manager API successfully"
