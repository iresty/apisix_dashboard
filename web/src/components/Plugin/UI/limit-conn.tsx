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
import React from 'react';
import type { FormInstance } from 'antd/es/form';
import { Form, Input, InputNumber, Select, Switch } from 'antd';
import { useIntl } from 'umi';

type Props = {
  form: FormInstance;
  schema: Record<string, any> | undefined;
  ref?: any;
};

const FORM_ITEM_LAYOUT = {
  labelCol: {
    span: 6,
  },
  wrapperCol: {
    span: 8,
  },
};

const LimitConn: React.FC<Props> = ({ form, schema }) => {
  const { formatMessage } = useIntl();
  const propertires = schema?.properties;
  const onlyUseDefaultDelay = form.getFieldValue('only_use_default_delay')
    ? form.getFieldValue('only_use_default_delay')
    : false;
  return (
    <Form form={form} {...FORM_ITEM_LAYOUT}>
      <Form.Item
        label="conn"
        required
        name="conn"
        tooltip={formatMessage({ id: 'component.pluginForm.limit-conn.conn.tooltip' })}
      >
        <InputNumber min={propertires.conn.exclusiveMinimum} required />
      </Form.Item>
      <Form.Item
        label="burst"
        required
        name="burst"
        tooltip={formatMessage({ id: 'component.pluginForm.limit-conn.burst.tooltip' })}
      >
        <InputNumber min={propertires.burst.minimum} required />
      </Form.Item>
      <Form.Item
        label="default_conn_delay"
        required
        name="default_conn_delay"
        tooltip={formatMessage({
          id: 'component.pluginForm.limit-conn.default_conn_delay.tooltip',
        })}
      >
        <InputNumber step={0.001} min={propertires.default_conn_delay.exclusiveMinimum} required />
      </Form.Item>

      <Form.Item
        label="only_use_default_delay"
        name="only_use_default_delay"
        initialValue={propertires.only_use_default_delay.default}
        tooltip={formatMessage({
          id: 'component.pluginForm.limit-conn.only_use_default_delay.tooltip',
        })}
      >
        <Switch defaultChecked={onlyUseDefaultDelay} />
      </Form.Item>

      <Form.Item
        label="key"
        required
        name="key"
        tooltip={formatMessage({ id: 'component.pluginForm.limit-conn.key.tooltip' })}
      >
        <Select>
          {propertires.key.enum.map((item: string) => {
            return (
              <Select.Option value={item} key={item}>
                {item}
              </Select.Option>
            );
          })}
        </Select>
      </Form.Item>
      <Form.Item
        label="rejected_code"
        name="rejected_code"
        initialValue={propertires.rejected_code.default}
        tooltip={formatMessage({ id: 'component.pluginForm.limit-conn.rejected_code.tooltip' })}
      >
        <InputNumber
          min={propertires.rejected_code.minimum}
          max={propertires.rejected_code.maximum}
          required
        />
      </Form.Item>
      <Form.Item
        label="rejected_msg"
        name="rejected_msg"
        tooltip={formatMessage({ id: 'component.pluginForm.limit-conn.rejected_msg.tooltip' })}
      >
        <Input />
      </Form.Item>
    </Form>
  );
};

export default LimitConn;
