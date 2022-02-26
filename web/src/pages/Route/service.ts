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
import { request } from 'umi';
import { pickBy, identity } from 'lodash';

import { transformStepData, transformRouteData, transformUpstreamNodes } from './transform';
import { transformLabelList } from '@/helpers';

export const create = (data: RouteModule.RequestData, mode?: RouteModule.RequestMode) =>
  request(`/routes`, {
    method: 'POST',
    data: mode === 'RawData' ? data : transformStepData(data),
  });

export const update = (
  rid: string,
  data: RouteModule.RequestData,
  mode?: RouteModule.RequestMode,
) =>
  request(`/routes/${rid}`, {
    method: 'PUT',
    data: mode === 'RawData' ? data : transformStepData(data),
  });

export const fetchItem = (rid: number) =>
  request(`/routes/${rid}`).then((data) => transformRouteData(data.data));

export const fetchList = ({ current = 1, pageSize = 10, ...res }) => {
  const { labels = [], API_VERSION = [], status } = res;

  return request<Res<ResListData<RouteModule.ResponseBody>>>('/routes', {
    params: {
      name: res.name,
      uri: res.uri,
      label: labels.concat(API_VERSION).join(','),
      page: current,
      page_size: pageSize,
      status,
    },
  }).then(({ data }) => {
    return {
      data: data.rows,
      total: data.total_size,
    };
  });
};

export const remove = (rid: string) => request(`/routes/${rid}`, { method: 'DELETE' });

export const checkUniqueName = (name = '', exclude = '') =>
  request('/notexist/routes', {
    params: pickBy(
      {
        name,
        exclude,
      },
      identity,
    ),
  });

export const fetchUpstreamItem = (sid: string) => {
  return request(`/upstreams/${sid}`).then(({ nodes, timeout, id }) => {
    return {
      upstreamHostList: transformUpstreamNodes(nodes),
      timeout,
      upstream_id: id,
    };
  });
};

export const checkHostWithSSL = (hosts: string[]) =>
  request('/check_ssl_exists', {
    method: 'POST',
    data: { hosts },
  });

export const fetchLabelList = () =>
  request('/labels/route').then(({ data }) => transformLabelList(data.rows) as LabelList);

export const updateRouteStatus = (rid: string, status: RouteModule.RouteStatus) =>
  request(`/routes/${rid}`, {
    method: 'PATCH',
    data: { status },
  });

export const debugRoute = (headers: any, data: undefined | string | FormData) => {
  return request('/debug-request-forwarding', {
    method: 'post',
    data,
    headers,
  });
};

export const fetchServiceList = () =>
  request('/services').then(({ data }) => ({
    data: data.rows,
    total: data.total_size,
  }));

export const exportRoutes = (ids?: string) => {
  return request(`/export/routes/${ids}`);
};

export const importRoutes = (formData: FormData) => {
  return request('/import/routes', {
    method: 'POST',
    data: formData,
    requestType: 'form',
  });
};
