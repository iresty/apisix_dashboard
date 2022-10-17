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

#Executes commands from the intended directory "/api/test/docker".
cd "$(dirname "$0")"

main() {
  #welcome message & check machine configuration
  welcome
  check

  if [ $# -eq 0 ]; then
    help
    give_up "\nNo arguments provided."
  fi

  UPFLAG="-d"
  UP=0
  for arg in "$@"; do
    case "$arg" in
    up | -u | -U)
      UP=1
      ;;
    --non-detach)
      UPFLAG=""
      ;;
    down | -d | -D)
      down
      ;;
    build | -b | -B)
      build
      ;;
    help | -h | --help | -H)
      help
      ;;
    *)
      echo -e "Invalid Argument.\n"
      help
      give_up
      ;;
    esac
  done

  if [ $UP -eq 1 ]; then
    up "$UPFLAG"
  fi
  echo "Execution complete."
}

check() {
  set -e

  if ! docker --version &>/dev/null; then
    give_up "The script depends on docker. Please proceed after the installation."
  fi

  if ! docker-compose --version &>/dev/null; then
    give_up "The script depends on docker-compose. Please proceed after the installation."
  fi

  if ! sed --version &>/dev/null; then
    give_up "The script depends on sed (GNU Stream Editor). Please proceed after the installation."
  fi
}

CONF_PATH=$(pwd | sed -e 's/test\/docker/conf\//g')
YAML_FILE="${CONF_PATH}conf.yaml"
BACKUP_FILE="${CONF_PATH}.backup.yaml"

up() {
  set -e
  #creating backup of current config
  if [ ! -f "$BACKUP_FILE" ]; then
    cp "$YAML_FILE" "$BACKUP_FILE"

    #modifying config
    sed -i 's/127.0.0.1:2379/172.16.238.10:2379/' "$YAML_FILE"
    sed -i 's/127.0.0.1:9180/172.16.238.30:9180/' "$YAML_FILE"
    sed -i 's@127.0.0.1@0.0.0.0/0@' "$YAML_FILE"
    sed -i '/172.16.238.10:2379/a\      - 172.16.238.11:2379' "$YAML_FILE"
    sed -i '/172.16.238.10:2379/a\      - 172.16.238.12:2379' "$YAML_FILE"
    sed -i 's@# - dubbo-proxy@- dubbo-proxy@' "$YAML_FILE"

  fi

  if [ ! -d "grpc-server-example" ]; then
    wget https://github.com/api7/grpc_server_example/archive/refs/tags/20210819.tar.gz
    mkdir grpc-server-example
    tar -xzvf 20210819.tar.gz -C grpc-server-example
    docker build -t grpc_server_example:latest ./grpc-server-example/grpc_server_example-20210819
  fi

  #Spinning up services
  if [ -z "$1" ]; then
    docker-compose up
  else
    docker-compose up "$1"
  fi

  ../shell/wait_for_services.sh
}

down() {
  #restoring old configuration.
  mv "$BACKUP_FILE" "$YAML_FILE"

  docker-compose down -v
}

build() {
  docker-compose build
}

help() {
  echo -e 'Usage: ./setup.sh [OPTIONS]...\n'
  echo "Single argument is mandatory."
  echo "Supported Command Arguments"
  echo -e ' up,\t-u\t spins up all the services in detach mode.'
  echo -e ' --non-detach\t pass the flag with "up / -u" to spin up in non-detach mode.'
  echo -e ' down,\t-d\t stop containers & delete containers, networks, volumes.'
  echo -e ' build,\t-b\t rebuild all the images from dockerfile.'
  echo -e ' help,\t-h\t info about how to use this script.'

  echo -e '\nFull documentation at <https://github.com/apache/apisix-dashboard/blob/master/docs/en/latest/back-end-tests.md>'
}

give_up() {
  echo -e "$1"
  exit 1
}

welcome() {
  echo "WELCOME TO THE SETUP SCRIPT"
  echo "It spins up a fleet of services required for local development."
  echo "Capable of making the necessary changes in yaml file & revert back to original state."
  echo "Services: "
  echo "1. Three etcd nodes."
  echo "2. Two apisix nodes."
  echo -e "3. Single node of\n\t manager-api \n\t skywalking \n\t upstream-node \n\t upstream-grpc \n\t upstream-echo"
  echo -e "=====================================================\n"
}

main "$@"
