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

FROM golang:1.19-alpine AS build-env

WORKDIR /workspace

COPY ./ ./

RUN mkdir -p /go/apisix-dashboard/config \
    && go build -o /go/apisix-dashboard/apisix-dashboard ./api \
    && mv ./api/config/* /go/apisix-dashboard/config/ \
    && wget https://github.com/api7/dag-to-lua/archive/v1.1.tar.gz && tar -zxvf v1.1.tar.gz \
    && mkdir -p /go/apisix-dashboard/dag-to-lua && mv ./dag-to-lua-1.1/lib/* /go/apisix-dashboard/dag-to-lua/


FROM alpine:3.14

WORKDIR /go/apisix-dashboard

RUN mkdir logs

EXPOSE 9000

ENTRYPOINT ["/go/apisix-dashboard/apisix-dashboard"]
CMD ["-test.coverprofile=./testdata/integrationcover.out"]

COPY --from=build-env /go/apisix-dashboard/ /go/apisix-dashboard/