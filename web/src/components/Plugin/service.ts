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
import { omit } from 'lodash';
import { request } from 'umi';

import type { PluginItem } from './data';
import { PLUGIN_LIST, PluginType, PluginState } from './data';

const cached: {
  list: PluginComponent.Meta[];
} = {
  list: [],
};

export const fetchList = async ({ enablePluginList = {} } = {}) => {
  let res: PluginItem[];
  if (cached.list.length) {
    res = cached.list;
  } else {
    res = await request<Res<PluginComponent.Meta[]>>('/plugins?all=true').then((data) => {
      const typedData = data.data.map((item) => ({
        ...item,
        type: PLUGIN_LIST[item.name]?.type || 'other',
        originType: item.type,
        hidden: PLUGIN_LIST[item.name]?.hidden || false,
      }));

      let finalList: PluginComponent.Meta[] = [];

      Object.values(PluginType).forEach((type) => {
        finalList = finalList.concat(typedData.filter((item) => item.type === type));
      });

      if (cached.list.length === 0) {
        cached.list = finalList;
      }

      return finalList;
    });
  }
  res.map((item) => {
    const isEnable = enablePluginList[item.name]?.disable === false;
    // eslint-disable-next-line no-param-reassign
    item.state = isEnable ? PluginState.enable : PluginState.disable;
    return item;
  });

  return Promise.resolve(res);
};

/**
 * cache plugin schema by schemaType
 * default schema is route for plugins in route
 * support schema: consumer for plugins in consumer
 */
const cachedPluginSchema: Record<string, any> = {
  route: {},
  consumer: {},
};
export const fetchSchema = async (
  name: string,
  schemaType: PluginComponent.Schema,
): Promise<any> => {
  if (!cachedPluginSchema[schemaType][name]) {
    const queryString = schemaType !== 'route' ? `?schema_type=${schemaType}` : '';
    cachedPluginSchema[schemaType][name] = (
      await request(`/schema/plugins/${name}${queryString}`)
    ).data;
    // for plugins schema returned with properties: [], which will cause parse error
    if (JSON.stringify(cachedPluginSchema[schemaType][name].properties) === '[]') {
      cachedPluginSchema[schemaType][name] = omit(
        cachedPluginSchema[schemaType][name],
        'properties',
      );
    }
  }
  return cachedPluginSchema[schemaType][name];
};

export const fetchPluginTemplateList = () => {
  return request<Res<ResListData<PluginTemplateModule.ResEntity>>>('/plugin_configs').then(
    (data) => {
      return data.data.rows;
    },
  );
};
