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
import { Button, Table, Modal, Form, Select, Input, Space, notification } from 'antd';
import { useIntl } from 'umi';

import PanelSection from '@/components/PanelSection';

const MatchingRulesView: React.FC<RouteModule.Step1PassProps> = ({
  advancedMatchingRules,
  disabled,
  onChange = () => {},
}) => {
  const [visible, setVisible] = useState(false);
  const [mode, setMode] = useState<RouteModule.ModalType>('CREATE');
  const [namePlaceholder, setNamePlaceholder] = useState('');
  const [modalForm] = Form.useForm();

  const [operatorValueSample, setOperatorValueSample] = useState('');

  const { Option } = Select;

  const { formatMessage } = useIntl();

  const onOk = () => {
    modalForm.validateFields().then((value: RouteModule.MatchingRule) => {
      if (value.operator === 'IN') {
        try {
          JSON.parse(value.value as string);
        } catch (error) {
          notification.warning({
            message: formatMessage({ id: 'page.route.fields.vars.invalid' }),
            description: formatMessage({ id: 'page.route.fields.vars.in.invalid' }),
          });
          return;
        }
      }
      if (mode === 'EDIT') {
        const key = modalForm.getFieldValue('key');
        onChange({
          action: 'advancedMatchingRulesChange',
          data: advancedMatchingRules.map((rule) => {
            if (rule.key === key) {
              return { ...value, key };
            }
            return rule;
          }),
        });
      } else {
        const rule = {
          ...value,
          key: Math.random().toString(36).slice(2),
        };
        onChange({
          action: 'advancedMatchingRulesChange',
          data: advancedMatchingRules.concat(rule),
        });
      }
      modalForm.resetFields();
      setVisible(false);
    });
  };

  const handleEdit = (record: RouteModule.MatchingRule) => {
    setMode('EDIT');
    setVisible(true);
    modalForm.setFieldsValue(record);
  };

  const handleRemove = (key: string) => {
    onChange({
      action: 'advancedMatchingRulesChange',
      data: advancedMatchingRules.filter((item) => item.key !== key),
    });
  };

  const OperatorRenderText = {
    '==': formatMessage({ id: 'page.route.equal' }),
    '~=': formatMessage({ id: 'page.route.unequal' }),
    '>': formatMessage({ id: 'page.route.greaterThan' }),
    '<': formatMessage({ id: 'page.route.lessThan' }),
    '~~': formatMessage({ id: 'page.route.regexMatch' }),
    '~*': formatMessage({ id: 'page.route.caseInsensitiveRegexMatch' }),
    IN: formatMessage({ id: 'page.route.in' }),
    HAS: formatMessage({ id: 'page.route.has' }),
    '!': formatMessage({ id: 'page.route.reverse' }),
  };

  const columns = [
    {
      title: formatMessage({ id: 'page.route.parameterPosition' }),
      key: 'position',
      render: (text: RouteModule.MatchingRule) => {
        let renderText;
        switch (text.position) {
          case 'http':
            renderText = formatMessage({ id: 'page.route.httpRequestHeader' });
            break;
          case 'arg':
            renderText = formatMessage({ id: 'page.route.requestParameter' });
            break;
          case 'cookie':
            renderText = 'Cookie';
            break;
          case 'buildin':
            renderText = formatMessage({ id: 'page.route.buildinParameter' });
            break;
          default:
            renderText = '';
        }
        return renderText;
      },
    },
    {
      title: formatMessage({ id: 'page.route.parameterName' }),
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: formatMessage({ id: 'page.route.operationalCharacter' }),
      key: 'operator',
      render: (text: RouteModule.MatchingRule) => OperatorRenderText[text.operator],
    },
    {
      title: formatMessage({ id: 'page.route.value' }),
      dataIndex: 'value',
      key: 'value',
    },
    disabled
      ? {}
      : {
          title: formatMessage({ id: 'component.global.operation' }),
          key: 'action',
          render: (_: any, record: RouteModule.MatchingRule) => (
            <Space size="middle">
              <a onClick={() => handleEdit(record)}>
                {formatMessage({ id: 'component.global.edit' })}
              </a>
              <a onClick={() => handleRemove(record.key)}>
                {formatMessage({ id: 'component.global.delete' })}
              </a>
            </Space>
          ),
        },
  ].filter((item) => Object.keys(item).length);

  const renderModal = () => {
    return (
      <Modal
        title={
          mode === 'EDIT'
            ? formatMessage({ id: 'page.route.rule.edit' })
            : formatMessage({ id: 'page.route.rule.create' })
        }
        centered
        visible
        onOk={onOk}
        onCancel={() => {
          setVisible(false);
          modalForm.resetFields();
        }}
        okText={formatMessage({ id: 'component.global.confirm' })}
        cancelText={formatMessage({ id: 'component.global.cancel' })}
        maskClosable={false}
        destroyOnClose
      >
        <Form form={modalForm} layout="vertical">
          <Form.Item
            label={formatMessage({ id: 'page.route.parameterPosition' })}
            name="position"
            rules={[
              {
                required: true,
                message: `${formatMessage({ id: 'component.global.pleaseEnter' })} ${formatMessage({
                  id: 'page.route.parameterPosition',
                })}`,
              },
            ]}
          >
            <Select
              onChange={(value) => {
                if (value === 'http') {
                  setNamePlaceholder(
                    formatMessage({ id: 'page.route.input.placeholder.parameterNameHttpHeader' }),
                  );
                } else {
                  setNamePlaceholder(
                    formatMessage({
                      id: 'page.route.input.placeholder.parameterNameRequestParameter',
                    }),
                  );
                }
              }}
            >
              <Option value="http">{formatMessage({ id: 'page.route.httpRequestHeader' })}</Option>
              <Option value="arg">{formatMessage({ id: 'page.route.requestParameter' })}</Option>
              <Option value="cookie">Cookie</Option>
              <Option value="buildin">
                {formatMessage({ id: 'page.route.buildinParameter' })}
              </Option>
            </Select>
          </Form.Item>
          <Form.Item
            label={formatMessage({ id: 'page.route.parameterName' })}
            name="name"
            rules={[
              {
                required: true,
                message: `${formatMessage({ id: 'component.global.pleaseEnter' })} ${formatMessage({
                  id: 'page.route.parameterName',
                })}`,
              },
              {
                pattern: new RegExp(/^([a-zA-Z][a-zA-Z0-9_-]*$)/, 'g'),
                message: formatMessage({ id: 'component.global.input.ruleMessage.name' }),
              },
            ]}
            tooltip={formatMessage({
              id: 'page.route.form.itemRulesRequiredMessage.parameterName',
            })}
            extra={namePlaceholder}
          >
            <Input />
          </Form.Item>
          <Form.Item
            label={formatMessage({ id: 'page.route.operationalCharacter' })}
            name="operator"
            rules={[
              {
                required: true,
                message: `${formatMessage({ id: 'component.global.pleaseChoose' })} ${formatMessage(
                  {
                    id: 'page.route.operationalCharacter',
                  },
                )}`,
              },
            ]}
          >
            <Select
              onChange={(e: string) => {
                switch (e) {
                  case 'IN':
                    setOperatorValueSample(
                      formatMessage({ id: 'page.route.advanced-match.operator.sample.IN' }),
                    );
                    break;
                  case '~~':
                  case '~*':
                    setOperatorValueSample(
                      formatMessage({ id: 'page.route.advanced-match.operator.sample.~~' }),
                    );
                    break;
                  default:
                    setOperatorValueSample('');
                }
              }}
            >
              {Object.keys(OperatorRenderText).map((item) => (
                <Option value={item} key={item}>
                  {OperatorRenderText[item]}
                </Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item
            label={formatMessage({ id: 'page.route.value' })}
            name="value"
            rules={[
              {
                required: true,
                message: `${formatMessage({ id: 'component.global.pleaseEnter' })} ${formatMessage({
                  id: 'page.route.value',
                })}`,
              },
            ]}
            extra={operatorValueSample}
          >
            <Input />
          </Form.Item>
        </Form>
      </Modal>
    );
  };

  return (
    <PanelSection
      title={formatMessage({ id: 'page.route.panelSection.title.advancedMatchRule' })}
      tooltip={formatMessage({ id: 'page.route.advanced-match.tooltip' })}
    >
      {!disabled && (
        <Button
          onClick={() => {
            setMode('CREATE');
            setVisible(true);
          }}
          type="primary"
          style={{
            marginBottom: 16,
          }}
        >
          {formatMessage({ id: 'component.global.add' })}
        </Button>
      )}
      <Table key="table" bordered dataSource={advancedMatchingRules} columns={columns} />
      {/* NOTE: tricky way, switch visible on Modal component will occur error */}
      {visible ? renderModal() : null}
    </PanelSection>
  );
};

export default MatchingRulesView;
