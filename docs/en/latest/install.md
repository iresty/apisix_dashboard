---
title: Basic Deploy
---

<!--
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
-->

Installing Apache APISIX Dashboard on Linux is easy.
Now, we provide Docker image and RPM installation package.

## Docker {#docker}

We recommend using Docker to run Dashboard:

```shell
docker pull apache/apisix-dashboard
docker run -d --name dashboard \
           -p 9000:9000        \
           -v <CONFIG_FILE>:/usr/local/apisix-dashboard/conf/conf.yaml \
           apache/apisix-dashboard
```

:::note
Please replace `<CONFIG_FILE>` to your configure file path.

Config file example:

https://github.com/apache/apisix-dashboard/blob/master/api/conf/conf.yaml
:::

## RPM {#rpm}

**NOTE:** Only CentOS 7 is supported currently.

### Install

```shell
# 1. install RPM package
sudo yum install -y https://github.com/apache/apisix-dashboard/releases/download/v3.0.1/apisix-dashboard-3.0.1-0.el7.x86_64.rpm
```

### Launch

```shell
# run dashboard in the shell
sudo manager-api -p /usr/local/apisix/dashboard/

# or run dashboard as a service
systemctl start apisix-dashboard
```

Without changing the configuration, visit `http://127.0.0.1:9000` to use the dashboard with GUI, where the default username and password are `admin`.

## Source {#source}

The Dashboard project contains both `manager-api` and `web`, but `web` is _optional_.

The `manager-api` and `web` will be included in this build guide product.

### Prerequisites {#source-prerequisites}

Before using source codes to build, make sure that the following dependencies are installed in your environment.

For `manager-api`:

1. [Golang](https://golang.org/dl/) 1.16+

> Tip: For users in mainland China, you can use the following command to speed up the module downloads.

```sh
$ go env -w GOPROXY=https://goproxy.cn,direct
```

For `web`:

1. [Node.js](https://nodejs.org/en/download/) current LTS (14.x+)
2. [Yarn](https://yarnpkg.com/getting-started/install)

### Download {#source-download}

```shell
git clone -b release/3.0 https://github.com/apache/apisix-dashboard.git && cd apisix-dashboard
```

### Build {#source-build}

```shell
cd apisix-dashboard
make build
```

When the build is complete, the results are stored in the root `output` directory.

Note: `make build` will build `manager-api` and `web`, use the `make help` command to see more commands.

### Launch {#source-launch}

1. After the build is complete and before you start, make sure the following dependencies are installed and running in your environment.

- [etcd](https://etcd.io/docs/v3.4.0/dl-build/) 3.4.0+

2. Check and modify the configuration information in `output/conf/conf.yaml` according to your deployment environment.

3. Launch the Dashboard

```shell
cd ./output

./manager-api
```

4. Without changing the configuration, visit `http://127.0.0.1:9000` to use the dashboard with GUI, where the default username and password are `admin`.

### Service {#source-service}

You will need to handle your own service management when deploying using the source code compilation method. We provide a service file template for operating systems that use the Systemd service manager.

1. Install

```shell
mkdir -p /usr/local/apisix-dashboard
cp -rf ./output/* /usr/local/apisix-dashboard
```

2. Create service unit

Copy the following or use this [**file**](https://github.com/apache/apisix-dashboard/tree/master/service/apisix-dashboard.service) directly, you need to copy it to the `/usr/lib/systemd/system` directory and execute the `systemctl daemon-reload` command.

```shell
# copy service unit
cp ./api/service/apisix-dashboard.service /usr/lib/systemd/system/apisix-dashboard.service
systemctl daemon-reload

# or: If you need to modify the service unit, you can use the following command
echo "[Unit]
Description=apisix-dashboard
Conflicts=apisix-dashboard.service
After=network-online.target

[Service]
WorkingDirectory=/usr/local/apisix-dashboard
ExecStart=/usr/local/apisix-dashboard/manager-api -c /usr/local/apisix-dashboard/conf/conf.yaml" > /usr/lib/systemd/system/apisix-dashboard.service
```

3. Manage service

You can use the following command to manage the service.

```shell
# start apisix-dashboard
systemctl start apisix-dashboard

# stop apisix-dashboard
systemctl stop apisix-dashboard

# check apisix-dashboard status
systemctl status apisix-dashboard
```
