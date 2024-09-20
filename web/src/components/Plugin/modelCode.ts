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
 * Model Code List of Models of Monaco editor
 */

/**
 * Model code of authentication type as fllows:
 */

export const authzcasbin = `{
    "model_path":
    "policy_path":
    "username":
}`;

export const authzkeycloak = `{
    "token_endpoint":
    "permissions":
    "audience":
}
`;

export const forwardauth = `{
    "uri":
    "request_headers":
    "upstream_headers":
    "client_headers":
}
`;

export const opa = `{
    "type":
    "request": {
        "scheme":
        "path":
        "headers": {
            "user-agent":
            "accept":
            "host":
        },
        "query":
        "port":
        "method":
        "host":
    },
    "var": {
        "timestamp":
        "server_addr":
        "server_port":
        "remote_port":
        "remote_addr":
    },
    "route":
    "service":
    "consumer":
}
`;

export const openidconnect = `{
    "client_id":
    "client_secret":
    "discovery":
    "introspection_endpoint":
    "bearer_only":
    "realm":
    "introspection_endpoint_auth_method":
}
`;

/**
 * Model code of security type as fllows:
 */

export const csrf = `{
    "key":
}
`;

export const iprestriction = `{
    "whitelist":
}
`;

export const uarestriction = `{
    "bypass_missing":
    "allowlist":
    "denylist":
}
`;

export const uriblocker = `{
    "block_rules":
}
`;

/**
 * Model code of traffic type as fllows:
 */

export const clientcontrol = `{
    "max_body_size":
}
`;

export const trafficsplit = `{
    "rules": [
        {
            "weighted_upstreams": [
                {
                    "upstream": {
                        "name":
                        "type":
                        "nodes":
                        "timeout": {
                            "connect":
                            "send":
                            "read":
                        }
                    },
                    "weight":
                },
                {
                    "weight":
                }
            ]
        }
    ]
}
`;

/**
 * Model code of serverless type as fllows:
 */

export const awslambda = `{
  "token_endpoint":
  "permissions":
  "audience":
}
`;

export const azurefunctions = `{
    "azure-functions": {
        "function_uri":
        "authorization": {
            "apikey":
        }
    }
}
`;

export const openwhisk = `{
    "api_host":
    "service_token":
    "namespace":
    "action":
}
`;

/**
 * Model code of observability type as fllows:
 */

export const clickhouselogger = `{
    "user":
    "password":
    "database":
    "logtable":
    "endpoint_addr":
}
`;

export const filelogger = `{
    "path":
}
`;

export const googlecloudlogging = `{
    "auth_config":{
        "project_id":
        "private_key":
        "token_uri":
        "scopes":
        "entries_uri":
    },
    "resource":{
        "type":
    },
    "log_id":
    "inactive_timeout":
    "max_retry_count":
    "buffer_duration":
    "retry_delay":
    "batch_max_size":
}
`;

export const httplogger = `{
    "uri":
}
`;

export const lokilogger = `{
    "endpoint_addrs": []
}
`;

export const kafkalogger = `{
    "broker_list" :
    "kafka_topic" :
    "key" :
    "batch_max_size":
    "name":
}
`;

export const loggly = `{
    "nameserver_list":
    "topic":
    "batch_max_size":
    "name":
}
`;

export const rocketmqlogger = `{
  "token_endpoint":
  "permissions":
  "audience":
}
`;

export const skywalking = `{
    "sample_ratio":
}
`;

export const skywalkinglogger = `{
    "endpoint_addr":
}
`;

export const slslogger = `{
    "host":
    "port":
    "project":
    "logstore":
    "access_key_id":
    "access_key_secret":
    "timeout":
}
`;

export const splunkheclogging = `{
    "endpoint":{
        "uri":
        "token":
        "channel":
        "timeout":
    },
    "buffer_duration":
    "max_retry_count":
    "retry_delay":
    "inactive_timeout":
    "batch_max_size":
}
`;

export const syslog = `{
    "host":
    "port":
    "flush_limit":
}
`;

export const tcplogger = `{
    "host":
    "port":
    "tls":
    "batch_max_size":
    "name":
}
`;

export const zipkin = `{
    "endpoint":
    "sample_ratio":
    "service_name":
    "server_addr":
}
`;

/**
 * Model code of other type as fllows:
 */

export const extpluginprereq = `{
    "conf":
}
`;

export const realip = `{
    "source":
    "trusted_addresses":
}
`;
