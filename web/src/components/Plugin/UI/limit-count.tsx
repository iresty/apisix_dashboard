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
import React, { useState } from 'react';
import type { FormInstance } from 'antd/es/form';
import { Button, Col, Form, Input, InputNumber, Row, Select } from 'antd';
import { MinusCircleOutlined, PlusOutlined } from '@ant-design/icons';
import { useIntl } from 'umi';

type Props = {
  form: FormInstance;
  schema: Record<string, any> | undefined;
};

type PolicyProps = 'local' | 'redis' | 'redis-cluster';

type FormsProps = {
  schema: Record<string, any> | undefined;
};

const FORM_ITEM_LAYOUT = {
  labelCol: {
    span: 7,
  },
  wrapperCol: {
    span: 10,
  },
};

const FORM_ITEM_WITHOUT_LABEL = {
  wrapperCol: {
    span: 10,
    offset: 7,
  },
};

const removeBtnStyle = {
  marginLeft: 20,
  display: 'flex',
  alignItems: 'center',
};

const RedisForm: React.FC<FormsProps> = ({ schema }) => {
  const { formatMessage } = useIntl();
  const properties = schema?.properties;

  return (
    <>
      <Form.Item
        label="redis_host"
        name="redis_host"
        tooltip={formatMessage({ id: 'component.pluginForm.limit-count.redis_host.tooltip' })}
        rules={[
          {
            required: true,
            message: `${formatMessage({ id: 'component.global.pleaseEnter' })} redis_host`,
          },
          {
            min: properties.redis_host.minLength,
            message: formatMessage({
              id: 'component.pluginForm.limit-count.atLeast2Characters.rule',
            }),
          },
        ]}
        validateTrigger={['onChange', 'onBlur', 'onClick']}
      >
        <Input />
      </Form.Item>
      <Form.Item
        initialValue={properties.redis_port.default}
        label="redis_port"
        name="redis_port"
        tooltip={formatMessage({ id: 'component.pluginForm.limit-count.redis_port.tooltip' })}
      >
        <InputNumber
          min={properties.redis_port.minimum}
          max={properties.redis_port.maximum || 65535}
        />
      </Form.Item>
      <Form.Item
        label="redis_password"
        name="redis_password"
        tooltip={formatMessage({ id: 'component.pluginForm.limit-count.redis_password.tooltip' })}
        rules={[
          {
            required: true,
            message: `${formatMessage({ id: 'component.global.pleaseEnter' })} redis_password`,
          },
        ]}
        validateTrigger={['onChange', 'onBlur', 'onClick']}
      >
        <Input />
      </Form.Item>
      <Form.Item
        initialValue={properties.redis_database.default}
        label="redis_database"
        name="redis_database"
        tooltip={formatMessage({ id: 'component.pluginForm.limit-count.redis_database.tooltip' })}
      >
        <InputNumber min={properties.redis_database.minimum} />
      </Form.Item>
      <Form.Item
        initialValue={properties.redis_timeout.default}
        label="redis_timeout"
        name="redis_timeout"
        tooltip={formatMessage({ id: 'component.pluginForm.limit-count.redis_timeout.tooltip' })}
      >
        <InputNumber min={properties.redis_timeout.minimum} />
      </Form.Item>
    </>
  );
};

const RedisClusterForm: React.FC<FormsProps> = ({ schema }) => {
  const { formatMessage } = useIntl();
  const properties = schema?.properties;
  const nodesPro = properties.redis_cluster_nodes;
  const { maxLength, minLength } = nodesPro.items;
  const nodesInit = Array(nodesPro.minItems).join('.').split('.');

  return (
    <>
      <Form.Item
        label="redis_cluster_name"
        name="redis_cluster_name"
        validateTrigger={['onChange', 'onBlur', 'onClick']}
        rules={[
          {
            required: true,
            message: `${formatMessage({ id: 'component.global.pleaseEnter' })} redis_cluster_name`,
          },
        ]}
        tooltip={formatMessage({
          id: 'component.pluginForm.limit-count.redis_cluster_name.tooltip',
        })}
      >
        <Input />
      </Form.Item>
      <Form.List name="redis_cluster_nodes" initialValue={nodesInit}>
        {(fields, { add, remove }) => {
          return (
            <div>
              <Form.Item
                label="redis_cluster_nodes"
                style={{ marginBottom: 0 }}
                required
                tooltip={formatMessage({
                  id: 'component.pluginForm.limit-count.redis_cluster_nodes.tooltip',
                })}
              >
                {fields.map((field, index) => (
                  <Row style={{ marginBottom: 10 }} gutter={16} key={index}>
                    <Col>
                      <Form.Item
                        {...field}
                        noStyle
                        validateTrigger={['onChange', 'onBlur', 'onClick']}
                        rules={[
                          {
                            required: true,
                            message: `${formatMessage({
                              id: 'component.global.pleaseEnter',
                            })} redis_cluster_node`,
                          },
                          {
                            min: 2,
                            message: formatMessage({
                              id: 'component.pluginForm.limit-count.atLeast2Characters.rule',
                            }),
                          },
                        ]}
                      >
                        <Input />
                      </Form.Item>
                    </Col>
                    <Col style={{ ...removeBtnStyle, marginLeft: -10 }}>
                      {fields.length > minLength && (
                        <MinusCircleOutlined
                          className="dynamic-delete-button"
                          onClick={() => {
                            remove(field.name);
                          }}
                        />
                      )}
                    </Col>
                  </Row>
                ))}
              </Form.Item>

              <Form.Item {...FORM_ITEM_WITHOUT_LABEL}>
                {fields.length < maxLength && (
                  <Button
                    type="dashed"
                    onClick={() => {
                      add();
                    }}
                  >
                    <PlusOutlined /> {formatMessage({ id: 'component.global.add' })}
                  </Button>
                )}
              </Form.Item>
            </div>
          );
        }}
      </Form.List>
      <Form.Item
        label="redis_password"
        name="redis_password"
        tooltip={formatMessage({ id: 'component.pluginForm.limit-count.redis_password.tooltip' })}
      >
        <Input />
      </Form.Item>
      <Form.Item
        initialValue={properties.redis_timeout.default}
        label="redis_timeout"
        name="redis_timeout"
        tooltip={formatMessage({ id: 'component.pluginForm.limit-count.redis_timeout.tooltip' })}
      >
        <InputNumber min={properties.redis_timeout.minimum} />
      </Form.Item>
    </>
  );
};

