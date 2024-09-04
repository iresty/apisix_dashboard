---
title: API doc of Manager API
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

Manager API directly operates ETCD and provides data management for Apache APISIX, provides APIs for Front-end or other clients.

**License:** [Apache License 2.0](http://www.apache.org/licenses/LICENSE-2.0)

### /apisix/admin/migrate/export

#### GET

##### Summary

Export a config file for migrate.

##### Parameters

None.

##### Responses

A file for download.

### /apisix/admin/migrate/import

##### Summary

Import the config file for restore config.

#### POST

##### Parameters (FORM)

| Name | Located in | Description                             | Required | Schema |
| ---- | ---------- | --------------------------------------- | -------- | ------ |
| mode | body(form) | import mode (return, skip or overwrite) | Yes      | string |
| file | body(form) | file to upload                          | Yes      | string |

##### Responses

| Code  | Description     | Schema                |
| ----- | --------------- | --------------------- |
| 0     | import success  | [ApiError](#ApiError) |
| 20001 | Config conflict | [ApiError](#ApiError) |

### /apisix/admin/check_ssl_cert

#### POST

##### Summary

verify SSL cert and key.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| cert | body       | cert of SSL | Yes      | string |
| key  | body       | key of SSL  | Yes      | string |

##### Responses

| Code    | Description       | Schema                |
| ------- | ----------------- | --------------------- |
| 0       | SSL verify passed | [ApiError](#ApiError) |
| default | unexpected error  | [ApiError](#ApiError) |

### /apisix/admin/check_ssl_exists

#### POST

##### Summary

Check whether the SSL exists.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| cert | body       | cert of SSL | Yes      | string |
| key  | body       | key of SSL  | Yes      | string |

##### Responses

| Code    | Description      | Schema                |
| ------- | ---------------- | --------------------- |
| 0       | SSL exists       | [ApiError](#ApiError) |
| default | unexpected error | [ApiError](#ApiError) |

### /apisix/admin/consumers

#### GET

##### Summary

Return the consumer list according to the specified page number and page size, and can search consumers by username.

##### Parameters

| Name      | Located in | Description          | Required | Schema  |
| --------- | ---------- | -------------------- | -------- | ------- |
| page      | query      | page number          | No       | integer |
| page_size | query      | page size            | No       | integer |
| username  | query      | username of consumer | No       | string  |

##### Responses

| Code    | Description      | Schema                    |
| ------- | ---------------- | ------------------------- |
| 0       | list response    | [ [consumer](#consumer) ] |
| default | unexpected error | [ApiError](#ApiError)     |

### /apisix/admin/notexist/routes

#### GET

##### Summary

Return result of route exists checking by name and exclude id.

##### Parameters

| Name    | Located in | Description                       | Required | Schema |
| ------- | ---------- | --------------------------------- | -------- | ------ |
| name    | query      | name of route                     | No       | string |
| exclude | query      | id of route that exclude checking | No       | string |

##### Responses

| Code    | Description      | Schema                |
| ------- | ---------------- | --------------------- |
| 0       | route not exists | [ApiError](#ApiError) |
| default | unexpected error | [ApiError](#ApiError) |

### /apisix/admin/routes

#### GET

##### Summary

Return the route list according to the specified page number and page size, and can search routes by name and uri.

##### Parameters

| Name      | Located in | Description    | Required | Schema  |
| --------- | ---------- | -------------- | -------- | ------- |
| page      | query      | page number    | No       | integer |
| page_size | query      | page size      | No       | integer |
| name      | query      | name of route  | No       | string  |
| uri       | query      | uri of route   | No       | string  |
| label     | query      | label of route | No       | string  |

##### Responses

| Code    | Description      | Schema                |
| ------- | ---------------- | --------------------- |
| 0       | list response    | [ [route](#route) ]   |
| default | unexpected error | [ApiError](#ApiError) |

### /apisix/admin/services

#### GET

##### Summary

Return the service list according to the specified page number and page size, and can search services by name.

##### Parameters

| Name      | Located in | Description     | Required | Schema  |
| --------- | ---------- | --------------- | -------- | ------- |
| page      | query      | page number     | No       | integer |
| page_size | query      | page size       | No       | integer |
| name      | query      | name of service | No       | string  |

##### Responses

| Code    | Description      | Schema                  |
| ------- | ---------------- | ----------------------- |
| 0       | list response    | [ [service](#service) ] |
| default | unexpected error | [ApiError](#ApiError)   |

### /apisix/admin/ssl

#### GET

##### Summary

Return the SSL list according to the specified page number and page size, and can SSLs search by sni.

##### Parameters

| Name      | Located in | Description | Required | Schema  |
| --------- | ---------- | ----------- | -------- | ------- |
| page      | query      | page number | No       | integer |
| page_size | query      | page size   | No       | integer |
| sni       | query      | sni of SSL  | No       | string  |

##### Responses

| Code    | Description      | Schema                |
| ------- | ---------------- | --------------------- |
| 0       | list response    | [ [ssl](#ssl) ]       |
| default | unexpected error | [ApiError](#ApiError) |

### /apisix/admin/upstreams

#### GET

##### Summary

Return the upstream list according to the specified page number and page size, and can search upstreams by name.

##### Parameters

| Name      | Located in | Description      | Required | Schema  |
| --------- | ---------- | ---------------- | -------- | ------- |
| page      | query      | page number      | No       | integer |
| page_size | query      | page size        | No       | integer |
| name      | query      | name of upstream | No       | string  |

##### Responses

| Code    | Description      | Schema                    |
| ------- | ---------------- | ------------------------- |
| 0       | list response    | [ [upstream](#upstream) ] |
| default | unexpected error | [ApiError](#ApiError)     |

### /apisix/admin/user/login

#### POST

##### Summary

user login.

##### Parameters

| Name     | Located in | Description | Required | Schema |
| -------- | ---------- | ----------- | -------- | ------ |
| username | body       | user name   | Yes      | string |
| password | body       | password    | Yes      | string |

##### Responses

| Code    | Description      | Schema                |
| ------- | ---------------- | --------------------- |
| 0       | login success    | [ApiError](#ApiError) |
| default | unexpected error | [ApiError](#ApiError) |

### Models

#### ApiError

| Name    | Type   | Description      | Required |
| ------- | ------ | ---------------- | -------- |
| code    | long   | response code    | No       |
| message | string | response message | No       |

#### BaseInfo

| Name        | Type   | Description | Required |
| ----------- | ------ | ----------- | -------- |
| create_time | long   |             | No       |
| id          | object |             | No       |
| update_time | long   |             | No       |

#### Consumer

| Name        | Type   | Description | Required |
| ----------- | ------ | ----------- | -------- |
| create_time | long   |             | No       |
| desc        | string |             | No       |
| id          | object |             | No       |
| labels      | object |             | No       |
| plugins     | object |             | No       |
| update_time | long   |             | No       |
| username    | string |             | No       |

#### LoginInput

| Name     | Type   | Description | Required |
| -------- | ------ | ----------- | -------- |
| password | string | password    | No       |
| username | string | user name   | No       |

#### Route

| Name             | Type                        | Description | Required |
| ---------------- | --------------------------- | ----------- | -------- |
| create_time      | long                        |             | No       |
| desc             | string                      |             | No       |
| enable_websocket | boolean                     |             | No       |
| filter_func      | string                      |             | No       |
| host             | string                      |             | No       |
| hosts            | [ string ]                  |             | No       |
| id               | object                      |             | No       |
| labels           | object                      |             | No       |
| methods          | [ string ]                  |             | No       |
| name             | string                      |             | No       |
| plugins          | object                      |             | No       |
| priority         | long                        |             | No       |
| remote_addr      | string                      |             | No       |
| remote_addrs     | [ string ]                  |             | No       |
| script           | object                      |             | No       |
| service_id       | object                      |             | No       |
| service_protocol | string                      |             | No       |
| update_time      | long                        |             | No       |
| upstream         | [UpstreamDef](#UpstreamDef) |             | No       |
| upstream_id      | object                      |             | No       |
| uri              | string                      |             | No       |
| uris             | [ string ]                  |             | No       |
| vars             | object                      |             | No       |

#### SSL

| Name           | Type       | Description | Required |
| -------------- | ---------- | ----------- | -------- |
| cert           | string     |             | No       |
| certs          | [ string ] |             | No       |
| create_time    | long       |             | No       |
| exptime        | long       |             | No       |
| id             | object     |             | No       |
| key            | string     |             | No       |
| keys           | [ string ] |             | No       |
| labels         | object     |             | No       |
| sni            | string     |             | No       |
| snis           | [ string ] |             | No       |
| status         | long       |             | No       |
| update_time    | long       |             | No       |
| validity_end   | long       |             | No       |
| validity_start | long       |             | No       |

#### Service

| Name             | Type                        | Description | Required |
| ---------------- | --------------------------- | ----------- | -------- |
| create_time      | long                        |             | No       |
| desc             | string                      |             | No       |
| enable_websocket | boolean                     |             | No       |
| id               | object                      |             | No       |
| labels           | object                      |             | No       |
| name             | string                      |             | No       |
| plugins          | object                      |             | No       |
| script           | string                      |             | No       |
| update_time      | long                        |             | No       |
| upstream         | [UpstreamDef](#UpstreamDef) |             | No       |
| upstream_id      | object                      |             | No       |

#### Upstream

| Name                | Type   | Description | Required |
| ------------------- | ------ | ----------- | -------- |
| checks              | object |             | No       |
| create_time         | long   |             | No       |
| desc                | string |             | No       |
| hash_on             | string |             | No       |
| id                  | object |             | No       |
| k8s_deployment_info | object |             | No       |
| key                 | string |             | No       |
| labels              | object |             | No       |
| name                | string |             | No       |
| nodes               | object |             | No       |
| pass_host           | string |             | No       |
| retries             | long   |             | No       |
| service_name        | string |             | No       |
| timeout             | object |             | No       |
| type                | string |             | No       |
| update_time         | long   |             | No       |
| upstream_host       | string |             | No       |

#### UpstreamDef

| Name                | Type   | Description | Required |
| ------------------- | ------ | ----------- | -------- |
| checks              | object |             | No       |
| desc                | string |             | No       |
| hash_on             | string |             | No       |
| k8s_deployment_info | object |             | No       |
| key                 | string |             | No       |
| labels              | object |             | No       |
| name                | string |             | No       |
| nodes               | object |             | No       |
| pass_host           | string |             | No       |
| retries             | long   |             | No       |
| service_name        | string |             | No       |
| timeout             | object |             | No       |
| type                | string |             | No       |
| upstream_host       | string |             | No       |
