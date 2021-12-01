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
export default {
  'component.plugin.tip1': 'NOTE: After customizing the plugin, you need to update schema.json.',
  'component.plugin.tip2': 'How to update?',
  'component.select.pluginTemplate': 'Select a plugin template',
  'component.step.select.pluginTemplate.select.option': 'Custom',
  'component.plugin.pluginTemplate.tip1':
    '1. When a route already have plugins field configured, the plugins in the plugin template will be merged into it.',
  'component.plugin.pluginTemplate.tip2':
    '2. The same plugin in the plugin template will override one in the plugins.',
  'component.plugin.enable': 'Enable',
  'component.plugin.disable': 'Disable',
  'component.plugin.authentication': 'Authentication',
  'component.plugin.security': 'Security',
  'component.plugin.traffic': 'Traffic Control',
  'component.plugin.serverless': 'Serverless',
  'component.plugin.observability': 'Observability',
  'component.plugin.other': 'Other',

  // cors
  'component.pluginForm.cors.allow_origins.tooltip':
    'Which Origins is allowed to enable CORS, format as：scheme://host:port, for example: https://somehost.com:8081. Multiple origin use , to split. When allow_credential is false, you can use * to indicate allow any origin. you also can allow all any origins forcefully using ** even already enable allow_credential, but it will bring some security risks.',
  'component.pluginForm.cors.allow_origins.extra': 'For example: https://somehost.com:8081',
  'component.pluginForm.cors.allow_methods.tooltip':
    'Which Method is allowed to enable CORS, such as: GET, POST etc. Multiple method use , to split. When allow_credential is false, you can use * to indicate allow all any method. You also can allow any method forcefully using ** even already enable allow_credential, but it will bring some security risks.',
  'component.pluginForm.cors.allow_headers.tooltip':
    'Which headers are allowed to set in request when access cross-origin resource. Multiple value use , to split. When allow_credential is false, you can use * to indicate allow all request headers. You also can allow any header forcefully using ** even already enable allow_credential, but it will bring some security risks.',
  'component.pluginForm.cors.expose_headers.tooltip':
    'Which headers are allowed to set in response when access cross-origin resource. Multiple value use , to split. When allow_credential is false, you can use * to indicate allow any header. You also can allow any header forcefully using ** even already enable allow_credential, but it will bring some security risks.',
  'component.pluginForm.cors.max_age.tooltip':
    'Maximum number of seconds the results can be cached. Within this time range, the browser will reuse the last check result. -1 means no cache. Please note that the maximum value is depended on browser, please refer to MDN for details.',
  'component.pluginForm.cors.allow_credential.tooltip':
    "If you set this option to true, you can not use '*' for other options.",
  'component.pluginForm.cors.allow_origins_by_regex.tooltip':
    'Use regex expressions to match which origin is allowed to enable CORS, for example, [".*.test.com"] can use to match all subdomain of test.com.',

  // referer-restriction
  'component.pluginForm.referer-restriction.whitelist.tooltip':
    'List of hostname to whitelist. The hostname can be started with * as a wildcard.',
  'component.pluginForm.referer-restriction.blacklist.tooltip':
    'List of hostname to blacklist. The hostname can be started with * as a wildcard.',
  'component.pluginForm.referer-restriction.listEmpty.tooltip': 'List empty',
  'component.pluginForm.referer-restriction.bypass_missing.tooltip':
    'Whether to bypass the check when the Referer header is missing or malformed.',
  'component.pluginForm.referer-restriction.message.tooltip':
    'Message returned in case access is not allowed.',

  // api-breaker
  'component.pluginForm.api-breaker.break_response_code.tooltip':
    'Return error code when unhealthy.',
  'component.pluginForm.api-breaker.max_breaker_sec.tooltip': 'Maximum breaker time(seconds).',
  'component.pluginForm.api-breaker.unhealthy.http_statuses.tooltip':
    'Status codes when unhealthy.',
  'component.pluginForm.api-breaker.unhealthy.failures.tooltip':
    'Number of consecutive error requests that triggered an unhealthy state.',
  'component.pluginForm.api-breaker.healthy.http_statuses.tooltip': 'Status codes when healthy.',
  'component.pluginForm.api-breaker.healthy.successes.tooltip':
    'Number of consecutive normal requests that trigger health status.',

  // proxy-mirror
  'component.pluginForm.proxy-mirror.host.tooltip':
    'Specify a mirror service address, e.g. http://127.0.0.1:9797 (address needs to contain schema: http or https, not URI part)',
  'component.pluginForm.proxy-mirror.host.extra': 'e.g. http://127.0.0.1:9797',
  'component.pluginForm.proxy-mirror.host.ruletip':
    'address needs to contain schema: http or https, not URI part',
  'component.pluginForm.proxy-mirror.sample_ratio.tooltip':
    'the sample ratio that requests will be mirrored.',

  // limit-conn
  'component.pluginForm.limit-conn.conn.tooltip':
    'the maximum number of concurrent requests allowed. Requests exceeding this ratio (and below conn + burst) will get delayed(the latency seconds is configured by default_conn_delay) to conform to this threshold.',
  'component.pluginForm.limit-conn.burst.tooltip':
    'the number of excessive concurrent requests (or connections) allowed to be delayed.',
  'component.pluginForm.limit-conn.default_conn_delay.tooltip':
    'the latency seconds of request when concurrent requests exceeding conn but below (conn + burst).',
  'component.pluginForm.limit-conn.key.tooltip':
    'to limit the concurrency level. For example, one can use the host name (or server zone) as the key so that we limit concurrency per host name. Otherwise, we can also use the client address as the key so that we can avoid a single client from flooding our service with too many parallel connections or requests. Now accept those as key: "remote_addr"(client\'s IP), "server_addr"(server\'s IP), "X-Forwarded-For/X-Real-IP" in request header, "consumer_name"(consumer\'s username).',
  'component.pluginForm.limit-conn.rejected_code.tooltip':
    'returned when the request exceeds conn + burst will be rejected.',
  'component.pluginForm.limit-conn.rejected_msg.tooltip':
    'the response body returned when the request exceeds conn + burst will be rejected.',
  'component.pluginForm.limit-conn.only_use_default_delay.tooltip':
    'enable the strict mode of the latency seconds. If you set this option to true, it will run strictly according to the latency seconds you set without additional calculation logic.',
  'component.pluginForm.limit-conn.allow_degradation.tooltip':
    'Whether to enable plugin degradation when the limit-conn function is temporarily unavailable. Allow requests to continue when the value is set to true, default false.',

  // limit-req
  'component.pluginForm.limit-req.rate.tooltip':
    'The specified request rate (number per second) threshold. Requests exceeding this rate (and below burst) will get delayed to conform to the rate.',
  'component.pluginForm.limit-req.burst.tooltip':
    'The number of excessive requests per second allowed to be delayed. Requests exceeding this hard limit will get rejected immediately.',
  'component.pluginForm.limit-req.key.tooltip':
    'The user specified key to limit the rate, now accept those as key: "remote_addr"(client\'s IP), "server_addr"(server\'s IP), "X-Forwarded-For/X-Real-IP" in request header, "consumer_name"(consumer\'s username).',
  'component.pluginForm.limit-req.rejected_code.tooltip':
    'The HTTP status code returned when the request exceeds the threshold is rejected.',
  'component.pluginForm.limit-req.nodelay.tooltip':
    'If nodelay flag is true, bursted requests will not get delayed',

  'component.plugin.form': 'Form',
  'component.plugin.format-codes.disable': 'Format JSON or YAML data',
  'component.plugin.editor': 'Plugin Editor',
  'component.plugin.noConfigurationRequired': "Doesn't need configuration",

  // limit-count
  'component.pluginForm.limit-count.count.tooltip': 'The specified number of requests threshold.',
  'component.pluginForm.limit-count.time_window.tooltip':
    'The time window in seconds before the request count is reset.',
  'component.pluginForm.limit-count.key.tooltip':
    'The user specified key to limit the count, now accept those as key: "remote_addr"(client\'s IP), "server_addr"(server\'s IP), "X-Forwarded-For/X-Real-IP" in request header, "consumer_name"(consumer\'s username) and "service_id".',
  'component.pluginForm.limit-count.rejected_code.tooltip':
    'The HTTP status code returned when the request exceeds the threshold is rejected, default 503.',
  'component.pluginForm.limit-count.policy.tooltip':
    'The rate-limiting policies to use for retrieving and incrementing the limits. Available values are local(the counters will be stored locally in-memory on the node) and redis(counters are stored on a Redis server and will be shared across the nodes, usually use it to do the global speed limit) and redis-cluster(the same function as redis, only use Redis cluster pattern).',
  'component.pluginForm.limit-count.redis_host.tooltip':
    'When using the redis policy, this property specifies the address of the Redis server.',
  'component.pluginForm.limit-count.redis_port.tooltip':
    'When using the redis policy, this property specifies the port of the Redis server.',
  'component.pluginForm.limit-count.redis_password.tooltip':
    'When using the redis policy, this property specifies the password of the Redis server.',
  'component.pluginForm.limit-count.redis_database.tooltip':
    'When using the redis policy, this property specifies the database you selected of the Redis server, and only for non Redis cluster mode (single instance mode or Redis public cloud service that provides single entry).',
  'component.pluginForm.limit-count.redis_timeout.tooltip':
    'When using the redis policy, this property specifies the timeout in milliseconds of any command submitted to the Redis server.',
  'component.pluginForm.limit-count.redis_cluster_nodes.tooltip':
    'When using redis-cluster policy, this property is a list of addresses of Redis cluster service nodes (at least two nodes).',
  'component.pluginForm.limit-count.redis_cluster_name.tooltip':
    'When using redis-cluster policy, this property is the name of Redis cluster service nodes.',
  'component.pluginForm.limit-count.atLeast2Characters.rule': 'Please enter at least 2 characters',

  // request-id
  'component.pluginForm.request-id.header_name.tooltip': 'Request ID header name',
  'component.pluginForm.request-id.include_in_response.tooltip':
    'Option to include the unique request ID in the response header',
};
