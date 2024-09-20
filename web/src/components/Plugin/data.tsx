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
  'basic-auth': <IconFont name="iconbasic-auth" />,
  'hmac-auth': <IconFont name="iconhmac-auth" />,
  cors: <IconFont name="iconcors" />,
  'wolf-rbac': <IconFont name="iconwolf-rbac" />,
  'key-auth': <IconFont name="iconkey-auth" />,
  'request-validation': <IconFont name="iconrequest-validation" />,
  'fault-injection': <IconFont name="iconfault-injection" />,
  'consumer-restriction': <IconFont name="iconconsumer-restriction" />,
  'api-breaker': <IconFont name="iconapi-breaker" />,
  'ip-restriction': <IconFont name="iconip-restriction" />,
  'uri-blocker': <IconFont name="iconuri-blocker" />,
  'referer-restriction': <IconFont name="iconreferer-restriction" />,
  'limit-conn': <IconFont name="iconlimit-connect" />,
  'limit-req': <IconFont name="iconlimit-req" />,
  'limit-count': <IconFont name="iconlimit-count" />,
  'traffic-split': <IconFont name="icontraffic-split" />,
  'serverless-post-function': <IconFont name="iconserverless-post-function" />,
  'serverless-pre-function': <IconFont name="iconserverless-pre-function" />,
  'response-rewrite': <IconFont name="iconresponse-rewrite" />,
  'proxy-rewrite': <IconFont name="iconproxy-rewrite" />,
  'proxy-mirror': <IconFont name="iconproxy-mirror" />,
  'proxy-cache': <IconFont name="iconproxy-cache" />,
  'grpc-transcode': <IconFont name="icongrpc-transcode" />,
  'batch-requests': <IconFont name="iconbatch-request" />,
  zipkin: <IconFont name="iconzipkin" />,
  'udp-logger': <IconFont name="iconudp-logger" />,
  'error-log-logger': <IconFont name="iconerror-logger" />,
  'http-logger': <IconFont name="iconhttp-logger" />,
  'loki-logger': <IconFont name="iconloki-logger" />,
  'request-id': <IconFont name="iconrequest-id" />,
  'sls-logger': <IconFont name="iconsls-logger" />,
  syslog: <IconFont name="iconsys-logger" />,
  'tcp-logger': <IconFont name="icontcp-logger" />,
};

// This list is used to filter out plugins that cannot be displayed in the plugins list.
export const PLUGIN_FILTER_LIST: Record<string, { list: PluginComponent.ReferPage[] }> = {
  redirect: { list: ['route'] }, // Filter out the redirect plugin on the route page.
  'proxy-rewrite': { list: ['route'] },
};

export enum PluginType {
  authentication = 'authentication',
  security = 'security',
  traffic = 'traffic',
  serverless = 'serverless',
  observability = 'observability',
  other = 'other',
}

/**
 * Plugin List that contains type field
 */
