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
declare namespace RouteModule {
  type Operator = '==' | '~=' | '>' | '<' | '~~' | '~*' | 'IN' | 'HAS' | '!';

  type VarPosition = 'arg' | 'post_arg' | 'http' | 'cookie' | 'buildin';

  type RequestProtocol = 'https' | 'http' | 'websocket';

  type RequestMode = 'Normal' | 'RawData';

  type BaseData = {
    id?: number;
    name: string;
    desc: string;
    uris: string[];
    hosts: string[];
    status: number;
  };

  type Step3Data = {
    plugins: PluginPage.PluginData;
    //  TEMP
    script: any;
    plugin_config_id?: string;
  };

  type UpstreamHost = {
    host: string;
    port: number;
    weight: number;
  };

  type UpstreamHeader = {
    header_name: string;
    header_value: string;
    key: string;
  };

  type ModalType = 'CREATE' | 'EDIT';

  type Redirect = {
    ret_code?: number;
    uri?: string;
    http_to_https?: boolean;
  };

  // Request Body or Response Data for API
  type Body = {
    id?: number;
    status: number;
    name: string;
    labels: Record<string, string>;
    desc: string;
    priority?: number;
    methods: HttpMethod[];
    uri?: string;
    uris: string[];
    host?: string;
    hosts?: string[];
    remote_addr?: string;
    remote_addrs?: string[];
    upstream: UpstreamComponent.ResponseData;
    vars: [string, Operator, string | any[]][];
    upstream_path?: {
      type?: string;
      from?: string;
      to: string;
    };
    upstream_id?: string;
    plugins: Record<string, any>;
    plugin_config_id?: string;
    script: Record<string, any>;
    url?: string;
    enable_websocket?: boolean;
    service_id?: string;
  };

  type MatchingRule = {
    position: VarPosition;
    name: string;
    operator: Operator;
    value: string | any[];
    key: string;
  };

  type Step1PassProps = {
    form: FormInstance;
    upstreamForm?: FormInstance;
    advancedMatchingRules: MatchingRule[];
    disabled?: boolean;
    isEdit?: boolean;
    onChange?: (data: { action: string; data: T }) => void;
  };

  type Form1Data = {
    name: string;
    desc: string;
    custom_version_label: string;
    custom_normal_labels: string[];
    priority: number;
    websocket: boolean;
    hosts: string[];
    uris: string[];
    remote_addrs: string[];
    methods: HttpMethod[];
    redirectOption: 'forceHttps' | 'customRedirect' | 'disabled';
    redirectURI?: string;
    ret_code?: number;
    status: number;
    enable_websocket?: boolean;
    service_id: string;
    proxyRewrite: ProxyRewrite;
    URIRewriteType: number;
    hostRewriteType: number;
  };
  type Kvobject = {
    key: string;
    value: string;
  };
  type ProxyRewrite = {
    scheme?: 'keep' | 'http' | 'https';
    uri?: string;
    regex_uri?: string[];
    host?: string;
    method?: string;
    kvHeaders?: Kvobject[];
    headers?: Record<string, string>;
  };

  type AdvancedMatchingRules = {
    advancedMatchingRules: MatchingRule[];
  };

  type Step2PassProps = {
    form: FormInstance;
    disabled?: boolean;
    upstreamRef: any;
    hasServiceId: boolean;
  };

  type RequestData = {
    form1Data: Form1Data;
    form2Data: UpstreamComponent.ResponseData;
    step3Data: Step3Data;
    advancedMatchingRules: MatchingRule[];
  };

  type RequestBody = {
    name?: string;
    desc?: string;
    uri: string;
    host?: string;
    hosts?: string[];
    remote_addr?: string;
    remote_addrs?: string[];
    methods?: HttpMethod[];
    priority?: number;
    vars?: [string, Operator, string][];
    filter_func?: string;
    plugins?: Record<string, any>;
    script?: Record<string, any>;
    // TODO:
    upstream?: any;
    upstream_id?: string;
    service_id?: string;
    service_protocol?: 'grpc' | 'http';
  };

  type ResponseBody = {
    id: string;
    methods: HttpMethod[];
    name: string;
    remote_addrs: string[];
    script: any;
    desc?: string;
    labels: Record<string, string>;
    upstream: {
      checks: UpstreamModule.HealthCheck;
      create_time: number;
      id: string;
      nodes: {
        port: number;
      }[];
      timeout: UpstreamModule.Timeout;
      type: UpstreamModule.Type;
    };
    uri: string;
    uris?: string[];
    host: string;
    hosts?: string[];
    create_time: number;
    update_time: number;
    status: number;
  };

  type RouteStatus = 0 | 1;

  // TODO： grpc and websocket
  type debugRequest = {
    url: string;
    request_protocol: RequestProtocol | 'grpc';
    method: string;
    body_params?: any;
    header_params?: any;
  };

  type debugResponse = {
    code: number;
    message: string;
    data: any;
    header: Record<string, string[]>;
  };

  type authData = {
    authType: string;
    username?: string;
    password?: string;
    Authorization?: string;
    apikey?: string;
  };
  type debugRequestParamsFormData = {
    check: boolean;
    key: string;
    value: any;
    type?: string;
  };
  type DebugViewProps = {
    form: FormInstance;
    name?: string;
    inputType?: 'param' | 'header';
  };
  type DebugBodyType = 'none' | 'x-www-form-urlencoded' | 'raw input' | 'form-data';
  type DebugDrawProps = {
    visible: boolean;
    data: Record<string, any>;
    onClose: () => void;
  };
}
