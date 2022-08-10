#!/usr/bin/env bats

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

# !!!NOTICE!!!
# Try to run cli tests (such as virtual machines) in an isolated environment.
# This test will register the systemd service and run the test. Although
# environmental cleaning will eventually be running, it may still pollute
# your environment.
#
# For developers who have not modified the CLI part of the program, you do
# not need to run this test locally, but rely on the test results of GitHub
# workflow. If you have any improvement suggestions for this test case, you
# are welcome to submit a PR for modification.

VERSION=$(cat ./VERSION)
KERNEL=$(uname -s)
CONF_FILE="/usr/local/apisix-dashboard/conf/conf.yaml"
APISIX_PROFILE_CONF_FILE="/usr/local/apisix-dashboard/conf/conf-test.yaml"
LOG_FILE="/usr/local/apisix-dashboard/logs/error.log"
ACCESS_LOG_FILE="/usr/local/apisix-dashboard/logs/access.log"
SERVICE_NAME="apisix-dashboard"

if [[ -f ../.githash ]]; then
  GITHASH=$(cat ../.githash)
  if [[ ! $GITHASH =~ ^[a-z0-9]{7}$ ]]; then
    echo "failed: verify .githash content failed"
    exit 1
  fi
else
  GITHASH=$(HASH="ref: HEAD"; while [[ $HASH == ref\:* ]]; do HASH="$(cat "../.git/$(echo $HASH | cut -d \  -f 2)")"; done; echo ${HASH:0:7})
fi

recover_conf() {
  run cp -rf ./conf/conf.yaml ${CONF_FILE}
  run cp -rf ./conf/conf.yaml ${APISIX_PROFILE_CONF_FILE}
  [ "$status" -eq 0 ]
}
check_logfile() {
  [ -f $LOG_FILE ]
}
clean_logfile() {
  echo > $LOG_FILE
}

recover_service_file() {
  run cp -f ./service/apisix-dashboard.service /usr/lib/systemd/system/${SERVICE_NAME}.service
  run systemctl daemon-reload
  [ "$status" -eq 0 ]
}

start_dashboard() {
  run systemctl start ${SERVICE_NAME}
  [ "$status" -eq 0 ]
  sleep $1
}

stop_dashboard() {
  run systemctl stop ${SERVICE_NAME}
  [ "$status" -eq 0 ]
  sleep $1
}

