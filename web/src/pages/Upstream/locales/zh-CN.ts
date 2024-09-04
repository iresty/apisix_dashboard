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
  'page.upstream.step.select.upstream': '选择上游服务',
  'page.upstream.step.select.upstream.select.option': '手动填写',
  'page.upstream.step.select.upstream.select.option.serviceSelected':
    '手动填写（当前配置将会覆盖已绑定的服务）',
  'page.upstream.step.select.upstream.select.none': '无',
  'page.upstream.step.backend.server.domain.or.ip': '后端服务域名或 IP',
  'page.upstream.form.item-label.node.domain.or.ip': '目标节点',
  'page.upstream.step.input.domain.name.or.ip': '请输入域名或 IP',
  'page.upstream.step.domain.name.or.ip': '主机名或 IP',
  'page.upstream.step.host': '主机名',
  'page.upstream.step.valid.domain.name.or.ip': '请输入合法的域名或 IP',
  'page.upstream.step.input.port': '请输入',
  'page.upstream.step.port': '端口',
  'page.upstream.step.input.weight': '请输入权重',
  'page.upstream.step.weight': '权重',
  'page.upstream.step.create': '创建',
  'page.upstream.step.name': '名称',
  'page.upstream.step.name.should.unique': '名称需全局唯一',
  'page.upstream.step.input.upstream.name': '请输入上游服务的名称',
  'page.upstream.step.description': '描述',
  'page.upstream.step.input.description': '请输入上游服务的描述',
  'page.upstream.step.type': '负载均衡算法',
  'page.upstream.step.pass-host': 'Host 请求头',
  'page.upstream.step.pass-host.pass': '保持与客户端请求一致的主机名',
  'page.upstream.step.pass-host.node': '使用目标节点列表中的主机名或 IP',
  'page.upstream.step.pass-host.rewrite': '自定义 Host 请求头（即将废弃）',
  'page.upstream.step.pass-host.upstream_host': '自定义主机名',
  'page.upstream.step.connect.timeout': '连接超时',
  'page.upstream.step.connect.timeout.desc': '建立从请求到上游服务器的连接的超时时间',
  'page.upstream.step.input.connect.timeout': '请输入连接超时时间',
  'page.upstream.step.send.timeout': '发送超时',
  'page.upstream.step.send.timeout.desc': '发送数据到上游服务器的超时时间',
  'page.upstream.step.input.send.timeout': '请输入发送超时时间',
  'page.upstream.step.read.timeout': '接收超时',
  'page.upstream.step.read.timeout.desc': '从上游服务器接收数据的超时时间',
  'page.upstream.step.input.read.timeout': '请输入接收超时时间',
  'page.upstream.step.healthyCheck.healthy.check': '健康检查',
  'page.upstream.step.healthyCheck.healthy': '健康',
  'page.upstream.step.healthyCheck.unhealthy': '不健康',
  'page.upstream.step.healthyCheck.healthy.status': '健康状态',
  'page.upstream.step.healthyCheck.unhealthyStatus': '不健康状态',
  'page.upstream.step.healthyCheck.active': '主动检查',
  'page.upstream.step.healthyCheck.active.timeout': '超时时间',
  'page.upstream.step.input.healthyCheck.active.timeout': '请输入超时时间',
  'page.upstream.step.input.healthyCheck.activeInterval': '请输入间隔时间',
  'page.upstream.step.input.healthyCheck.active.req_headers': '请输入请求头',
  'page.upstream.step.healthyCheck.passive': '被动检查',
  'page.upstream.step.healthyCheck.passive.http_statuses': '状态码',
  'page.upstream.step.input.healthyCheck.passive.http_statuses': '请输入状态码',
  'page.upstream.step.healthyCheck.passive.tcp_failures': 'TCP 失败次数',
  'page.upstream.step.input.healthyCheck.passive.tcp_failures': '请输入 TCP 失败次数',
  'page.upstream.step.keepalive_pool': '连接池',
  'page.upstream.notificationMessage.enableHealthCheckFirst': '请先启用探活健康检查。',
  'page.upstream.upstream_host.required': '请输入自定义 Host 请求头',

  'page.upstream.create': '创建上游服务',
  'page.upstream.configure': '配置上游服务',
  'page.upstream.create.upstream.successfully': '创建上游服务成功',
  'page.upstream.edit.upstream.successfully': '更新上游服务成功',
  'page.upstream.create.basic.info': '基础信息',
  'page.upstream.create.preview': '预览',

  'page.upstream.list.id': 'ID',
  'page.upstream.list.name': '名称',
  'page.upstream.list.type': '类型',
  'page.upstream.list.description': '描述',
  'page.upstream.list.edit.time': '更新时间',
  'page.upstream.list.operation': '操作',
  'page.upstream.list.edit': '配置',
  'page.upstream.list.confirm.delete': '确定删除该条记录吗？',
  'page.upstream.list.confirm': '确定',
  'page.upstream.list.cancel': '取消',
  'page.upstream.list.delete.successfully': '删除记录成功',
  'page.upstream.list.delete': '删除',
  'page.upstream.list': '上游列表',
  'page.upstream.list.input': '请输入',
  'page.upstream.list.create': '创建',

  'page.upstream.type.roundrobin': '带权轮询（Round Robin）',
  'page.upstream.type.chash': '一致性哈希（CHash）',
  'page.upstream.type.ewma': '指数加权移动平均法（EWMA）',
  'page.upstream.type.least_conn': '最小连接数（least_conn）',

  'page.upstream.list.content':
    '上游列表包含了已创建的上游服务（即后端服务），可以对上游服务的多个目标节点进行负载均衡和健康检查。',

  'page.upstream.checks.active.timeout.description': '主动健康检查的套接字的超时时间',
  'page.upstream.checks.active.unhealthy.interval.description':
    '对不健康的上游服务目标节点进行主动健康检查的间隔时间，默认值为0，表示对不健康节点不进行主动健康检查。',
  'page.upstream.checks.passive.healthy.http_statuses.description':
    '当被动健康检查的探针返回值是 HTTP 状态码列表的某一个值时，代表健康状态是由代理流量产生的。',
  'page.upstream.checks.passive.unhealthy.http_statuses.description':
    '当被动健康检查的探针返回值是 HTTP 状态码列表的某一个值时，代表不健康状态是由代理流量产生的。',
  'page.upstream.checks.passive.unhealthy.http_failures.description':
    '由被动健康检查所观察，代理流量中 HTTP 失败的次数。如果达到此值，则认为上游服务目标节点是不健康的。',
  'page.upstream.checks.passive.unhealthy.tcp_failures.description':
    '被动健康检查所观察到的代理流量中 TCP 失败的次数。如果达到此值，则认为上游服务目标节点是不健康的。',
  'page.upstream.scheme': '协议',

  'page.upstream.other.configuration.invalid': '请检查上游配置',

  'page.upstream.service-name.required': '请输入自定义服务名',
  'page.upstream.step.service-name': '服务名称',
  'page.upstream.step.service-name.pass': '保持与客户端请求的主机名一致',
  'page.upstream.step.service-name.rewrite': '自定义服务名称',
  'page.upstream.step.service-name.custom_name': '自定义服务名称',
};
