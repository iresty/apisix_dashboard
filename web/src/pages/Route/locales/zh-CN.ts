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
  // global
  'page.route.parameterPosition': '参数位置',
  'page.route.httpRequestHeader': 'HTTP 请求头',
  'page.route.requestParameter': '请求参数',
  'page.route.postRequestParameter': 'POST 请求参数',
  'page.route.builtinParameter': '内置参数',
  'page.route.parameterName': '参数名称',
  'page.route.operationalCharacter': '运算符',
  'page.route.equal': '等于（==）',
  'page.route.unequal': '不等于（~=）',
  'page.route.greaterThan': '大于（>）',
  'page.route.lessThan': '小于（<）',
  'page.route.regexMatch': '正则匹配（～～）',
  'page.route.caseInsensitiveRegexMatch': '不区分大小写的正则匹配（~*）',
  'page.route.in': 'IN',
  'page.route.has': 'HAS',
  'page.route.reverse': '非(!)',
  'page.route.rule': '规则',
  'page.route.host': '域名',
  'page.route.path': '路径',
  'page.route.remoteAddrs': '客户端地址',
  'page.route.value': '参数值',
  'page.route.httpHeaderName': 'HTTP 请求头名称',
  'page.route.status': '状态',
  'page.route.groupName': '分组名称',
  'page.route.offline': '下线',
  'page.route.publish': '发布',
  'page.route.published': '已发布',
  'page.route.unpublished': '未发布',
  'page.route.onlineDebug': '在线调试',
  'page.route.pluginTemplateConfig': '插件模版配置',
  'page.route.service': '绑定服务',
  'page.route.instructions': '说明',
  'page.route.import': '导入',
  'page.route.createRoute': '创建路由',
  'page.route.editRoute': '编辑路由',

  // button
  'page.route.button.returnList': '返回路由列表',
  'page.route.button.send': '发送请求',

  // input
  'page.route.input.placeholder.parameterNameHttpHeader': '请求头名称，例如：HOST',
  'page.route.input.placeholder.parameterNameRequestParameter': '参数名称，例如：id',
  'page.route.input.placeholder.redirectCustom': '例如：/foo/index.html',
  'page.route.input.placeholder.requestUrl': '请输入合法的请求地址',
  'page.route.input.placeholder.paramKey': '参数名称',
  'page.route.input.placeholder.paramType': '参数类型',
  'page.route.input.placeholder.paramValue': '参数值',
  // form
  'page.route.form.itemRulesRequiredMessage.parameterName': '仅支持字母和数字，且只能以字母开头',
  'page.route.form.itemRulesPatternMessage.apiNameRule': '最大长度应仅为 100',
  'page.route.form.itemLabel.httpMethod': 'HTTP 方法',
  'page.route.form.itemLabel.scheme': '协议',
  'page.route.form.itemLabel.priority': '优先级',
  'page.route.form.itemLabel.redirect': '重定向',
  'page.route.form.itemLabel.redirectCustom': '自定义重定向',
  'page.route.form.itemLabel.URIRewriteType': '路径改写',
  'page.route.form.itemLabel.hostRewriteType': '域名改写',
  'page.route.form.itemLabel.methodRewrite': 'HTTP 方法改写',
  'page.route.form.itemLabel.headerRewrite': '请求头改写',
  'page.route.form.itemLabel.redirectURI': '重定向路径',
  'page.route.form.itemExtraMessage.domain': '路由匹配的域名列表。支持泛域名，如：*.test.com',
  'page.route.form.itemRulesPatternMessage.domain':
    '仅支持字母、数字、-、_和 *，但 * 需要在开头位置。',
  'page.route.form.itemExtraMessage1.path':
    'HTTP 请求路径，如 /foo/index.html，支持请求路径前缀 /foo/*。/* 代表所有路径',
  'page.route.form.itemExtraMessage1.remoteAddrs':
    '客户端与服务器握手时 IP，即客户端 IP，例如：192.168.1.101，192.168.1.0/24，::1，fe80::1，fe80::1/64',
  'page.route.form.itemRulesPatternMessage.remoteAddrs':
    '请输入合法的 IP 地址，例如：192.168.1.101，192.168.1.0/24，::1，fe80::1，fe80::1/64',
  'page.route.form.itemLabel.username': '用户名',
  'page.route.form.itemLabel.password': '密 码',
  'page.route.form.itemLabel.token': 'Token',
  'page.route.form.itemLabel.apikey': 'Apikey',

  // select
  'page.route.select.option.enableHttps': '启用 HTTPS',
  'page.route.select.option.configCustom': '自定义',
  'page.route.select.option.forbidden': '禁用',
  'page.route.select.option.redirect301': '301（永久重定向）',
  'page.route.select.option.redirect302': '302（临时重定向）',
  'page.route.select.option.inputManually': '手动填写',
  'page.route.select.option.methodRewriteNone': '不改写',

  // steps
  'page.route.steps.stepTitle.defineApiRequest': '设置路由信息',
  'page.route.steps.stepTitle.defineApiBackendServe': '设置上游服务',

  // panelSection
  'page.route.panelSection.title.nameDescription': '基本信息',
  'page.route.panelSection.title.httpOverrideRequestHeader': 'HTTP 请求头改写',
  'page.route.panelSection.title.requestOverride': '请求改写',
  'page.route.panelSection.title.requestConfigBasicDefine': '匹配条件',
  'page.route.panelSection.title.advancedMatchRule': '高级匹配条件',
  'page.route.PanelSection.title.requestLine': '请求行定义',
  'page.route.PanelSection.title.defineRequestParams': '请求参数定义',
  'page.route.PanelSection.title.responseResult': '请求响应结果',
  'page.route.Model.inputPath': '请重写通配符*的路径',
  'page.route.Select.method': '请求方法',
  'page.route.Select.protocol': '协议',
  'page.route.Select.requestPath': '选择请求路径',

  'page.route.httpAction': '行为',
  'page.route.httpOverrideOrCreate': '重写/创建',

  'page.route.form.itemLabel.domainNameOrIp': '域名/IP',
  'page.route.form.itemExtraMessage.domainNameOrIp': '使用域名时，默认解析本地：/etc/resolv.conf',
  'page.route.portNumber': '端口',
  'page.route.weight': '权重',

  'page.route.radio.static': '静态改写',
  'page.route.radio.regex': '正则改写',
  // Need a better translation
  'page.route.form.itemLabel.regex': '匹配正则表达式',
  'page.route.form.itemLabel.template': '转发路径模版',
  'page.route.form.itemHelp.status': '路由创建后是否可以使用， 默认值为 false',
  'page.route.radio.staySame': '保持原样',
  'page.route.input.placeholder.newPath': '例如：/foo/bar/index.html',
  'page.route.form.itemLabel.newPath': '新路径',
  'page.route.form.itemLabel.newHost': '新域名',
  'page.route.popconfirm.title.offline': '确定下线该路由吗？',
  'page.route.debug.showResultAfterSendRequest': '发送请求后在此查看响应结果',
  'page.route.TabPane.queryParams': '查询参数',
  'page.route.TabPane.bodyParams': '请求体参数',
  'page.route.TabPane.headerParams': '请求头参数',
  'page.route.TabPane.authentication': '认证',
  'page.route.TabPane.response': '响应结果',
  'page.route.TabPane.header': '响应请求头参数',
  'page.route.debugWithoutAuth': '当前请求不启用任何认证方式。',
  'page.route.button.exportOpenApi': '导出 OpenAPI',
  'page.route.exportRoutesTips': '请选择导出文件的类型',
  'page.route.button.importOpenApi': '导入 OpenAPI',
  'page.route.button.selectFile': '请选择上传文件',
  'page.route.list': '路由列表',
  'page.route.tooltip.pluginOrchOnlySupportChrome': '插件编排仅支持 Chrome 浏览器。',
  'page.route.tooltip.pluginOrchWithoutProxyRewrite':
    '当步骤一中 配置了 请求改写时，不可使用插件编排模式。',
  'page.route.tooltip.pluginOrchWithoutRedirect':
    '当步骤一中 重定向 选择为 启用 HTTPS 时，不可使用插件编排模式。',

  'page.route.tabs.normalMode': '普通模式',
  'page.route.tabs.orchestration': '编排模式',

  'page.route.list.description':
    '路由（Route）是请求的入口点，它定义了客户端请求与服务之间的匹配规则。路由可以与服务（Service）、上游（Upstream）关联，一个服务可对应一组路由，一个路由可以对应一个上游对象（一组后端服务节点），因此，每个匹配到路由的请求将被网关代理到路由绑定的上游服务中。',

  'page.route.configuration.name.rules.required.description': '请输入路由名称',
  'page.route.configuration.name.placeholder': '请输入路由名称',
  'page.route.configuration.desc.tooltip': '路由的描述信息',
  'page.route.configuration.publish.tooltip': '用于控制路由创建后，是否立即发布到网关',
  'page.route.configuration.version.placeholder': '请输入路由版本号',
  'page.route.configuration.version.tooltip': '路由的版本号，如 V1',
  'page.route.configuration.normal-labels.tooltip': '为路由增加自定义标签，可用于路由分组。',

  'page.route.configuration.path.rules.required.description': '请输入有效的 HTTP 请求路径',
  'page.route.configuration.path.placeholder': '请输入 HTTP 请求路径',
  'page.route.configuration.remote_addrs.placeholder': '请输入客户端地址',
  'page.route.configuration.host.placeholder': '请输入 HTTP 请求域名',

  'page.route.service.none': '不绑定服务',

  'page.route.rule.create': '创建规则',
  'page.route.rule.edit': '编辑规则',

  'page.route.advanced-match.operator.sample.IN': '请输入数组，示例：["1", "2"]',
  'page.route.advanced-match.operator.sample.~~': '请输入正则表达式，示例：[a-z]+',
  'page.route.fields.service_id.invalid': '请检查路由绑定的服务',
  'page.route.fields.service_id.without-upstream': '如果不绑定服务，则必须设置上游服务（步骤 2）',
  'page.route.advanced-match.tooltip':
    '支持通过请求头，请求参数、Cookie 进行路由匹配，可应用于灰度发布，蓝绿测试等场景。',
  'page.route.advanced-match.message': '提示',
  'page.route.advanced-match.tips.requestParameter': '请求参数：请求 URL 中的 Query 部分',
  'page.route.advanced-match.tips.postRequestParameter':
    'POST 请求参数：仅支持 POST x-www-form-urlencoded 表单',
  'page.route.advanced-match.tips.builtinParameter': '内置参数：Nginx 内部参数',

  'page.route.fields.custom.redirectOption.tooltip': '在此配置 redirect 插件',
  'page.route.fields.service_id.tooltip': '绑定服务（Service）对象，以便复用其中的配置。',

  'page.route.fields.vars.invalid': '请检查高级匹配条件配置',
  'page.route.fields.vars.in.invalid': '使用 IN 操作符时，请输入数组格式的参数值。',
};
