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
declare namespace UpstreamModule {
  type Node = Record<string, number | string>;
  type Type = 'roundrobin' | 'chash' | 'ewma';
  type DiscoveryType = 'dns' | 'consul' | 'consul_kv' | 'nacos' | 'eureka';
  type DiscoveryArgs = {
    group_name?: string;
    namespace_id?: string;
  };

  type KeepalivePool = {
    size?: number;
    idle_timeout?: number;
    requests?: number;
  };

  type Timeout = Record<'connect' | 'send' | 'read', number>;

  type HealthCheck = {
    active?: {
      timeout?: number;
      http_path: string;
      host: string;
      healthy: {
        interval: number;
        successes: number;
      };
      unhealthy: {
        interval: number;
        http_failures: number;
      };
      req_headers?: string[];
      port?: number;
    };
    passive?: {
      healthy: {
        http_statuses: number[];
        successes: number;
      };
      unhealthy: {
        http_statuses: number[];
        http_failures: number;
        tcp_failures: number;
      };
    };
  };

  type UpstreamHost = {
    host: string;
    port: number;
    weight: number;
  };

  type K8SDeploymentInfo = {
    namespace: string;
    deploy_name: string;
    service_name: string;
    backend_type: string;
    port: number;
  };

  type RequestBody = {
    id: string;
    upstream_id: string;
    type: Type;
    upstream_type: string;
    discovery_type?: DiscoveryType;
    service_name_type?: 'pass' | 'rewrite';
    service_name?: string;
    discovery_args?: DiscoveryArgs;
    nodes?: UpstreamComponent.SubmitNode;
    hash_on?: 'vars' | 'header' | 'cookie' | 'consumer';
    key?: string;
    checks?: HealthCheck;
    retries?: number;
    retry_timeout?: number;
    enable_websocket?: boolean;
    timeout?: Timeout;
    name?: string;
    desc?: string;
    pass_host?: 'pass' | 'node' | 'rewrite';
    upstream_host: UpstreamHost[];
    keepalive_pool: KeepalivePool;

    // Custom Fields that need to be omitted
    custom?: {};
    submitNodes?: Node[];
  };

  // TODO: typing
  type ResponseBody = Record<string, unknown> & RequestBody;
}
