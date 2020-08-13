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
  'route.create.define.api.request': '定义 API 请求',
  'route.create.define.api.backend.server': '定义 API 后端服务',
  'route.create.plugin.configuration': '插件配置',

  'route.result.submit.success': '提交成功',
  'route.result.return.list': '返回路由列表',
  'route.result.create': '创建路由',

  'route.match.parameter.position': '参数位置',
  'route.match.http.request.header': 'HTTP 请求头',
  'route.match.request.parameter': '请求参数',
  'route.match.parameter.name': '参数名称',
  'route.match.operational.character': '运算符',
  'route.match.equal': '等于',
  'route.match.unequal': '不等于',
  'route.match.greater.than': '大于',
  'route.match.less.than': '小于',
  'route.match.regex.match': '正则匹配',
  'route.match.parameter.value': '参数值',
  'route.match.operation': '操作',
  'route.match.edit': '编辑',
  'route.match.delete': '删除',
  'route.match.edit.rule': '编辑规则',
  'route.match.create.rule': '创建规则',
  'route.match.confirm': '确定',
  'route.match.cancel': '取消',
  'route.match.select.parameter.position': '请选择参数位置',
  'route.match.request.header.example': '请求头键名，例如：HOST',
  'route.match.parameter.name.example': '参数名称，例如：id',
  'route.match.input.parameter.name': '请输入参数名称',
  'route.match.parameter.name.rule': '仅支持字母、数字、- 和 _ ，且只能以字母开头',
  'route.match.rule': '仅支持字母和数字，且只能以字母开头',
  'route.match.choose.operational.character': '请选择运算符',
  'route.match.value': '值',
  'route.match.input.parameter.value': '请输入参数值',
  'route.match.advanced.match.rule': '高级路由匹配条件',
  'route.match.create': '创建',

  'route.meta.name.description': '名称及其描述',
  'route.meta.api.name': 'API 名称',
  'route.meta.input.api.name': '请输入 API 名称',
  'route.meta.api.name.rule': '最大长度100，仅支持字母、数字、- 和 _，且只能以字母开头',
  'rotue.meta.api.rule': '仅支持字母、数字、- 和 _，且只能以字母开头',
  'route.meta.description': '描述',
  'route.meta.description.rule': '不超过 200 个字符',

  'route.request.config.domain.name': '域名',
  'route.request.config.domain.or.ip': '域名或IP，支持泛域名，如：*.test.com',
  'route.request.config.input.domain.name': '请输入域名',
  'route.request.config.domain.name.rule': '仅支持字母、数字和 * ，且 * 只能是在开头，支持单个 * ',
  'route.request.config.create': '创建',
  'route.request.config.path': '路径',
  'route.request.config.path.description1':
    '1. 请求路径，如 /foo/index.html，支持请求路径前缀 /foo/* ；',
  'route.request.config.path.description2': '2. /* 代表所有路径',
  'route.request.config.input.path': '请输入请求路径',
  'route.request.config.path.rule': '以 / 开头，且 * 只能在最后',
  'route.request.config.basic.define': '请求基础定义',
  'route.request.config.protocol': '协议',
  'route.request.config.choose.protocol': '请选择协议',
  'route.request.config.http.method': 'HTTP 方法',
  'route.request.config.choose.http.method': '请选择 HTTP 方法',
  'route.request.config.redirect': '重定向',
  'route.request.config.enable.https': '启用 HTTPS',
  'route.request.config.custom': '自定义',
  'route.request.config.forbidden': '禁用',
  'route.request.config.redirect.custom': '自定义重定向',
  'route.request.config.redirect.custom.example': '例如：/foo/index.html',
  'route.request.config.redirect.301': '301（永久重定向）',
  'route.request.config.redirect.302': '302（临时重定向）',

  'route.http.request.header.name': 'HTTP 请求头名称',
  'route.http.action': '行为',
  'route.http.override.or.create': '重写/创建',
  'route.http.delete': '删除',
  'route.http.value': '值',
  'route.http.operation': '操作',
  'route.http.edit': '编辑',
  'route.http.edit.request.header': '编辑请求头',
  'route.http.operate.request.header': '操作请求头',
  'route.http.confirm': '确定',
  'route.http.cancel': '取消',
  'route.http.input.request.header.name': '请输入 HTTP 请求头名称',
  'route.http.select.actions': '请选择行为',
  'route.http.input.value': '请输入值',
  'route.http.override.request.header': 'HTTP 请求头改写',

  'route.request.override.input': '手动填写',
  'route.request.override.domain.name.or.ip': '域名/IP',
  'route.request.override.use.domain.name.default.analysis':
    '使用域名时，默认解析本地：/etc/resolv.conf',
  'route.request.override.input.domain.or.ip': '请输入域名/IP',
  'route.request.override.domain.or.ip.rules': '仅支持字母、数字和 . ',
  'route.request.override.input.port.number': '请输入端口号',
  'route.request.override.port.number': '端口号',
  'route.request.override.input.weight': '请输入权重',
  'route.request.override.weight': '权重',
  'route.request.override.create': '创建',
  'route.request.override': '请求改写',
  'route.request.override.protocol': '协议',
  'route.request.override.select.protocol': '请选择协议',
  'route.request.override.stay.same': '保持原样',
  'route.request.override.path': '请求路径',
  'route.request.override.edit': '编辑',
  'route.request.override.new.path': '新路径',
  'route.request.override.input.path': '请输入请求路径',
  'route.request.override.path.example': '例如：/foo/bar/index.html',
  'route.request.override.upstream': '上游',
  'route.request.override.connection.timeout': '连接超时',
  'route.request.override.input.connection.timeout': '请输入连接超时时间',
  'route.request.override.send.timeout': '发送超时',
  'route.request.override.inout.send.timeout': '请输入发送超时时间',
  'route.request.override.receive.timeout': '接收超时',
  'route.request.override.inout.receive.timeout': '请输入接收超时时间',

  'route.constants.define.api.request': '定义 API 请求',
  'route.constants.preview': '预览',
  'route.constants.define.api.backend.serve': '定义 API 后端服务',
  'route.constants.plugin.configuration': '插件配置',

  'route.create.management': '路由管理',

  'route.list.name': '名称',
  'route.list.domain.name': '域名',
  'route.list.path': '路径',
  'route.list.description': '描述',
  'route.list.edit.time': '编辑时间',
  'route.list.operation': '操作',
  'route.list.edit': '编辑',
  'route.list.delete.confrim': '确定删除该路由吗？',
  'route.list.delete.success': '删除成功！',
  'route.list.confirm': '确认',
  'route.list.cancel': '取消',
  'route.list.delete': '删除',
  'route.list': '路由列表',
  'route.list.input': '请输入',
  'route.list.create': '创建',
  'page.route.radio.static': '静态重写',
  'page.route.radio.regx': '正则重写',
  'page.route.form.itemLabel.from': '原路径',
};
