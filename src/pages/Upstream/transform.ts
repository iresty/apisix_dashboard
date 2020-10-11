import { pickBy, identity, omit } from 'lodash';

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
export const transformRequest = (
  formData: UpstreamModule.RequestBody,
): UpstreamModule.RequestBody | undefined => {
  let data = pickBy(formData, identity) as UpstreamModule.RequestBody;
  const { type, hash_on, key, k8s_deployment_info, nodes, pass_host, upstream_host } = data;
  data.checks = pickBy(data.checks, identity);
  if (Object.keys(data.checks).length === 0) {
    data = omit(data, 'checks');
  }
  if (nodes && k8s_deployment_info) {
    return undefined;
  }

  if (!nodes && !k8s_deployment_info) {
    return undefined;
  }

  if (type === 'chash') {
    if (!hash_on) {
      return undefined;
    }

    if (hash_on !== 'consumer' && !key) {
      return undefined;
    }
  }

  if (pass_host === 'rewrite' && !upstream_host) {
    return undefined;
  }

  if (nodes) {
    return data;
  }

  return undefined;
};