const LimitCount: React.FC<Props> = ({ form, schema }) => {
  const propertires = schema?.properties;
  const [policy, setPoicy] = useState<PolicyProps>(propertires.policy.default);
  const { formatMessage } = useIntl();
  const dependSchema = schema?.dependencies.policy.oneOf;

  return (
    <Form form={form} {...FORM_ITEM_LAYOUT}>
      <Form.Item
        label="count"
        name="count"
        rules={[
          {
            required: true,
            message: `${formatMessage({ id: 'component.global.pleaseEnter' })} count`,
          },
        ]}
        tooltip={formatMessage({ id: 'component.pluginForm.limit-count.count.tooltip' })}
        validateTrigger={['onChange', 'onBlur', 'onClick']}
      >
        <InputNumber min={propertires.count.exclusiveMinimum} />
      </Form.Item>
      <Form.Item
        label="time_window"
        name="time_window"
        rules={[
          {
            required: true,
            message: `${formatMessage({ id: 'component.global.pleaseEnter' })} time_window`,
          },
        ]}
        tooltip={formatMessage({ id: 'component.pluginForm.limit-count.time_window.tooltip' })}
        validateTrigger={['onChange', 'onBlur', 'onClick']}
      >
        <InputNumber min={propertires.time_window.exclusiveMinimum} />
      </Form.Item>
      <Form.Item
        initialValue={propertires.key.default}
        label="key"
        name="key"
        tooltip={formatMessage({ id: 'component.pluginForm.limit-count.key.tooltip' })}
      >
        <Select>
          {propertires.key.enum.map((item: string) => (
            <Select.Option value={item} key={item}>
              {item}
            </Select.Option>
          ))}
        </Select>
      </Form.Item>
      <Form.Item
        initialValue={propertires.rejected_code.default}
        label="rejected_code"
        name="rejected_code"
        tooltip={formatMessage({ id: 'component.pluginForm.limit-count.rejected_code.tooltip' })}
      >
        <InputNumber
          min={propertires.rejected_code.minimum}
          max={propertires.rejected_code.maximum}
        />
      </Form.Item>
      <Form.Item
        label="rejected_msg"
        name="rejected_msg"
        tooltip={formatMessage({ id: 'component.pluginForm.limit-count.rejected_msg.tooltip' })}
      >
        <Input />
      </Form.Item>
      <Form.Item
        initialValue={policy}
        label="policy"
        name="policy"
        tooltip={formatMessage({ id: 'component.pluginForm.limit-count.policy.tooltip' })}
      >
        <Select
          onChange={(e: PolicyProps) => {
            setPoicy(e);
          }}
        >
          {propertires.policy.enum.map((item: string) => (
            <Select.Option value={item} key={item}>
              {item}
            </Select.Option>
          ))}
        </Select>
      </Form.Item>
      <Form.Item
        shouldUpdate={(prev, next) => prev.policy !== next.policy}
        style={{ display: 'none' }}
      >
        {() => {
          setPoicy(form.getFieldValue('policy'));
        }}
      </Form.Item>
      {Boolean(policy === 'redis') && <RedisForm schema={dependSchema[1]} />}
      {Boolean(policy === 'redis-cluster') && <RedisClusterForm schema={dependSchema[2]} />}
    </Form>
  );
};

export default LimitCount;
