---
title: FAQ
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

### 1. Vue.js version of the Dashboard

If you need a Vue.js build of the Apache APISIX Dashboard 1.0, use the [master-vue branch](https://github.com/apache/apisix-dashboard/tree/master-vue).

### 2. What are the differences between Dashboard version 2.0 and version 1.5?

The 2.0 version of the dashboard removed MySQL from [version 1.5](https://github.com/apache/apisix-dashboard/tree/backup-1.5-latest) and will operate directly on etcd.

### 3. Etcd compatibility issues

If you are using Apache APISIX below v2.0, be aware that the data from the etcd v2 API is [not compatible](https://etcd.io/docs/v3.4.0/op-guide/v2-migration/) with the data from the v3 API. Apache APISIX Dashboard v2.0 and above uses the etcd v3 API, and apisix 1.5 and below uses the etcd v2 API.

### 4. After modifying the plugin schema or creating a custom plugin in Apache APISIX, why can't I find it on the dashboard?

Since the Dashboard caches the jsonschema data of the plugins in Apache APISIX, you need to synchronize the data in the Dashboard after you create your custom plugins in Apache APISIX, which currently **only supports manual operation**. Please follow the following guide.

1. Confirm that your APISIX is running and has enabled control API (enabled by default and only runs local access)
   Refer to the beginning in:
   [https://apisix.apache.org/docs/apisix/control-api](https://apisix.apache.org/docs/apisix/control-api)

2. Execute the following commands to export jsonchema on your APISIX server (if it is configured for non-local access, it does not need to be executed on your APISIX server, and the access IP and port should be modified accordingly)

```sh
curl 127.0.0.1:9090/v1/schema > schema.json
```

Refer to [https://apisix.apache.org/docs/apisix/control-api#get-v1schema](https://apisix.apache.org/docs/apisix/control-api#get-v1schema)

3. Copy the exported `schema.json` to the `conf` directory in the Dashboard working directory (About working directory, please refer to https://github.com/apache/apisix-dashboard/blob/release/2.9.0/docs/en/latest/deploy.md#working-directory)

4. Restart the Manager API

### 5. How to write API documentation

We use [go-swagger](https://github.com/go-swagger/go-swagger) to generate Swagger 2.0 documents, and then convert them to markdown format so that they can be viewed directly in the github repository. Specific steps are as follows:

1. Write comments according to [Specification](https://goswagger.io/use/spec.html). For details, please refer to the existing example `api/internal/handler/route/route.go` in this project.

2. Use the `go-swagger` tool to generate Swagger 2.0 documents.

```shell
$ swagger generate spec -o ./docs/en/latest/api/api.yaml --scan-models
```

3. Use the `swagger-markdown` tool to convert Swagger 2.0 documents into markdown documents.

```shell
$ swagger-markdown -i ./docs/en/latest/api/api.yaml
```

### 6. How to allow all IPs to access APISIX Dashboard

1. Allow all IPv4 access

By default, the IPv4 range of `127.0.0.0/24` is allowed to access `APISIX Dashboard`. If you want to allow all IPv4 access, then just configure `conf.allow_list` in the configuration file of `conf/conf.yaml` as follows:

```yaml
conf:
  allow_list:
    - 0.0.0.0/0
```

2. Allow all IPv6 access

By default, the IPv6 range of `::1` is allowed to access `APISIX Dashboard`. If you want to allow all IPv6 access, then just configure `conf.allow_list` in the configuration file of `conf/conf.yaml` as follows:

```yaml
conf:
  allow_list:
    - ::/0
```

3. Allow all IP access

If you want to allow all IPs to access `APISIX Dashboard`, you only need to do the following configuration in the configuration file of `conf/conf.yaml`:

```yaml
conf:
  allow_list:
```

Restart `manager-api`, all IPs can access `APISIX Dashboard`.

Note: You can use this method in development and test environment to allow all IPs to access your `APISIX Dashboard` instance, but it is not safe to use it in a production environment. In production environment, please only authorize specific IP addresses or address ranges to access your instance.

### 7. What is the default strategy when import a duplicate route?

Currently we reject import duplicate route, that is to say when you import a route which has the same attributes, all of the `uri` `uris` `host` `hosts` `remote_addr` `remote_addrs` `priority` `vars` and `filter_func`, as the existing route, you will get an error while importing a route from OAS3.0.

### 8. APISIX dashboard add grafana cross-domain problem

Modifying the Grafana configuration:

1. Enable anonymous access:

```shell
# grep 'auth.anonymous' -A 3 defaults.ini
[auth.anonymous]
# enable anonymous access
enabled = true
```

2. Allow access via iframe

```shell
# grep 'allow_embedding' defaults.ini
allow_embedding = true
···
```

### 9. APISIX dashboard configured domain name, the embedded Grafana can't login

If the domain name of the address is configured as HTTPS, the embedded grafana will jump to the login page after logging in. You can refer to this solution:

It's best for Grafana to configure the domain name in the same way. Otherwise there will be problems with address resolution.
