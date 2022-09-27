---
title: Development Guide
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

The Dashboard contains both `manager-api` and `web` parts, so you need to start the development environment separately.

## Prerequisites

Before development, refer to this [guide](./install.md) to install dependencies.

## Clone the project

```sh
$ git clone -b release/2.13 https://github.com/apache/apisix-dashboard.git
```

## Start developing

```sh
$ cd apisix-dashboard
```

### manager-api

1. Please change the configuration in `api/conf/conf.yaml`.

2. In the root directory, launch development mode.

```sh
$ make api-run
```

3. In the root directory, stop development mode.

```sh
$ make api-stop
```

4. Please refer to the [FAQ](./FAQ.md) about the problem of displaying exception in the dashboard after adding custom plugins or modifying plugin's schema.

5. If writing an back end E2E test, please refer to the [Back End E2E Writing Guide](./back-end-tests.md)

### web

1. Go to the `web` directory.

```sh
$ cd ./web
```

2. Please change the `manager-api` address in the `web/.env` file. If you follow this guidelines, the address may need to be set as below.

> All commands here are for Linux environment, other systems please use the corresponding commands for your platform. You are also welcome to contribute your own methods.

```bash
echo "SERVE_URL_DEV=http://localhost:9000" > web/.env
```

If you don't want to create the file, you can also export the variable.

```bash
export SERVE_URL_DEV=http://localhost:9000
```

3. Launch development mode [(pnpm)](https://pnpm.io/installation)

```shell
# You can also use pnpm as a tool to download dependent packages.
$ pnpm install --unsafe-perm

$ pnpm start

# If some dependent packages are not installed, use the following command. The higher version of NPM forbids automatic installation of peer-to-peer dependencies, but we need them, so conflicts arise.

# In In the .npmrc configuration file, add strict peer dependencies = false, which means that the strict peer dependency mode will be turned off.
$ npm config set strict-peer-dependencies=false

# Change it to true if you want to automatically install peer dependencies.
$ npm config set auto-install-peers=true

```

> If there is an error about gyp during pnpm install, please ignore it and go ahead! If monaco-editor cannot be found during startup, please go to the original address of the file and copy it .public, you can use the command 'find -name monaco-editor'

4. If writing an front end E2E test, please refer to the [Front End E2E Writing Guide](./front-end-e2e.md)
