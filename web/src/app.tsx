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
import { history } from 'umi';
import type { RequestConfig } from 'umi';
import type { Settings as LayoutSettings } from '@ant-design/pro-layout';
import { isPlainObject } from 'lodash';

import RightContent from '@/components/RightContent';
import Footer from '@/components/Footer';
import { queryCurrent } from '@/services/user';
import { getMenuData, errorHandler, getUrlQuery } from '@/helpers';

import './libs/iconfont';
import defaultSettings from '../config/defaultSettings';

export async function getInitialState(): Promise<{
  currentUser?: API.CurrentUser;
  settings?: LayoutSettings;
}> {
  const token = localStorage.getItem('token');
  if (!token) {
    const redirect = getUrlQuery('redirect') || '/';
    history.replace(`/user/login?redirect=${redirect}`);
  }

  const currentUser = await queryCurrent();
  return {
    currentUser,
    settings: defaultSettings,
  };
}

export const layout = ({ initialState }: { initialState: { settings?: LayoutSettings } }) => {
  return {
    headerRender: undefined,
    rightContentRender: () => <RightContent />,
    disableContentMargin: false,
    footerRender: () => <Footer />,
    menuHeaderRender: undefined,
    menuDataRender: getMenuData,
    ...initialState?.settings,
  };
};

/* eslint no-param-reassign: ["error", { "props": true, "ignorePropertyModificationsFor": ["obj"] }] */
const nullValueFilter = (obj: Record<string, any>) => {
  Object.entries(obj).forEach(([key, value]) => {
    if (isPlainObject(value)) {
      nullValueFilter(value);
    } else if ([null, undefined].includes(value)) {
      delete obj[key];
    }
  });
};

export const request: RequestConfig = {
  prefix: '/apisix/admin',
  errorHandler,
  credentials: 'same-origin',
  requestInterceptors: [
    (url, options) => {
      const newOptions = { ...options };
      if (newOptions.data) {
        nullValueFilter(newOptions.data);
      }
      newOptions.headers = {
        ...options.headers,
        Authorization: localStorage.getItem('token') || '',
      };
      return {
        url,
        options: { ...newOptions, interceptors: true },
      };
    },
  ],
  responseInterceptors: [
    async (res) => {
      if (!res.ok) {
        // NOTE: http code >= 400, using errorHandler
        return res;
      }

      const data = await res.json();
      const { success = false } = data as Res<any>;
      if (!success) {
        // eslint-disable-next-line
        return Promise.reject({ response: res, data });
      }
      return data;
    },
  ],
};
