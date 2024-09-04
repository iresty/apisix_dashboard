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
  'component.plugin.tip1': 'नोट: प्लगइन को कस्टमाइज़ करने के बाद, आपको schema.json को अपडेट करने की आवश्यकता है।',
  'component.plugin.tip2': 'अपडेट कैसे करें?',
  'component.select.pluginTemplate': 'प्लगइन टेम्पलेट का चयन करें',
  'component.step.select.pluginTemplate.select.option': 'कस्टम',
  'component.plugin.pluginTemplate.tip1':
    '1. जब किसी रूट में पहले से ही प्लगइन्स फ़ील्ड को कॉन्फ़िगर किया जा चुका है, तो प्लगइन टेम्पलेट में मौजूद प्लगइन्स इसमें मर्ज हो जाएंगे।',
  'component.plugin.pluginTemplate.tip2':
    '2. प्लगइन टेम्पलेट में समान प्लगइन उसमें मौजूद प्लगइन को अधिलेखित करेगा।',
  'component.plugin.enable': 'सक्षम करें',
  'component.plugin.disable': 'संपादित करें',
  'component.plugin.authentication': 'प्रमाणीकरण',
  'component.plugin.security': 'सुरक्षा',
  'component.plugin.traffic': 'ट्रैफ़िक नियंत्रण',
  'component.plugin.serverless': 'सर्वरलेस',
  'component.plugin.observability': 'अवलोकन',
  'component.plugin.other': 'अन्य',
  'component.plugin.all': 'सभी',
  // cors
  'component.pluginForm.cors.allow_origins.tooltip':
    'कौन से उत्सर्ग को सक्षम करने की अनुमति है, प्रारूप का रूप में: योजना://होस्ट:पोर्ट, उदाहरण के लिए: https://somehost.com:8081। एक से अधिक उत्सर्ग का उपयोग करने के लिए , का उपयोग करें। जब allow_credential नहीं है, तो आप * का उपयोग कर सकते हैं ताकि किसी भी उत्सर्ग को सक्षम किया जा सके। आप ** का उपयोग भी कर सकते हैं जो पहले से ही allow_credential को सक्षम करने के बावजूद सभी किसी भी उत्सर्ग को सक्षम कर सकता है, लेकिन यह कुछ सुरक्षा जोखिम लेकर आ सकता है।',
  'component.pluginForm.cors.allow_origins.extra': 'उदाहरण के लिए: https://somehost.com:8081',
  'component.pluginForm.cors.allow_methods.tooltip':
    'कौन से विधि को सक्षम करने की अनुमति है, जैसे: GET, POST इत्यादि। एक से अधिक विधि का उपयोग करने के लिए , का उपयोग करें। जब allow_credential नहीं है, तो आप * का उपयोग कर सकते हैं ताकि किसी भी विधि को सक्षम किया जा सके। आप ** का उपयोग भी कर सकते हैं जो पहले से ही allow_credential को सक्षम करने के बावजूद सभी किसी भी विधि को सक्षम कर सकता है, लेकिन यह कुछ सुरक्षा जोखिम लेकर आ सकता है।',
  'component.pluginForm.cors.allow_headers.tooltip':
    'जब क्रॉस-उत्सर्ग संसाधन तक पहुँचने के लिए अनुरोध करता है, तो कौन-कौन से हैडर सेट करने की अनुमति है। कई मूल्यों का उपयोग करने के लिए , का उपयोग करें। जब allow_credential नहीं है, तो आप * का उपयोग कर सकते हैं ताकि किसी भी अनुरोध हैडर को सक्षम किया जा सके। आप ** का उपयोग भी कर सकते हैं जो पहले से ही allow_credential को सक्षम करने के बावजूद सभी किसी भी हैडर को सक्षम कर सकता है, लेकिन यह कुछ सुरक्षा जोखिम लेकर आ सकता है।',
  'component.pluginForm.cors.expose_headers.tooltip':
    'जब क्रॉस-उत्सर्ग संसाधन तक पहुँचने के लिए प्रतिसाद में सेट करने की अनुमति है, तो कौन-कौन से हैडर सेट करने की अनुमति है। कई मूल्यों का उपयोग करने के लिए , का उपयोग करें। जब allow_credential नहीं है, तो आप * का उपयोग कर सकते हैं ताकि किसी भी हैडर को सक्षम किया जा सके। आप ** का उपयोग भी कर सकते हैं जो पहले से ही allow_credential को सक्षम करने के बावजूद सभी किसी भी हैडर को सक्षम कर सकता है, लेकिन यह कुछ सुरक्षा जोखिम लेकर आ सकता है।',
  'component.pluginForm.cors.max_age.tooltip':
    'परिणामों को कैश करने के लिए अधिकतम सेकंड की संख्या। इस समय सीमा के भीतर, ब्राउज़र पिछले चेक के परिणाम का पुनः उपयोग करेगा। -1 का मतलब है कोई कैश नहीं किया जाएगा। कृपया ध्यान दें कि अधिकतम मूल्य ब्राउज़र पर निर्भर करता है, कृपया विवरण के लिए MDN का संदर्भ दें।',
  'component.pluginForm.cors.allow_credential.tooltip':
    'यदि आप इस विकल्प को सत्य में सेट करते हैं, तो आप अन्य विकल्पों के लिए * का उपयोग नहीं कर सकते हैं।',
  'component.pluginForm.cors.allow_origins_by_metadata.tooltip':
    'यहाँ सूचीबद्ध allow_origins का संदर्भ करके कौन-सा उत्सर्ग सक्षम है, उसे मिलता है।',
  'component.pluginForm.cors.allow_origins_by_regex.tooltip':
    'कौन सा उत्सर्ग सक्षम है, इसे सक्षम करने के लिए रेजेक्ट अभिव्यक्तियों का उपयोग करें। प्रत्येक इनपुट बॉक्स को केवल एक, स्थानीय रेगुलर एक्सप्रेशन के साथ कॉन्फ़िगर किया जा सकता है, जैसे ".*.test.com" जो test.com के किसी भी सबडोमेन के साथ मेल खाता है।',
  // referer-restriction
  'component.pluginForm.referer-restriction.whitelist.tooltip':
    'व्हाइटलिस्ट करने के लिए होस्टनाम की सूची। होस्टनाम * के साथ शुरू हो सकता है जैसा कि वाइल्डकार्ड के रूप में।',
  'component.pluginForm.referer-restriction.blacklist.tooltip':
    'ब्लैकलिस्ट करने के लिए होस्टनाम की सूची। होस्टनाम * के साथ शुरू हो सकता है जैसा कि वाइल्डकार्ड के रूप में।',
  'component.pluginForm.referer-restriction.listEmpty.tooltip': 'सूची खाली है',
  'component.pluginForm.referer-restriction.bypass_missing.tooltip':
    'क्या रेफ़रर हेडर गायब या रूपांतरित है जब नियमित रूप से चेक करने पर जाने की आवश्यकता है।',
  'component.pluginForm.referer-restriction.message.tooltip':
    'यदि पहुँच नहीं होने की अनुमति नहीं है, तो संदेश वापस दिया जाएगा।',
  // api-breaker
  'component.pluginForm.api-breaker.break_response_code.tooltip':
    'अशान्त होने पर त्रुटि कोड वापस दें।',
  'component.pluginForm.api-breaker.break_response_body.tooltip':
    'जब अपस्ट्रीम अशान्त होता है, तो प्रतिसाद संदेश का शरीर वापस दें।',
  'component.pluginForm.api-breaker.break_response_headers.tooltip':
    'अशान्त होने पर हेडर वापस दें।',
  'component.pluginForm.api-breaker.max_breaker_sec.tooltip': 'अधिकतम ब्रेकर समय (सेकंड)।',
  'component.pluginForm.api-breaker.unhealthy.http_statuses.tooltip':
    'अशान्त होने पर स्थिति कोड।',
  'component.pluginForm.api-breaker.unhealthy.failures.tooltip':
    'जब अशान्त स्थिति को ट्रिगर करने वाले सतत त्रुटि अनुरोधों की संख्या।',
  'component.pluginForm.api-breaker.healthy.http_statuses.tooltip': 'स्वस्थ होने पर स्थिति कोड।',
  'component.pluginForm.api-breaker.healthy.successes.tooltip':
    'स्वस्थ स्थिति को ट्रिगर करने वाले सतत सामान्य अनुरोधों की संख्या।',
  // proxy-mirror
  'component.pluginForm.proxy-mirror.host.tooltip':
    'मिरर सेवा पता को निर्दिष्ट करें, उदाहरण: http://127.0.0.1:9797 (पता में योजना शामिल होनी चाहिए: http या https, URI हिस्सा नहीं)',
  'component.pluginForm.proxy-mirror.host.extra': 'उदाहरण: http://127.0.0.1:9797',
  'component.pluginForm.proxy-mirror.host.ruletip':
    'पता में योजना शामिल होनी चाहिए: http या https, URI हिस्सा नहीं',
  'component.pluginForm.proxy-mirror.path.tooltip':
    'मिरर अनुरोध का पथ हिस्सा को निर्दिष्ट करें। इसके बिना, वर्तमान पथ का उपयोग किया जाएगा।',
  'component.pluginForm.proxy-mirror.path.ruletip': 'कृपया सही पथ दर्ज करें, उदाहरण के लिए /पथ',
  'component.pluginForm.proxy-mirror.sample_ratio.tooltip':
    'वे अनुरोध जिन्हें मिरर किया जाएगा, उनका नमूना अनुपात।',
