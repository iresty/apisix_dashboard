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

/**
 * Schema: https://github.com/apache/apisix/blob/master/apisix/schema_def.lua#L335
 */
declare namespace UpstreamComponent {
  type ActiveCheck = {};

  type PassiveCheck = {};

  type TLS = {
    client_cert: string;
    client_key: string;
  };

  type Node = {
    host: string;
    port: number;
    weight: number;
    priority?: number;
  };

  type SubmitNode = Record<string, number>;

  type Timeout = {
    connect: number;
    send: number;
    read: number;
  };

  type ResponseData = {
    nodes?: SubmitNode | Node[];
    retries?: number;
    timeout?: Timeout;
    tls?: TLS;
    type?: string;
    checks?: {
      active?: ActiveCheck;
      passive?: PassiveCheck;
    };
    hash_on?: string;
    key?: string;
    scheme?: string;
    discovery_type?: string;
    pass_host?: string;
    upstream_host?: string;
    name?: string;
    desc?: string;
    service_name_type?: string;
    service_name?: string;
    id?: string;
    upstream_id?: string;
    upstream_type?: string;

    // NOTE: custom field
    custom?: Record<string, any>;
  };
}
