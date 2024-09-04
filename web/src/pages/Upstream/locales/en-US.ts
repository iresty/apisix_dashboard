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
  'page.upstream.step.select.upstream': 'Upstream',
  'page.upstream.step.select.upstream.select.option': 'Custom',
  'page.upstream.step.select.upstream.select.option.serviceSelected':
    'Custom (The current configuration will override the bound service)',
  'page.upstream.step.select.upstream.select.none': 'None',
  'page.upstream.step.backend.server.domain.or.ip': 'Backend Server Host/IP',
  'page.upstream.form.item-label.node.domain.or.ip': 'Targets',
  'page.upstream.step.input.domain.name.or.ip': 'Please enter domain or IP',
  'page.upstream.step.valid.domain.name.or.ip': 'Please enter valid a domain or IP',
  'page.upstream.step.domain.name.or.ip': 'Hostname or IP',
  'page.upstream.step.host': 'Host',
  'page.upstream.step.input.port': 'Please enter port number',
  'page.upstream.step.port': 'Port',
  'page.upstream.step.input.weight': 'Please enter weight',
  'page.upstream.step.weight': 'Weight',
  'page.upstream.step.create': 'Create',
  'page.upstream.step.name': 'Name',
  'page.upstream.step.name.should.unique': 'Name should be unique',
  'page.upstream.step.input.upstream.name': 'Please enter upstream name',
  'page.upstream.step.description': 'Description',
  'page.upstream.step.input.description': "Please enter upstream's description",
  'page.upstream.step.type': 'Algorithm',
  'page.upstream.step.pass-host': 'Hostname',
  'page.upstream.step.pass-host.pass': 'Keep the same Host from client request',
  'page.upstream.step.pass-host.node': 'Use the domain or IP from Node List',
  'page.upstream.step.pass-host.rewrite': 'Custom Host (Will be deprecated in the future)',
  'page.upstream.step.pass-host.upstream_host': 'Custom Host',
  'page.upstream.step.connect.timeout': 'Connect Timeout',
  'page.upstream.step.connect.timeout.desc':
    'Timeout for establishing a connection from the request to the upstream server',
  'page.upstream.step.input.connect.timeout': 'Please enter connect timeout',
  'page.upstream.step.send.timeout': 'Send Timeout',
  'page.upstream.step.send.timeout.desc': 'Timeout for sending data to upstream servers',
  'page.upstream.step.input.send.timeout': 'Please enter send timeout',
  'page.upstream.step.read.timeout': 'Read Timeout',
  'page.upstream.step.read.timeout.desc': 'Timeout for receiving data from upstream servers',
  'page.upstream.step.input.read.timeout': 'Please enter read timeout',
  'page.upstream.step.healthyCheck.healthy.check': 'Health Check',
  'page.upstream.step.healthyCheck.healthy': 'Healthy',
  'page.upstream.step.healthyCheck.unhealthy': 'Unhealthy',
  'page.upstream.step.healthyCheck.healthy.status': 'Healthy Status',
  'page.upstream.step.healthyCheck.unhealthyStatus': 'Unhealthy Status',
  'page.upstream.step.healthyCheck.active': 'Active',
  'page.upstream.step.healthyCheck.active.timeout': 'Timeout',
  'page.upstream.step.input.healthyCheck.active.timeout': 'Please enter timeout',
  'page.upstream.step.input.healthyCheck.activeInterval': 'Please enter interval',
  'page.upstream.step.input.healthyCheck.active.req_headers': 'Please enter request headers',
  'page.upstream.step.healthyCheck.passive': 'Passive',
  'page.upstream.step.healthyCheck.passive.http_statuses': 'HTTP Status',
  'page.upstream.step.input.healthyCheck.passive.http_statuses': 'Please enter http status',
  'page.upstream.step.healthyCheck.passive.tcp_failures': 'TCP Failures',
  'page.upstream.step.input.healthyCheck.passive.tcp_failures': 'Please enter TCP failures',
  'page.upstream.step.keepalive_pool': 'Keepalive Pool',
  'page.upstream.notificationMessage.enableHealthCheckFirst': 'Please enable health check first.',
  'page.upstream.upstream_host.required': 'Please enter the custom Host',

  'page.upstream.create': 'Create Upstream',
  'page.upstream.configure': 'Configure Upstream',
  'page.upstream.create.upstream.successfully': 'Create Upstream Successfully',
  'page.upstream.edit.upstream.successfully': 'Configure Upstream Successfully',
  'page.upstream.create.basic.info': 'Basic Information',
  'page.upstream.create.preview': 'Preview',

  'page.upstream.list.id': 'ID',
  'page.upstream.list.name': 'Name',
  'page.upstream.list.type': 'Type',
  'page.upstream.list.description': 'Description',
  'page.upstream.list.edit.time': 'Update Time',
  'page.upstream.list.operation': 'Operation',
  'page.upstream.list.edit': 'Configure',
  'page.upstream.list.confirm.delete': 'Are you sure to delete ?',
  'page.upstream.list.confirm': 'Confirm',
  'page.upstream.list.cancel': 'Cancel',
  'page.upstream.list.delete.successfully': 'Delete Upstream Successfully',
  'page.upstream.list.delete': 'Delete',
  'page.upstream.list': 'Upstream List',
  'page.upstream.list.input': 'Please enter',
  'page.upstream.list.create': 'Create',

  'page.upstream.type.roundrobin': 'Round Robin',
  'page.upstream.type.chash': 'CHash',
  'page.upstream.type.ewma': 'EWMA',
  'page.upstream.type.least_conn': 'Least conn',

  'page.upstream.list.content':
    'The upstream list contains the created upstream services (i.e., backend services) and allows load balancing and health checking of multiple target nodes of the upstream services.',

  'page.upstream.checks.active.timeout.description':
    'Socket timeout for active checks (in seconds)',
  'page.upstream.checks.active.unhealthy.interval.description':
    'Interval between checks for unhealthy targets (in seconds)',
  'page.upstream.checks.passive.healthy.http_statuses.description':
    'Which HTTP statuses to consider a success',
  'page.upstream.checks.passive.unhealthy.http_statuses.description':
    'Which HTTP statuses to consider a failure',
  'page.upstream.checks.passive.unhealthy.http_failures.description':
    'Number of HTTP failures to consider a target unhealthy',
  'page.upstream.checks.passive.unhealthy.tcp_failures.description':
    'Number of TCP failures to consider a target unhealthy',
  'page.upstream.scheme': 'Scheme',

  'page.upstream.other.configuration.invalid': 'Please check the Upstream configuration',

  'page.upstream.service-name.required': 'Please enter the custom server ame',
  'page.upstream.step.service-name': 'Servername',
  'page.upstream.step.service-name.pass': 'Keep the same as Host from client request',
  'page.upstream.step.service-name.rewrite': 'Custom Server Name',
  'page.upstream.step.service-name.custom_name': 'Custom Server Name',

};