// limit-conn
  'component.pluginForm.limit-conn.conn.tooltip':
    'समय की सीमा तक अनुमत प्रतिसादी अनुरोधों की अधिकतम संख्या। इस अनुपालन के लिए (और conn + burst के बीच) बढ़ी हुई अनुरोध हेतु (latency सेकंड डिफ़ॉल्ट_conn_delay द्वारा कॉन्फ़िगर किया जाता है) को विलम्बित किया जाएगा।',
  'component.pluginForm.limit-conn.burst.tooltip':
    'विलंबित करने की अनुमति दी जाने वाली अधिशेष प्रतिसादी अनुरोधों (या कनेक्शन) की संख्या।',
  'component.pluginForm.limit-conn.default_conn_delay.tooltip':
    'जब प्रतिसादी अनुरोध conn से अधिक होते हैं, लेकिन (conn + burst) के बीच, तो अनुरोध की लैटेंसी सेकंड।',
  'component.pluginForm.limit-conn.key_type.tooltip':
    'कुंजी प्रकार, समर्थन: "var" (एकल var) और "var_combination" (var को संयुक्त करें)',
  'component.pluginForm.limit-conn.key.tooltip':
    'सीमा तक पहुँच को सीमित करने के लिए। उदाहरण के लिए, कोई यूज़ होस्ट का नाम (या सर्वर क्षेत्र) का उपयोग कर सकता है, ताकि हम होस्ट के नाम के प्रति सीमा लगा सकते हैं। अन्यथा, हम क्लाइयंट पता का उपयोग कर सकते हैं, ताकि हम एकल क्लाइयंट से हमारी सेवा को बहुतल संबंधित अनुरोध या कनेक्शनों से बचा सकते हैं।।',
  'component.pluginForm.limit-conn.rejected_code.tooltip':
    'जब अनुरोध conn + burst से अधिक होता है, तो यह अस्वीकृत किया जाएगा।',
  'component.pluginForm.limit-conn.rejected_msg.tooltip':
    'जब अनुरोध conn + burst से अधिक होता है, तो इसे अस्वीकृत किया जाएगा, तो यह प्रतिसाद शरीर।',
  'component.pluginForm.limit-conn.only_use_default_delay.tooltip':
    'सेकंड की शीघ्रता का कड़ा मोड सक्षम करें। यदि आप इस विकल्प को सत्य में सेट करते हैं, तो यह आपके द्वारा निर्धारित किए गए सेकंड के अनुसार कड़ा से चलेगा बिना कोई अतिरिक्त हिसाब तालिका होगा।',
  'component.pluginForm.limit-conn.allow_degradation.tooltip':
    'क्या plugin degradation को सक्षम करने के लिए अस्थायी रूप से अनुपलब्ध है जब limit-conn कार्य की अवस्था। यह मूल्य सत्य है, तो अनुरोध जारी रहने दिया जाएगा, डिफ़ॉल्ट फॉल्स है।',
  // limit-req
  'component.pluginForm.limit-req.rate.tooltip':
    'निर्दिष्ट अनुरोध दर (प्रति सेकंड संख्या) सीमा। इस दर से अधिक (और burst के बीच) अनुरोध हेतु विलंबित किया जाएगा ताकि दर से अनुपालन किया जा सके।',
  'component.pluginForm.limit-req.burst.tooltip':
    'विलंबित करने की अनुमति दी जाने वाली प्रति सेकंड अधिशेष अनुरोधों की संख्या। इस कड़ी सीमा को अधिग्रहण करने वाले अनुरोधों को तुरंत अस्वीकृत किया जाएगा।',
  'component.pluginForm.limit-req.key_type.tooltip':
    'कुंजी प्रकार, समर्थन: "var" (एकल var) और "var_combination" (var को संयुक्त करें)',
  'component.pluginForm.limit-req.key.tooltip': 'दर को सीमित करने के लिए उपयोगकर्ता द्वारा निर्दिष्ट कुंजी।',
  'component.pluginForm.limit-req.rejected_code.tooltip':
    'जब अनुरोध सीमा से अधिक होता है, तो इसे अस्वीकृत किया जाएगा, तब यह HTTP स्थिति कोड।',
  'component.pluginForm.limit-req.rejected_msg.tooltip':
    'जब अनुरोध सीमा से अधिक होता है, तो इसे अस्वीकृत किया जाएगा, तब यह प्रतिसाद शरीर।',
  'component.pluginForm.limit-req.nodelay.tooltip':
    'यदि nodelay फ्लैग सत्य है, तो बर्स्ट हुए अनुरोधों को विलंबित नहीं किया जाएगा',
  'component.plugin.form': 'फ़ॉर्म',
  'component.plugin.format-codes.disable': 'JSON या YAML डेटा को स्वरूपित करें',
  'component.plugin.editor': 'प्लगइन संपादक',
  'component.plugin.noConfigurationRequired': "कॉन्फ़िगरेशन की आवश्यकता नहीं है",
  // limit-count
  'component.pluginForm.limit-count.count.tooltip': 'निर्दिष्ट अनुरोध सीमा।',
  'component.pluginForm.limit-count.time_window.tooltip':
    'अनुरोध गणना को रीसेट करने से पहले समय विंडो सेकंड में।',
  'component.pluginForm.limit-count.key_type.tooltip':
    'कुंजी प्रकार, समर्थन: "var" (एकल var) और "var_combination" (var को संयुक्त करें)',
  'component.pluginForm.limit-count.key.tooltip': 'दर को सीमित करने के लिए उपयोगकर्ता द्वारा निर्दिष्ट कुंजी।',
  'component.pluginForm.limit-count.rejected_code.tooltip':
    'जब अनुरोध सीमा से अधिक होता है, तो इसे अस्वीकृत किया जाएगा, डिफ़ॉल्ट 503।',
  'component.pluginForm.limit-count.rejected_msg.tooltip':
    'जब अनुरोध सीमा से अधिक होता है, तो इसे अस्वीकृत किया जाएगा, डिफ़ॉल्ट 503।',
  'component.pluginForm.limit-count.policy.tooltip':
    'सीमा और बढ़ाने के लिए उपयोग करने के लिए रेट-लिमिटिंग नीतियां। उपलब्ध मूल्य हैं local(गिनती को स्थानीय रूप में नोड में मेमोरी में स्थानीय रूप से स्थित रहेगी) और redis(गिनती को एक रेडिस सर्वर पर स्थानीय रूप से स्थित किया जाएगा और नोडों के बीच साझा किया जाएगा, सामान्यत: इसका उपयोग वैश्विक गति सीमा करने के लिए करें) और redis-cluster(रेडिस क्लस्टर पैटर्न का उपयोग करने वाली एक सार्वजनिक स्थानीय रूप, समर्थन मोड के लिए हैं (केवल रेडिस क्लस्टर मोड के लिए उपयोग करें)।',
  'component.pluginForm.limit-count.allow_degradation.tooltip':
    'क्या plugin degradation को सक्षम करने के लिए अस्थायी रूप से अनुपलब्ध है (उदाहरण के लिए, रेडिस समय समाप्त हो गया है)। जब मूल्य को सत्य में सेट किया जाता है, तो अनुरोध जारी रहने दिया जाएगा।',
  'component.pluginForm.limit-count.show_limit_quota_header.tooltip':
    'क्या प्रतिसाद मुख्य में X-RateLimit-Limit और X-RateLimit-Remaining (जिससे मतलब है कि कुल अनुरोधों की संख्या और अब तक भेजे जा सकने वाले अनुरोधों की बची हुई संख्या) को जवाब है',
  'component.pluginForm.limit-count.group.tooltip':
    'उन रूटें जो एक ही समूह के साथ कॉन्फ़िगर हैं, वे एक ही काउंटर को साझा करेंगी',
  'component.pluginForm.limit-count.redis_host.tooltip':
    'रेडिस नीति का उपयोग करते समय, इस संपत्ति ने रेडिस सर्वर का पता दिया है।',
  'component.pluginForm.limit-count.redis_port.tooltip':
    'रेडिस नीति का उपयोग करते समय, इस संपत्ति ने रेडिस सर्वर का पोर्ट निर्दिष्ट किया है।',
  'component.pluginForm.limit-count.redis_password.tooltip':
    'रेडिस नीति का उपयोग करते समय, इस संपत्ति ने रेडिस सर्वर का पासवर्ड निर्दिष्ट किया है।',
  'component.pluginForm.limit-count.redis_database.tooltip':
    'रेडिस नीति का उपयोग करते समय, इस संपत्ति ने रेडिस सर्वर के चयनित डेटाबेस को निर्दिष्ट किया है, और यह केवल गैर रेडिस क्लस्टर मोड के लिए है (एकल इंस्टेंस मोड या रेडिस सार्वजनिक क्लाउड सेवा जो एकल प्रवेश प्रदान करती है)।',
  'component.pluginForm.limit-count.redis_timeout.tooltip':
    'रेडिस नीति का उपयोग करते समय, इस संपत्ति ने रेडिस सर्वर को सबमिट किए गए किसी भी कमांड के लिए समय सीमा निर्दिष्ट किया है।',
  'component.pluginForm.limit-count.redis_cluster_nodes.tooltip':
    'रेडिस-क्लस्टर नीति का उपयोग करते समय, इस संपत्ति में रेडिस क्लस्टर सेवा नोड का पता है (कम से कम दो नोड)।',
  'component.pluginForm.limit-count.redis_cluster_name.tooltip':
    'रेडिस-क्लस्टर नीति का उपयोग करते समय, इस संपत्ति में रेडिस क्लस्टर सेवा नोड का नाम है।',
  'component.pluginForm.limit-count.atLeast2Characters.rule': 'कृपया कम से कम 2 अक्षर दर्ज करें।',
};