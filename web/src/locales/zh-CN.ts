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
import { ActionBarZhCN } from '@/components/ActionBar';

import component from './zh-CN/component';
import globalHeader from './zh-CN/globalHeader';
import menu from './zh-CN/menu';
import pwa from './zh-CN/pwa';
import other from './zh-CN/other'
import settingDrawer from './zh-CN/settingDrawer';
import settings from './zh-CN/setting';
import Plugin from '../components/Plugin/locales/zh-CN';
import PluginFlow from '../components/PluginFlow/locales/zh-CN';
import RawDataEditor from '../components/RawDataEditor/locales/zh-CN';
import UpstreamComponent from '../components/Upstream/locales/zh-CN'

export default {
  'navBar.lang': '语言',
  'layout.user.link.help': '帮助',
  'layout.user.link.privacy': '隐私',
  'layout.user.link.terms': '条款',
  'app.preview.down.block': '下载此页面到本地项目',
  ...globalHeader,
  ...menu,
  ...settingDrawer,
  ...settings,
  ...pwa,
  ...component,
  ...other,
  ...ActionBarZhCN,
  ...Plugin,
  ...PluginFlow,
  ...RawDataEditor,
  ...UpstreamComponent
};
