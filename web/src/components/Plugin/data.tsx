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
import React from 'react';

import IconFont from '../IconFont';

export const PLUGIN_ICON_LIST: Record<string, any> = {
  prometheus: <IconFont name="iconPrometheus_software_logo" />,
  skywalking: <IconFont name="iconskywalking" />,
  'jwt-auth': <IconFont name="iconjwt-3" />,
  'authz-keycloak': <IconFont name="iconkeycloak_icon_32px" />,
  'openid-connect': <IconFont name="iconicons8-openid" />,
  'kafka-logger': <IconFont name="iconApache_kafka" />,
};

// This list is used to filter out plugins that cannot be displayed in the plugins list.
export const PLUGIN_FILTER_LIST: Record<string, { list: PluginComponent.ReferPage[] }> = {
  redirect: { list: ['route'] }, // Filter out the redirect plugin on the route page.
  'proxy-rewrite': { list: ['route'] },
};

export enum PluginType {
  general = "general",
  transformation = "transformation",
  authentication = "authentication",
  security = "security",
  traffic = "traffic",
  monitoring = "monitoring",
  loggers = "loggers",
  protocol = "protocol",
  other = "other"
}

/**
 * Plugin List that contains type field
*/
export const PLUGIN_LIST = {
  "hmac-auth": {
    "type": PluginType.authentication
  },
  "serverless-post-function": {
    "type": PluginType.general
  },
  "mqtt-proxy": {
    "type": PluginType.protocol
  },
  "response-rewrite": {
    "type": PluginType.transformation
  },
  "basic-auth": {
    "type": PluginType.authentication
  },
  "error-log-logger": {
    "type": PluginType.loggers
  },
  "fault-injection": {
    "type": PluginType.transformation
  },
  "limit-count": {
    "type": PluginType.traffic
  },
  "prometheus": {
    "type": PluginType.monitoring
  },
  "proxy-rewrite": {
    "type": PluginType.transformation
  },
  "syslog": {
    "type": PluginType.loggers
  },
  "traffic-split": {
    "type": PluginType.traffic
  },
  "jwt-auth": {
    "type": PluginType.authentication
  },
  "kafka-logger": {
    "type": PluginType.loggers
  },
  "limit-conn": {
    "type": PluginType.traffic
  },
  "udp-logger": {
    "type": PluginType.loggers
  },
  "zipkin": {
    "type": PluginType.monitoring
  },
  "echo": {
    "type": PluginType.general
  },
  "log-rotate": {
    "type": PluginType.loggers
  },
  "serverless-pre-function": {
    "type": PluginType.general
  },
  "dubbo-proxy": {
    "type": PluginType.protocol
  },
  "node-status": {
    "type": PluginType.monitoring
  },
  "referer-restriction": {
    "type": PluginType.security
  },
  "api-breaker": {
    "type": PluginType.traffic
  },
  "consumer-restriction": {
    "type": PluginType.security
  },
  "cors": {
    "type": PluginType.security
  },
  "limit-req": {
    "type": PluginType.traffic
  },
  "proxy-mirror": {
    "type": PluginType.traffic
  },
  "request-validation": {
    "type": PluginType.traffic
  },
  "example-plugin": {
    "type": PluginType.other
  },
  "ip-restriction": {
    "type": PluginType.security
  },
  "key-auth": {
    "type": PluginType.authentication
  },
  "proxy-cache": {
    "type": PluginType.traffic
  },
  "redirect": {
    "type": PluginType.general
  },
  "request-id": {
    "type": PluginType.traffic
  },
  "skywalking": {
    "type": PluginType.monitoring
  },
  "batch-requests": {
    "type": PluginType.general
  },
  "http-logger": {
    "type": PluginType.loggers
  },
  "openid-connect": {
    "type": PluginType.authentication
  },
  "sls-logger": {
    "type": PluginType.loggers
  },
  "tcp-logger": {
    "type": PluginType.loggers
  },
  "uri-blocker": {
    "type": PluginType.security
  },
  "wolf-rbac": {
    "type": PluginType.authentication
  },
  "authz-keycloak": {
    "type": PluginType.authentication
  },
  "grpc-transcode": {
    "type": PluginType.transformation
  },
  "server-info": {
    "type": PluginType.general
  }
}