### Test Case
#pre
@test "Build and Deploy APISIX Dashboard Manager API" {
  run go build -o ./manager-api -ldflags "-X github.com/apache/apisix-dashboard/api/internal/utils.version=${VERSION} -X github.com/apache/apisix-dashboard/api/internal/utils.gitHash=${GITHASH}" ./main.go
  [ "$status" -eq 0 ]

  # prepare service files
  mkdir -p /usr/local/apisix-dashboard/conf /usr/local/apisix-dashboard/logs
  cp ./conf/* /usr/local/apisix-dashboard/conf
  cp ./manager-api /usr/local/apisix-dashboard

  # create systemd service
  cp ./service/apisix-dashboard.service /usr/lib/systemd/system/${SERVICE_NAME}.service
  run systemctl daemon-reload
  [ "$status" -eq 0 ]
}

#1
@test "Check warn log level" {
  start_dashboard 3

  stop_dashboard 6

  check_logfile

  [ $(grep -c "INFO" "${LOG_FILE}") -eq '0' ]

  clean_logfile
}

#2
@test "Check info log level and signal" {
  if [[ $KERNEL = "Darwin" ]]; then
    sed -i "" 's/level: warn/level: info/' ${CONF_FILE}
  else
    sed -i 's/level: warn/level: info/' ${CONF_FILE}
  fi

  start_dashboard 3

  stop_dashboard 6

  check_logfile

  [ $(grep -c "server receive terminated" "${LOG_FILE}") -eq '1' ]

  clean_logfile
}

#3
@test "Check start info" {
  LOGLEVEL=$(cat "$CONF_FILE" | awk '$1=="level:"{print $2}')
  HOST=$(cat "$CONF_FILE" | awk '$1=="host:"{print $2}')
  PORT=$(cat "$CONF_FILE" | awk '$1=="port:"{print $2}')
  start_dashboard 3

  run systemctl status ${SERVICE_NAME}

  [ $(echo "$output" | grep -c "The manager-api is running successfully\!") -eq '1' ]
  [ $(echo "$output" | grep -c -w "${VERSION}") -eq '1' ]
  [ $(echo "$output" | grep -c "${GITHASH}") -eq '1' ]
  [ $(echo "$output" | grep -c "${LOGLEVEL}") -eq '1' ]
  [ $(echo "$output" | grep -c "${HOST}:${PORT}") -eq '1' ]

  stop_dashboard 6
}

#4
@test "Check version sub-command" {
  run /usr/local/apisix-dashboard/manager-api version

  [ $(echo "$output" | grep -c "$VERSION") -eq '1' ]
  [ $(echo "$output" | grep -c "$GITHASH") -eq '1' ]
}

#5
@test "Check static file server" {
  # create html directory
  mkdir -p /usr/local/apisix-dashboard/html
  echo "hi~" >> /usr/local/apisix-dashboard/html/index.html

  # start Manager API
  start_dashboard 3

  # request index page
  result=$(curl "http://127.0.0.1:9000")
  [ "$result" = "hi~" ]

  stop_dashboard 6

  recover_conf
}

#6
@test "Check invalid etcd endpoint" {
  recover_conf

  if [[ $KERNEL = "Darwin" ]]; then
    sed -i "" 's/127.0.0.1:2379/0.0.0.0:0/' ${CONF_FILE}
  else
    sed -i 's/127.0.0.1:2379/0.0.0.0:0/' ${CONF_FILE}
  fi

  start_dashboard 6

  run journalctl -u ${SERVICE_NAME}.service -n 30

  [ $(echo "$output" | grep -c "Error while dialing dial tcp") -eq '1' ]

  stop_dashboard 6
}

#7
@test "Check assess log" {
  recover_conf

  start_dashboard 3

  run curl http://127.0.0.1:9000/apisix/admin/user/login -H "Content-Type: application/json" -d '{"username":"admin", "password": "admin"}'
  [ "$status" -eq 0 ]

  stop_dashboard 6

  [ $(grep -c "/apisix/admin/user/login" "${ACCESS_LOG_FILE}") -ne '0' ]

  # check logging middleware
  [ $(grep -c "filter/logging.go" "${ACCESS_LOG_FILE}") -ne '0' ]
}

#8
@test "Check ip allow list" {
  recover_conf

  if [[ $KERNEL = "Darwin" ]]; then
    sed -i "" 's@- 127.0.0.1 @- 10.0.0.1 @' ${CONF_FILE}
  else
    sed -i 's@- 127.0.0.1 @- 10.0.0.1 @' ${CONF_FILE}
  fi

  start_dashboard 3

  run curl -k -i -m 20 -o /dev/null -s -w %{http_code} "http://127.0.0.1:9000"
  [ "$output" -eq 403 ]

  stop_dashboard 6
}

#9
@test "Check HTTPS server" {
  recover_conf

  if [[ $KERNEL = "Darwin" ]]; then
    sed -i "" 's@# ssl:@ssl:@' ${CONF_FILE}
    sed -i "" 's@#   port: 9001@  port: 9001@' ${CONF_FILE}
    sed -i "" "s@#   cert: \"/tmp/cert/example.crt\"@  cert: \"$(pwd)/test/certs/test2.crt\"@" ${CONF_FILE}
    sed -i "" "s@#   key:  \"/tmp/cert/example.key\"@  cert: \"$(pwd)/test/certs/test2.key\"@" ${CONF_FILE}
  else
    sed -i 's@# ssl:@ssl:@' ${CONF_FILE}
    sed -i 's@#   port: 9001@  port: 9001@' ${CONF_FILE}
    sed -i "s@#   cert: \"/tmp/cert/example.crt\"@  cert: \"$(pwd)/test/certs/test2.crt\"@" ${CONF_FILE}
    sed -i "s@#   key:  \"/tmp/cert/example.key\"@  key: \"$(pwd)/test/certs/test2.key\"@" ${CONF_FILE}
  fi

  start_dashboard 3

  run curl -k -i -m 20 -o /dev/null -s -w %{http_code} --resolve 'www.test2.com:9001:127.0.0.1' "https://www.test2.com:9001/apisix/admin/tool/version"
  [ "$output" -eq 200 ]

  stop_dashboard 6
}

#10
@test "Check etcd basic auth" {
  recover_conf

  # add root user
  curl -L http://localhost:2379/v3/auth/user/add -X POST -d '{"name": "root", "password": "root"}'

  # add root role
  curl -L http://localhost:2379/v3/auth/role/add -d '{"name": "root"}'

  # grant root role to root user
  curl -L http://localhost:2379/v3/auth/user/grant -d '{"user": "root", "role": "root"}'

  # enable auth
  curl -L http://localhost:2379/v3/auth/enable -d '{}'

  start_dashboard 3

  [ $(grep -c "etcdserver: user name is empty" ${LOG_FILE}) -ne '0' ]

  stop_dashboard 6

  # modify etcd auth config
  if [[ $KERNEL = "Darwin" ]]; then
    sed -i "" '1,$s/# username: "root"    # ignore etcd username if not enable etcd auth/username: "root"/g' ${CONF_FILE}
    sed -i "" '1,$s/# password: "123456"  # ignore etcd password if not enable etcd auth/password: "root"/g' ${CONF_FILE}
  else
    sed -i '1,$s/# username: "root"    # ignore etcd username if not enable etcd auth/username: "root"/g' ${CONF_FILE}
    sed -i '1,$s/# password: "123456"  # ignore etcd password if not enable etcd auth/password: "root"/g' ${CONF_FILE}
  fi

  start_dashboard 3

  # validate process is right by requesting login api
  run curl http://127.0.0.1:9000/apisix/admin/user/login -H "Content-Type: application/json" -d '{"username":"admin", "password": "admin"}'
  token=$(echo "$output" | sed 's/{/\n/g' | sed 's/,/\n/g' | grep "token" | sed 's/:/\n/g' | sed '1d' | sed 's/}//g'  | sed 's/"//g')

  [ -n "${token}" ]

  # more validation to make sure it's ok to access etcd
  run curl -ig -XPUT http://127.0.0.1:9000/apisix/admin/consumers -i -H "Content-Type: application/json" -H "Authorization: $token" -d '{"username":"etcd_basic_auth_test"}'
  respCode=$(echo "$output" | sed 's/{/\n/g'| sed 's/,/\n/g' | grep "code" | sed 's/:/\n/g' | sed '1d')
  respMessage=$(echo "$output" | sed 's/{/\n/g'| sed 's/,/\n/g' | grep "message" | sed 's/:/\n/g' | sed '1d')

  [ "$respCode" = "0" ]
  [ "$respMessage" = "\"\"" ]

  run curl "http://127.0.0.1:9000/apisix/admin/tool/version"

  [ $(echo "$output" | grep -c "${VERSION}") -eq '1' ]
  [ $(echo "$output" | grep -c "${GITHASH}") -eq '1' ]

  check_logfile

  stop_dashboard 6

  # disable etcd basic auth
  run curl -L http://localhost:2379/v3/auth/authenticate -X POST -d '{"name": "root", "password": "root"}'
  etcd_token=$(echo "$output" |grep -oE "token\".*\"(.*)\""|awk -F[:\"] '{print $4}')
  [ -n "${etcd_token}" ]

  run curl -L http://localhost:2379/v3/auth/disable -H "Authorization: ${etcd_token}" -X POST -d ''
  [ "$status" -eq 0 ]
}

#11
@test "Check etcd prefix" {
  recover_conf

  start_dashboard 3

  run curl http://127.0.0.1:9000/apisix/admin/user/login -H "Content-Type: application/json" -d '{"username":"admin", "password": "admin"}'
  [ "$status" -eq 0 ]

  token=$(echo "$output" | sed 's/{/\n/g' | sed 's/,/\n/g' | grep "token" | sed 's/:/\n/g' | sed '1d' | sed 's/}//g'  | sed 's/"//g')
  [ -n "${token}" ]

  prefix="/apisix"
  key_base64=$(echo -n $prefix/consumers/etcd_prefix_test | base64)

  run curl -ig -XPUT http://127.0.0.1:9000/apisix/admin/consumers -i -H "Content-Type: application/json" -H "Authorization: $token" -d '{"username":"etcd_prefix_test"}'
  [ "$status" -eq 0 ]

  run curl -L http://localhost:2379/v3/kv/range -X POST -d '{"key": "'"${key_base64}"'"}'
  [ "$status" -eq 0 ]

  count=$(echo "$output" | grep -oE "count.*([0-9]+)" | awk -F\" '{print $3}')
  [ "$count" ]
  [ "$count" -eq 1 ]

  stop_dashboard 6

  recover_conf

  # modify etcd prefix config to /apisix-test
  if [[ $KERNEL = "Darwin" ]]; then
    sed -i "" '1,$s/# prefix: \/apisix.*/prefix: \/apisix-test/g' ${CONF_FILE}
  else
    sed -i '1,$s/# prefix: \/apisix.*/prefix: \/apisix-test/g' ${CONF_FILE}
  fi

  start_dashboard 3

  run curl http://127.0.0.1:9000/apisix/admin/user/login -H "Content-Type: application/json" -d '{"username":"admin", "password": "admin"}'
  [ "$status" -eq 0 ]

  token=$(echo "$output" | sed 's/{/\n/g' | sed 's/,/\n/g' | grep "token" | sed 's/:/\n/g' | sed '1d' | sed 's/}//g'  | sed 's/"//g')
  [ -n "${token}" ]

  prefix="/apisix-test"
  key_base64=$(echo -n $prefix/consumers/etcd_prefix_test | base64)

  run curl -ig -XPUT http://127.0.0.1:9000/apisix/admin/consumers -i -H "Content-Type: application/json" -H "Authorization: $token" -d '{"username":"etcd_prefix_test"}'
  [ "$status" -eq 0 ]

  run curl -L http://localhost:2379/v3/kv/range -X POST -d '{"key": "'"${key_base64}"'"}'
  [ "$status" -eq 0 ]

  count=$(echo "$output" | grep -oE "count.*([0-9]+)" | awk -F\" '{print $3}')
  [ "$count" ]
  [ "$count" -eq 1 ]

  stop_dashboard 6
}

#12
@test "Check etcd mTLS" {
  recover_conf

  run ./etcd-v3.4.14-linux-amd64/etcd --name infra0 --data-dir infra0 \
        --client-cert-auth --trusted-ca-file=$(pwd)/test/certs/mtls_ca.pem --cert-file=$(pwd)/test/certs/mtls_server.pem --key-file=$(pwd)/test/certs/mtls_server-key.pem \
        --advertise-client-urls https://127.0.0.1:3379 --listen-client-urls https://127.0.0.1:3379 --listen-peer-urls http://127.0.0.1:3380 &

  if [[ $KERNEL = "Darwin" ]]; then
    sed -i "" "s@key_file: \"\"@key_file: \"$(pwd)/test/certs/mtls_client-key.pem\"@g" ${CONF_FILE}
    sed -i "" "s@cert_file: \"\"@cert_file: \"$(pwd)/test/certs/mtls_client.pem\"@g" ${CONF_FILE}
    sed -i "" "s@ca_file: \"\"@ca_file: \"$(pwd)/test/certs/mtls_ca.pem\"@g" ${CONF_FILE}
    sed -i "" 's/127.0.0.1:2379/127.0.0.1:3379/' ${CONF_FILE}
  else
    sed -i "s@key_file: \"\"@key_file: \"$(pwd)/test/certs/mtls_client-key.pem\"@g" ${CONF_FILE}
    sed -i "s@cert_file: \"\"@cert_file: \"$(pwd)/test/certs/mtls_client.pem\"@g" ${CONF_FILE}
    sed -i "s@ca_file: \"\"@ca_file: \"$(pwd)/test/certs/mtls_ca.pem\"@g" ${CONF_FILE}
    sed -i 's/127.0.0.1:2379/127.0.0.1:3379/' ${CONF_FILE}
  fi

  start_dashboard 3

  run curl http://127.0.0.1:9000/apisix/admin/user/login -H "Content-Type: application/json" -d '{"username":"admin", "password": "admin"}'
  [ "$status" -eq 0 ]

  token=$(echo "$output" | sed 's/{/\n/g' | sed 's/,/\n/g' | grep "token" | sed 's/:/\n/g' | sed '1d' | sed 's/}//g'  | sed 's/"//g')
  [ -n "${token}" ]

  run curl -ig -XPUT http://127.0.0.1:9000/apisix/admin/consumers -i -H "Content-Type: application/json" -H "Authorization: $token" -d '{"username":"etcd_mtls_test"}'
  respCode=$(echo "$output" | sed 's/{/\n/g'| sed 's/,/\n/g' | grep "code" | sed 's/:/\n/g' | sed '1d')
  respMessage=$(echo "$output" | sed 's/{/\n/g'| sed 's/,/\n/g' | grep "message" | sed 's/:/\n/g' | sed '1d')

  [ "$respCode" = "0" ]
  [ "$respMessage" = "\"\"" ]

  stop_dashboard 6
}

#13
@test "Check etcd bad data" {
  recover_conf

  run ./etcd-v3.4.14-linux-amd64/etcdctl put /apisix/routes/unique1 "{\"id\":}"
  [ "$status" -eq 0 ]
  sleep 2

  start_dashboard 3

  run journalctl -u ${SERVICE_NAME}.service -n 30

  [ $(echo "$output" | grep -c "Error occurred while initializing logical store:  /apisix/routes") -eq '1' ]
  [ $(echo "$output" | grep -c "Error: json unmarshal failed") -eq '1' ]

  run ./etcd-v3.4.14-linux-amd64/etcdctl del /apisix/routes/unique1
  [ "$status" -eq 0 ]

  stop_dashboard 6
}


#14
@test "Check APISIX_PROFILE" {
  recover_conf

  start_dashboard 3

  run journalctl -u ${SERVICE_NAME}.service -n 30
  [ $(echo "$output" | grep -c "conf.yaml") -eq '1' ]

  stop_dashboard 3

  sed -i 's#-c /usr/local/apisix-dashboard/conf/conf.yaml##g' /usr/lib/systemd/system/${SERVICE_NAME}.service
  sed -i '$a\Environment=APISIX_PROFILE=test' /usr/lib/systemd/system/${SERVICE_NAME}.service
  run systemctl daemon-reload

  start_dashboard 3

  run journalctl -u ${SERVICE_NAME}.service -n 30
  [ $(echo "$output" | grep -c "conf-test.yaml") -eq '1' ]

  stop_dashboard 3

  recover_service_file
}

#15
@test "Check Security configuration" {
  recover_conf

  start_dashboard 3

  # check response header without custom header
  run curl -i http://127.0.0.1:9000

  [ $(echo "$output" | grep -c "X-Frame-Options: deny") -eq '1' ]

  stop_dashboard 6

  sed -i 's@# security:@security:@' ${CONF_FILE}
  sed -i 's@#   x_frame_options: "deny"@  x_frame_options: "test"@' ${CONF_FILE}

  start_dashboard 3

  # check response header with custom header
  run curl -i http://127.0.0.1:9000

[ $(echo "$output" | grep -c "X-Frame-Options: test") -eq '1' ]

  stop_dashboard 6
}

#post
@test "Clean test environment" {
  # kill etcd
  pkill -f etcd

  # stop dashboard service
  stop_dashboard 0

  # clean configure and log files
  rm -rf /usr/local/apisix-dashboard
  rm /usr/lib/systemd/system/${SERVICE_NAME}.service

  # reload systemd services
  run systemctl daemon-reload
  [ "$status" -eq 0 ]
}