export const PLUGIN_LIST = {
  'api-breaker': {
    type: PluginType.security,
  },
  'authz-casbin': {
    type: PluginType.authentication,
  },
  'authz-casdoor': {
    type: PluginType.authentication,
  },
  'authz-keycloak': {
    type: PluginType.authentication,
  },
  'aws-lambda': {
    type: PluginType.serverless,
  },
  'azure-functions': {
    type: PluginType.serverless,
  },
  'basic-auth': {
    type: PluginType.authentication,
  },
  'batch-requests': {
    type: PluginType.other,
    hidden: true,
  },
  'clickhouse-logger': {
    type: PluginType.observability,
  },
  'client-control': {
    type: PluginType.traffic,
  },
  'consumer-restriction': {
    type: PluginType.security,
  },
  cors: {
    type: PluginType.security,
  },
  csrf: {
    type: PluginType.security,
  },
  datadog: {
    type: PluginType.observability,
  },
  'dubbo-proxy': {
    type: PluginType.other,
    hidden: true,
  },
  echo: {
    type: PluginType.other,
    hidden: true,
  },
  'error-log-logger': {
    type: PluginType.observability,
    hidden: true,
  },
  'example-plugin': {
    type: PluginType.other,
    hidden: true,
  },
  'ext-plugin-post-req': {
    type: PluginType.other,
  },
  'ext-plugin-post-resp': {
    type: PluginType.other,
  },
  'ext-plugin-pre-req': {
    type: PluginType.other,
  },
  'fault-injection': {
    type: PluginType.security,
  },
  'file-logger': {
    type: PluginType.observability,
  },
  'forward-auth': {
    type: PluginType.authentication,
  },
  'google-cloud-logging': {
    type: PluginType.observability,
  },
  'grpc-transcode': {
    type: PluginType.other,
  },
  'grpc-web': {
    type: PluginType.other,
  },
  gzip: {
    type: PluginType.other,
  },
  'hmac-auth': {
    type: PluginType.authentication,
  },
  'http-logger': {
    type: PluginType.observability,
  },
  'loki-logger': {
    type: PluginType.observability,
  },
  'ip-restriction': {
    type: PluginType.security,
  },
  'jwt-auth': {
    type: PluginType.authentication,
  },
  'kafka-logger': {
    type: PluginType.observability,
  },
  'kafka-proxy': {
    type: PluginType.other,
  },
  'key-auth': {
    type: PluginType.authentication,
  },
  'ldap-auth': {
    type: PluginType.authentication,
  },
  'limit-conn': {
    type: PluginType.traffic,
  },
  'limit-count': {
    type: PluginType.traffic,
  },
  'limit-req': {
    type: PluginType.traffic,
  },
  loggly: {
    type: PluginType.observability,
  },
  'log-rotate': {
    type: PluginType.observability,
    hidden: true,
  },
  mocking: {
    type: PluginType.other,
  },
  'node-status': {
    type: PluginType.other,
    hidden: true,
  },
  opa: {
    type: PluginType.authentication,
  },
  'openid-connect': {
    type: PluginType.authentication,
  },
  opentelemetry: {
    type: PluginType.observability,
  },
  openwhisk: {
    type: PluginType.serverless,
  },
  prometheus: {
    type: PluginType.observability,
  },
  'proxy-cache': {
    type: PluginType.other,
  },
  'proxy-control': {
    type: PluginType.other,
  },
  'proxy-mirror': {
    type: PluginType.other,
  },
  'proxy-rewrite': {
    type: PluginType.other,
  },
  'public-api': {
    type: PluginType.security,
  },
  'real-ip': {
    type: PluginType.other,
  },
  redirect: {
    type: PluginType.other,
    hidden: true,
  },
  'referer-restriction': {
    type: PluginType.security,
  },
  'request-id': {
    type: PluginType.observability,
  },
  'request-validation': {
    type: PluginType.security,
  },
  'response-rewrite': {
    type: PluginType.other,
  },
  'rocketmq-logger': {
    type: PluginType.observability,
  },
  'server-info': {
    type: PluginType.other,
    hidden: true,
  },
  'serverless-post-function': {
    type: PluginType.serverless,
  },
  'serverless-pre-function': {
    type: PluginType.serverless,
  },
  skywalking: {
    type: PluginType.observability,
  },
  'skywalking-logger': {
    type: PluginType.observability,
  },
  'sls-logger': {
    type: PluginType.observability,
  },
  'splunk-hec-logging': {
    type: PluginType.observability,
  },
  syslog: {
    type: PluginType.observability,
  },
  'tcp-logger': {
    type: PluginType.observability,
  },
  'traffic-split': {
    type: PluginType.traffic,
  },
  'ua-restriction': {
    type: PluginType.security,
  },
  'udp-logger': {
    type: PluginType.observability,
  },
  'uri-blocker': {
    type: PluginType.security,
  },
  'wolf-rbac': {
    type: PluginType.authentication,
  },
  zipkin: {
    type: PluginType.observability,
  },
  'mqtt-proxy': {
    type: PluginType.other,
    hidden: true,
  },
  'elasticsearch-logger': {
    type: PluginType.observability,
  },
  openfunction: {
    type: PluginType.serverless,
  },
  'tencent-cloud-cls': {
    type: PluginType.observability,
  },
  ai: {
    type: PluginType.other,
  },
  'cas-auth': {
    type: PluginType.authentication,
  },
};
