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

export const FORM_ITEM_LAYOUT = {
  labelCol: {
    span: 4,
  },
  wrapperCol: {
    span: 8,
  },
};

const LimitReq: React.FC<Props> = ({ form, schema }) => {
  const { formatMessage } = useIntl();
  const propertires = schema?.properties;
  return (
    <Form form={form} {...FORM_ITEM_LAYOUT}>
      <Form.Item
        label="rate"
        name="rate"
        rules={[
          {
            required: true,
            message: `${formatMessage({ id: 'component.global.pleaseEnter' })} rate`,
          },
        ]}
        tooltip={formatMessage({ id: 'component.pluginForm.limit-req.rate.tooltip' })}
        validateTrigger={['onChange', 'onBlur', 'onClick']}
      >
        <InputNumber min={propertires.rate.exclusiveMinimum} required />
      </Form.Item>
      <Form.Item
        label="burst"
        name="burst"
        rules={[
          {
            required: true,
            message: `${formatMessage({ id: 'component.global.pleaseEnter' })} burst`,
          },
        ]}
        tooltip={formatMessage({ id: 'component.pluginForm.limit-req.burst.tooltip' })}
        validateTrigger={['onChange', 'onBlur', 'onClick']}
      >
        <InputNumber min={propertires.burst.minimum} required />
      </Form.Item>
      <Form.Item
        label="key"
        name="key"
        rules={[
          {
            required: true,
            message: `${formatMessage({ id: 'component.global.pleaseChoose' })} key`,
          },
        ]}
        tooltip={formatMessage({ id: 'component.pluginForm.limit-req.key.tooltip' })}
        validateTrigger={['onChange', 'onBlur', 'onClick']}
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
        tooltip={formatMessage({ id: 'component.pluginForm.limit-req.rejected_code.tooltip' })}
      >
        <InputNumber
          min={propertires.rejected_code.minimum}
          max={propertires.rejected_code.maximum}
        />
      </Form.Item>
      <Form.Item
        label="rejected_msg"
        name="rejected_msg"
        tooltip={formatMessage({ id: 'component.pluginForm.limit-req.rejected_msg.tooltip' })}
      >
        <Input />
      </Form.Item>
      <Form.Item
        label="nodelay"
        name="nodelay"
        valuePropName="checked"
        tooltip={formatMessage({ id: 'component.pluginForm.limit-req.nodelay.tooltip' })}
      >
        <Switch defaultChecked={false} />
      </Form.Item>
    </Form>
  );
};

export default LimitReq;
