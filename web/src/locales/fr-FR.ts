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
import { ActionBarEnUS } from '@/components/ActionBar';

import Plugin from '../components/Plugin/locales/fr-Fr';
import PluginFlow from '../components/PluginFlow/locales/fr-Fr';
import RawDataEditor from '../components/RawDataEditor/locales/fr-FR';
import UpstreamComponent from '../components/Upstream/locales/fr-FR';
import component from './fr-FR/component';
import globalHeader from './fr-FR/globalHeader';
import menu from './fr-FR/menu';
import other from './fr-FR/other';
import pwa from './fr-FR/pwa';
import settings from './fr-FR/setting';
import settingDrawer from './fr-FR/settingDrawer';

export default {
  'navBar.lang': 'Langues',
  'layout.user.link.help': 'Aide',
  'layout.user.link.privacy': 'Confidentialité',
  'layout.user.link.terms': 'Termes',
  'app.preview.down.block': 'Téléchargez cette page sur votre projet local',
  ...globalHeader,
  ...menu,
  ...settingDrawer,
  ...settings,
  ...pwa,
  ...component,
  ...other,
  ...ActionBarEnUS,
  ...Plugin,
  ...PluginFlow,
  ...RawDataEditor,
  ...UpstreamComponent,
};
