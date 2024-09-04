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
import type { FormInstance } from 'antd';
import { Form, Input, Select } from 'antd';
import React from 'react';
import { useIntl } from 'umi';

import ServiceDiscoveryArgs from '@/components/Upstream/components/ServiceDiscoveryArgs';
import ServerName from '@/components/Upstream/components/ServerName';

type Props = {
  form: FormInstance;
  readonly?: boolean;
};

const discoveryType = {
  dns: {},
  consul: {},
  consul_kv: {},
  nacos: {
    args: ['group_name', 'namespace_id'],
  },
  eureka: {},
  kubernetes: {},
};

const ServiceDiscovery: React.FC<Props> = ({ readonly, form }) => {
  const { formatMessage } = useIntl();

  return (
    <React.Fragment>
      <Form.Item
        name="discovery_type"
        label={formatMessage({ id: 'component.upstream.fields.discovery_type' })}
        tooltip={formatMessage({ id: 'component.upstream.fields.discovery_type.tooltip' })}
        rules={[{ required: true }]}
      >
        <Select
          disabled={readonly}
          placeholder={formatMessage({
            id: 'component.upstream.fields.discovery_type.placeholder',
          })}
        >
          {Object.keys(discoveryType).map((item) => {
            return (
              <Select.Option key={item} value={item}>
                {formatMessage({ id: `component.upstream.fields.discovery_type.type.${item}` })}
              </Select.Option>
            );
          })}
        </Select>
      </Form.Item>
      <ServerName form={form} readonly={readonly} />
      <Form.Item shouldUpdate noStyle>
        {() => {
          if (!form.getFieldValue('discovery_type')) return null;

          const { args } = discoveryType[form.getFieldValue('discovery_type')];
          if (args && args.length > 0) {
            return <ServiceDiscoveryArgs readonly={readonly} args={args} />;
          }
          return null;
        }}
      </Form.Item>
    </React.Fragment>
  );
};

export default ServiceDiscovery;
