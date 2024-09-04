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
  'page.upstream.step.select.upstream': 'अपस्ट्रीम',
  'page.upstream.step.select.upstream.select.option': 'कस्टम',
  'page.upstream.step.select.upstream.select.option.serviceSelected':
    'कस्टम (वर्तमान कॉन्फ़िगरेशन सेवा से ओवरराइड करेगी)',
  'page.upstream.step.select.upstream.select.none': 'कोई नहीं',
  'page.upstream.step.backend.server.domain.or.ip': 'बैकएंड सर्वर होस्ट/IP',
  'page.upstream.form.item-label.node.domain.or.ip': 'लक्ष्य',
  'page.upstream.step.input.domain.name.or.ip': 'कृपया डोमेन या IP दर्ज करें',
  'page.upstream.step.valid.domain.name.or.ip': 'कृपया वैध डोमेन या IP दर्ज करें',
  'page.upstream.step.domain.name.or.ip': 'होस्टनाम या IP',
  'page.upstream.step.host': 'होस्ट',
  'page.upstream.step.input.port': 'कृपया पोर्ट नंबर दर्ज करें',
  'page.upstream.step.port': 'पोर्ट',
  'page.upstream.step.input.weight': 'कृपया वजन दर्ज करें',
  'page.upstream.step.weight': 'वजन',
  'page.upstream.step.create': 'बनाएं',
  'page.upstream.step.name': 'नाम',
  'page.upstream.step.name.should.unique': 'नाम अद्वितीय होना चाहिए',
  'page.upstream.step.input.upstream.name': 'कृपया अपस्ट्रीम नाम दर्ज करें',
  'page.upstream.step.description': 'विवरण',
  'page.upstream.step.input.description': 'कृपया अपस्ट्रीम का विवरण दर्ज करें',
  'page.upstream.step.type': 'एल्गोरिदम',
  'page.upstream.step.pass-host': 'होस्टनाम',
  'page.upstream.step.pass-host.pass': 'ग्राहक अनुरोध से समान होस्ट बनाए रखें',
  'page.upstream.step.pass-host.node': 'उपनोद सूची से डोमेन या IP का उपयोग करें',
  'page.upstream.step.pass-host.rewrite': 'कस्टम होस्ट (भविष्य में समाप्त हो जाएगा)',
  'page.upstream.step.pass-host.upstream_host': 'कस्टम होस्ट',
  'page.upstream.step.connect.timeout': 'कनेक्ट टाइमआउट',
  'page.upstream.step.connect.timeout.desc':
    'अनुरोध से अपस्ट्रीम सर्वर के साथ एक कनेक्शन स्थापित करने के लिए समय सीमा',
  'page.upstream.step.input.connect.timeout': 'कृपया कनेक्ट टाइमआउट दर्ज करें',
  'page.upstream.step.send.timeout': 'भेजें टाइमआउट',
  'page.upstream.step.send.timeout.desc': 'अपस्ट्रीम सर्वर्स को डेटा भेजने के लिए समय सीमा',
  'page.upstream.step.input.send.timeout': 'कृपया भेजें टाइमआउट दर्ज करें',
  'page.upstream.step.read.timeout': 'रीड टाइमआउट',
  'page.upstream.step.read.timeout.desc': 'अपस्ट्रीम सर्वर्स से डेटा प्राप्त करने के लिए समय सीमा',
  'page.upstream.step.input.read.timeout': 'कृपया रीड टाइमआउट दर्ज करें',
  'page.upstream.step.healthyCheck.healthy.check': 'स्वास्थ्य जाँच',
  'page.upstream.step.healthyCheck.healthy': 'स्वस्थ',
  'page.upstream.step.healthyCheck.unhealthy': 'अस्वस्थ',
  'page.upstream.step.healthyCheck.healthy.status': 'स्वस्थ स्थिति',
  'page.upstream.step.healthyCheck.unhealthyStatus': 'अस्वस्थ स्थिति',
  'page.upstream.step.healthyCheck.active': 'सक्रिय',
  'page.upstream.step.healthyCheck.active.timeout': 'टाइमआउट',
  'page.upstream.step.input.healthyCheck.active.timeout': 'कृपया टाइमआउट दर्ज करें',
  'page.upstream.step.input.healthyCheck.activeInterval': 'कृपया अंतराल दर्ज करें',
  'page.upstream.step.input.healthyCheck.active.req_headers': 'कृपया अनुरोध हेडर्स दर्ज करें',
  'page.upstream.step.healthyCheck.passive': 'निष्क्रिय',
  'page.upstream.step.healthyCheck.passive.http_statuses': 'HTTP स्थिति',
  'page.upstream.step.input.healthyCheck.passive.http_statuses': 'कृपया HTTP स्थिति दर्ज करें',
  'page.upstream.step.healthyCheck.passive.tcp_failures': 'TCP असफलता',
  'page.upstream.step.input.healthyCheck.passive.tcp_failures': 'कृपया TCP असफलता दर्ज करें',
  'page.upstream.step.keepalive_pool': 'कीपएलाइव पूल',
  'page.upstream.notificationMessage.enableHealthCheckFirst': 'कृपया पहले स्वास्थ्य जाँच सक्षम करें।',
  'page.upstream.upstream_host.required': 'कृपया कस्टम होस्ट दर्ज करें',
  'page.upstream.create': 'उपस्ट्रीम बनाएं',
  'page.upstream.configure': 'उपस्ट्रीम कॉन्फ़िगर करें',
  'page.upstream.create.upstream.successfully': 'उपस्ट्रीम सफलतापूर्वक बनाएं',
  'page.upstream.edit.upstream.successfully': 'उपस्ट्रीम सफलतापूर्वक कॉन्फ़िगर करें',
  'page.upstream.create.basic.info': 'मौलिक जानकारी',
  'page.upstream.create.preview': 'पूर्वावलोकन',
  'page.upstream.list.id': 'आईडी',
  'page.upstream.list.name': 'नाम',
  'page.upstream.list.type': 'प्रकार',
  'page.upstream.list.description': 'विवरण',
  'page.upstream.list.edit.time': 'अपडेट समय',
  'page.upstream.list.operation': 'क्रिया',
  'page.upstream.list.edit': 'कॉन्फ़िगर',
  'page.upstream.list.confirm.delete': 'क्या आपको हटाना है?',
  'page.upstream.list.confirm': 'कन्फर्म',
  'page.upstream.list.cancel': 'रद्द करें',
  'page.upstream.list.delete.successfully': 'उपस्ट्रीम सफलतापूर्वक हटाएं',
  'page.upstream.list.delete': 'हटाएं',
  'page.upstream.list': 'उपस्ट्रीम सूची',
  'page.upstream.list.input': 'कृपया दर्ज करें',
  'page.upstream.list.create': 'बनाएं',
  'page.upstream.type.roundrobin': 'राउंड रॉबिन',
  'page.upstream.type.chash': 'सीहैश',
  'page.upstream.type.ewma': 'ईडब्ल्यूएमए',
  'page.upstream.type.least_conn': 'सबसे कम कनेक्शन',
  'page.upstream.list.content':
    'उपस्ट्रीम सूची ने बनाई गई उपस्ट्रीम सेवाएँ (यानी, बैकएंड सेवाएँ) को समाहित करती हैं और उपस्ट्रीम सेवाओं के कई लक्ष्य नोड का लोड बैलेंसिंग और स्वास्थ्य जाँच करने की अनुमति देती है।',
  'page.upstream.checks.active.timeout.description':
    'सक्रिय जाँचों के लिए सॉकेट टाइमआउट (सेकंड में)',
  'page.upstream.checks.active.unhealthy.interval.description':
    'अस्वस्थ लक्ष्यों के लिए जाँचों के बीच अंतराल (सेकंड में)',
  'page.upstream.checks.passive.healthy.http_statuses.description':
    'कौन सी HTTP स्थितियों को सफल मानना चाहिए',
  'page.upstream.checks.passive.unhealthy.http_statuses.description':
    'कौन सी HTTP स्थितियों को असफल मानना चाहिए',
  'page.upstream.checks.passive.unhealthy.http_failures.description':
    'लक्ष्य को अस्वस्थ मानने के लिए HTTP असफलताओं की संख्या',
  'page.upstream.checks.passive.unhealthy.tcp_failures.description':
    'लक्ष्य को अस्वस्थ मानने के लिए TCP असफलताओं की संख्या',
  'page.upstream.scheme': 'स्कीम',
  'page.upstream.other.configuration.invalid': 'कृपया उपस्ट्रीम कॉन्फ़िगरेशन की जाँच करें',
};





