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
import { Form, Input, Switch } from 'antd';
import { useIntl } from 'umi';

type Props = {
  form: FormInstance;
};

const FORM_ITEM_LAYOUT = {
  labelCol: {
    span: 7,
  },
  wrapperCol: {
    span: 18,
  },
};

const RequestId: React.FC<Props> = ({ form }) => {
  const { formatMessage } = useIntl();
  return (
    <Form form={form} {...FORM_ITEM_LAYOUT}>
      <Form.Item
        label="header_name"
        name="header_name"
        initialValue="X-Request-Id"
        rules={[
          {
            message: `${formatMessage({ id: 'component.global.pleaseEnter' })} header_name`,
          },
        ]}
        tooltip={formatMessage({ id: 'component.pluginForm.request-id.header_name.tooltip' })}
        validateTrigger={['onChange', 'onBlur', 'onClick']}
      >
        <Input />
      </Form.Item>
      <Form.Item
        label="include_in_response"
        name="include_in_response"
        tooltip={formatMessage({
          id: 'component.pluginForm.request-id.include_in_response.tooltip',
        })}
        valuePropName="checked"
      >
        <Switch />
      </Form.Item>
    </Form>
  );
};

export default RequestId;
