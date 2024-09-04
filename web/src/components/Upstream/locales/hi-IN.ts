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
  'component.upstream.fields.tls.client_key': 'क्लाइंट कुंजी',
  'component.upstream.fields.tls.client_key.required': 'कृपया क्लाइंट कुंजी दर्ज करें',
  'component.upstream.fields.tls.client_cert': 'क्लाइंट सर्टिफिकेट',
  'component.upstream.fields.tls.client_cert.required': 'कृपया क्लाइंट सर्टिफिकेट दर्ज करें',
  'component.upstream.fields.upstream_type': 'अपस्ट्रीम प्रकार',
  'component.upstream.fields.upstream_type.node': 'नोड',
  'component.upstream.fields.upstream_type.service_discovery': 'सेवा खोज',
  'component.upstream.fields.discovery_type': 'खोज प्रकार',
  'component.upstream.fields.discovery_type.tooltip': 'खोज प्रकार',
  'component.upstream.fields.discovery_type.placeholder': 'कृपया खोज प्रकार चुनें',
  'component.upstream.fields.discovery_type.type.dns': 'डीएनएस',
  'component.upstream.fields.discovery_type.type.consul': 'कॉन्सुल',
  'component.upstream.fields.discovery_type.type.consul_kv': 'कॉन्सुल KV',
  'component.upstream.fields.discovery_type.type.nacos': 'नैकोस',
  'component.upstream.fields.discovery_type.type.eureka': 'युरेका',
  'component.upstream.fields.discovery_type.type.kubernetes': 'कुबरनेटीस',
  'component.upstream.fields.discovery_args.group_name': 'समूह का नाम',
  'component.upstream.fields.discovery_args.group_name.tooltip': 'समूह का नाम',
  'component.upstream.fields.discovery_args.group_name.placeholder': 'कृपया समूह का नाम दर्ज करें',
  'component.upstream.fields.discovery_args.namespace_id': 'नेमस्पेस आईडी',
  'component.upstream.fields.discovery_args.namespace_id.tooltip': 'नेमस्पेस आईडी',
  'component.upstream.fields.discovery_args.namespace_id.placeholder':
    'कृपया नेमस्पेस आईडी दर्ज करें',
  'component.upstream.fields.service_name': 'सेवा का नाम',
  'component.upstream.fields.service_name.tooltip': 'सेवा का नाम',
  'component.upstream.fields.service_name.placeholder': 'कृपया सेवा का नाम दर्ज करें',
  'component.upstream.fields.scheme.tooltip.stream':
    'यह प्रकार केवल स्ट्रीम रूट के लिए उपयोग किया जाता है, जो एक लेयर 4 प्रॉक्सी है। संदर्भ: https://apisix.apache.org/docs/apisix/stream-proxy/',
  'component.upstream.fields.scheme.tooltip.pubsub':
    'यह प्रकार केवल प्रकाशन सब्सक्रिप्शन में उपयोग होता है। संदर्भ: https://apisix.apache.org/docs/apisix/pubsub/',
  'component.upstream.fields.tls': 'टीएलएस',
  'component.upstream.fields.tls.tooltip': 'टीएलएस सर्टिफिकेट',
  'component.upstream.fields.hash_on': 'हैश ऑन',
  'component.upstream.fields.hash_on.tooltip': 'हैशिंग इनपुट के रूप में क्या उपयोग करें',
  'component.upstream.fields.key': 'कुंजी',
  'component.upstream.fields.key.tooltip': 'हैशिंग इनपुट के रूप में कुंजी',
  'component.upstream.fields.retries': 'पुनः प्रयास',
  'component.upstream.fields.retries.tooltip':
    'पुनः प्रयास की योजना अगले अपस्ट्रीम नोड को भेजती है। 0 का मूल्य पुनः प्रयास योजना को निष्क्रिय कर देता है और उपलब्ध बैकएंड नोड की संख्या का उपयोग करता है।',
  'component.upstream.fields.retry_timeout': 'पुनः प्रयास समय सीमा',
  'component.upstream.fields.retry_timeout.tooltip':
    'पुनः प्रयास जारी रखने के लिए सेकंड्स की मात्रा को सीमित करने के लिए एक संख्या कॉन्फ़िगर करें, और यदि पिछला अनुरोध और पुनः प्रयास लंबे समय तक ले लें, तो पुनः प्रयास जारी न करें। 0 का मतलब है पुनः प्रयास समय सीमा योजना को निष्क्रिय करें।',
  'component.upstream.fields.keepalive_pool': 'कीपअलाइव पूल',
  'component.upstream.fields.keepalive_pool.tooltip': 'अपस्ट्रीम के लिए स्वतंत्र कीपअलाइव पूल सेट करें',
  'component.upstream.fields.keepalive_pool.size': 'आकार',
  'component.upstream.fields.keepalive_pool.size.placeholder': 'कृपया आकार दर्ज करें',
  'component.upstream.fields.keepalive_pool.idle_timeout': 'निष्क्रिय समय सीमा',
  'component.upstream.fields.keepalive_pool.idle_timeout.placeholder':
    'कृपया निष्क्रिय समय सीमा दर्ज करें',
  'component.upstream.fields.keepalive_pool.requests': 'अनुरोध',
  'component.upstream.fields.keepalive_pool.requests.placeholder': 'कृपया अनुरोध दर्ज करें',
  'component.upstream.fields.checks.active.type': 'प्रकार',
  'component.upstream.fields.checks.active.type.tooltip':
    'क्या HTTP या HTTPS का उपयोग करके सक्रिय स्वास्थ्य जाँचें करनी चाहिए, या क्या केवल टीसीपी कनेक्शन का प्रयास करना चाहिए।',
  'component.upstream.fields.checks.active.concurrency': 'समवर्थन',
  'component.upstream.fields.checks.active.concurrency.tooltip':
    'सक्रिय स्वास्थ्य जाँचों में सक्रिय स्वास्थ्य जाँचों में सक्रिय निर्देशित साक्षर संख्या।',
  'component.upstream.fields.checks.active.host': 'होस्ट',
  'component.upstream.fields.checks.active.host.required': 'कृपया होस्टनाम दर्ज करें',
  'component.upstream.fields.checks.active.host.tooltip':
    'सक्रिय स्वास्थ्य जाँच करने के लिए उपयोग किए जाने वाले HTTP अनुरोध का होस्टनाम',
  'component.upstream.fields.checks.active.host.scope': 'केवल अक्षर, संख्या और . का समर्थन किया जाता है',
  'component.upstream.fields.checks.active.port': 'पोर्ट',
  'component.upstream.fields.checks.active.http_path': 'HTTP पथ',
  'component.upstream.fields.checks.active.http_path.tooltip':
    'उद्दीप्ति के लिए HTTP GET अनुरोध जारी करते समय उपयोग करने के लिए पथ। डिफ़ॉल्ट मूल्य है /।',
  'component.upstream.fields.checks.active.http_path.placeholder':
    'कृपया HTTP अनुरोध पथ दर्ज करें',
  'component.upstream.fields.checks.active.https_verify_certificate': 'HTTPs सर्टिफिकेट सत्यापित करें',
  'component.upstream.fields.checks.active.https_verify_certificate.tooltip':
    'क्या सक्रिय स्वास्थ्य जाँच करने के लिए HTTPS का उपयोग करते समय दूरस्त होस्ट के SSL सर्टिफिकेट की मान्यता की जाए।',
  'component.upstream.fields.checks.active.req_headers': 'अनुरोध हैडर्स',
  'component.upstream.fields.checks.active.req_headers.tooltip':
    'अतिरिक्त अनुरोध हैडर्स, उदाहरण: User-Agent: curl/7.29.0',
  'component.upstream.fields.checks.active.healthy.interval': 'अंतराल',
  'component.upstream.fields.checks.active.healthy.interval.tooltip':
    'स्वस्थ लक्ष्यों के लिए स्वास्थ्य जाँचों के लिए अंतराल (सेकंड में)',
  'component.upstream.fields.checks.active.healthy.successes': 'सफलताएँ',
  'component.upstream.fields.checks.active.healthy.successes.tooltip':
    'एक लक्ष्य को स्वस्थ मानने के लिए सफलताओं की संख्या',
  'component.upstream.fields.checks.active.healthy.successes.required':
    'कृपया सफलताओं की संख्या दर्ज करें',
  'component.upstream.fields.checks.active.healthy.http_statuses': 'HTTP स्थितियाँ',
  'component.upstream.fields.checks.active.healthy.http_statuses.tooltip':
    'सक्रिय स्वास्थ्य जाँच में सक्षमता को संकेत करने के लिए HTTP स्थितियों का एक अर्रेय।',
  'component.upstream.fields.checks.active.unhealthy.timeouts': 'समय सीमा से बाहर',
  'component.upstream.fields.checks.active.unhealthy.timeouts.tooltip':
    'सक्रिय प्रोब्स में समय सीमा से बाहर जाने के लिए अस्वस्थ लक्ष्यों की संख्या।',
  'component.upstream.fields.checks.active.unhealthy.http_failures': 'HTTP असफलताएं',
  'component.upstream.fields.checks.active.unhealthy.http_failures.tooltip':
    'एक लक्ष्य को अस्वस्थ मानने के लिए HTTP असफलताओं की संख्या',
  'component.upstream.fields.checks.active.unhealthy.http_failures.required':
    'कृपया HTTP असफलताओं की संख्या दर्ज करें',
  'component.upstream.fields.checks.active.unhealthy.tcp_failures': 'TCP असफलताएं',
  'component.upstream.fields.checks.active.unhealthy.tcp_failures.tooltip':
    'एक लक्ष्य को अस्वस्थ मानने के लिए TCP असफलताओं की संख्या',
  'component.upstream.fields.checks.active.unhealthy.tcp_failures.required':
    'कृपया TCP असफलताओं की संख्या दर्ज करें',
  'component.upstream.fields.checks.active.unhealthy.interval': 'अंतराल',
  'component.upstream.fields.checks.active.unhealthy.interval.tooltip':
    'अस्वस्थ लक्ष्यों के लिए सक्रिय स्वास्थ्य जाँचों के लिए अंतराल (सेकंड में)। शून्य का मतलब है कि स्वस्थ लक्ष्यों के लिए सक्रिय प्रोब्स को किया जाना चाहिए नहीं है।',
  'component.upstream.fields.checks.active.unhealthy.required':
    'कृपया अस्वस्थ अंतराल दर्ज करें',
  'component.upstream.fields.checks.passive.healthy.successes': 'सफलताएँ',
  'component.upstream.fields.checks.passive.healthy.successes.tooltip':
    'एक लक्ष्य को स्वस्थ मानने के लिए सफलताओं की संख्या',
  'component.upstream.fields.checks.passive.healthy.successes.required':
    'कृपया सफलताओं की संख्या दर्ज करें',
  'component.upstream.fields.checks.passive.unhealthy.timeouts': 'समय सीमा से बाहर',
  'component.upstream.fields.checks.passive.unhealthy.timeouts.tooltip':
    'सक्रिय स्वास्थ्य जाँचों द्वारा देखे जाने वाले अस्वस्थ लक्ष्यों को मान्यता देने के लिए प्रॉक्सी ट्रैफिक में समय सीमा से बाहर जाने वाले टाइमआउट्स की संख्या।',
  'component.upstream.other.none': 'कोई नहीं (सेवा को बाइंड करते समय ही उपलब्ध है)',
  'component.upstream.other.pass_host-with-multiple-nodes.title':
    'कृपया लक्ष्य नोड कॉन्फ़िगरेशन की जाँच करें',
  'component.upstream.other.pass_host-with-multiple-nodes':
    'लक्ष्य नोड सूची में होस्ट नाम या आईपी का उपयोग करते समय, सुनिश्चित करें कि केवल एक लक्ष्य नोड है',
  'component.upstream.other.health-check.passive-only':
    'जब पैसिव स्वास्थ्य जाँच सक्षम है, तब सक्रिय स्वास्थ्य जाँच को भी सक्षम करना आवश्यक है।',
  'component.upstream.other.health-check.invalid': 'कृपया स्वास्थ्य जाँच कॉन्फ़िगरेशन की जाँच करें',
};