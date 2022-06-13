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
import Form from 'antd/es/form';
import { Button, Input, Radio, Row, Col, Select } from 'antd';
import { PlusOutlined, MinusCircleOutlined } from '@ant-design/icons';
import { useIntl } from 'umi';

import PanelSection from '@/components/PanelSection';

import {
  FORM_ITEM_WITHOUT_LABEL,
  SCHEME_REWRITE,
  URI_REWRITE_TYPE,
  HOST_REWRITE_TYPE,
} from '@/pages/Route/constants';

const removeBtnStyle = {
  marginLeft: 20,
  display: 'flex',
  alignItems: 'center',
};

/**
 * https://apisix.apache.org/docs/apisix/plugins/proxy-rewrite
 * UI for ProxyRewrite plugin
 */
const ProxyRewrite: React.FC<RouteModule.Step1PassProps> = ({ form, disabled }) => {
  const { formatMessage } = useIntl();

  const getUriRewriteItems = () => {
    switch (form.getFieldValue('URIRewriteType')) {
      case URI_REWRITE_TYPE.STATIC:
        return (
          <Form.Item
            label={formatMessage({ id: 'page.route.form.itemLabel.newPath' })}
            name={['proxyRewrite', 'uri']}
            rules={[
              {
                required: true,
                message: `${formatMessage({ id: 'component.global.pleaseEnter' })} ${formatMessage({
                  id: 'page.route.form.itemLabel.newPath',
                })}`,
              },
            ]}
          >
            <Input
              placeholder={`${formatMessage({
                id: 'component.global.pleaseEnter',
              })} ${formatMessage({ id: 'page.route.form.itemLabel.newPath' })}`}
              disabled={disabled}
            />
          </Form.Item>
        );
      case URI_REWRITE_TYPE.REGEXP:
        return (
          <Form.List name={['proxyRewrite', 'regex_uri']} initialValue={['', '']}>
            {(fields) =>
              fields.map((field, index) => {
                switch (index) {
                  case 0:
                    return (
                      <Form.Item
                        label={formatMessage({ id: 'page.route.form.itemLabel.regex' })}
                        name={field.name}
                        key={field.name}
                        rules={[
                          {
                            required: true,
                            message: `${formatMessage({
                              id: 'component.global.pleaseEnter',
                            })} ${formatMessage({ id: 'page.route.form.itemLabel.regex' })}`,
                          },
                        ]}
                      >
                        <Input
                          placeholder={`${formatMessage({
                            id: 'component.global.pleaseEnter',
                          })} ${formatMessage({ id: 'page.route.form.itemLabel.regex' })}`}
                          disabled={disabled}
                        />
                      </Form.Item>
                    );
                  case 1:
                    return (
                      <Form.Item
                        label={formatMessage({ id: 'page.route.form.itemLabel.template' })}
                        name={field.name}
                        key={field.name}
                      >
                        <Input
                          placeholder={`${formatMessage({
                            id: 'component.global.pleaseEnter',
                          })} ${formatMessage({ id: 'page.route.form.itemLabel.template' })}`}
                          disabled={disabled}
                        />
                      </Form.Item>
                    );
                  default:
                    return null;
                }
              })
            }
          </Form.List>
        );
      case URI_REWRITE_TYPE.KEEP:
      default:
        return null;
    }
  };

  const getHostRewriteItems = () => {
    switch (form.getFieldValue('hostRewriteType')) {
      case HOST_REWRITE_TYPE.REWRITE:
        return (
          <Form.Item
            label={formatMessage({ id: 'page.route.form.itemLabel.newHost' })}
            name={['proxyRewrite', 'host']}
            rules={[
              {
                required: true,
                message: `${formatMessage({ id: 'component.global.pleaseEnter' })} ${formatMessage({
                  id: 'page.route.form.itemLabel.newHost',
                })}`,
              },
            ]}
          >
            <Input
              placeholder={`${formatMessage({
                id: 'component.global.pleaseEnter',
              })} ${formatMessage({ id: 'page.route.form.itemLabel.newHost' })}`}
              disabled={disabled}
            />
          </Form.Item>
        );
      case HOST_REWRITE_TYPE.KEEP:
      default:
        return null;
    }
  };

  const SchemeComponent: React.FC = () => {
    const options = [
      {
        value: SCHEME_REWRITE.KEEP,
        label: formatMessage({ id: 'page.route.radio.staySame' }),
      },
      {
        value: SCHEME_REWRITE.HTTP,
        label: SCHEME_REWRITE.HTTP.toLocaleUpperCase(),
      },
      {
        value: SCHEME_REWRITE.HTTPS,
        label: SCHEME_REWRITE.HTTPS.toLocaleUpperCase(),
      },
    ];

    return (
      <Form.Item
        label={formatMessage({ id: 'page.route.form.itemLabel.scheme' })}
        name={['proxyRewrite', 'scheme']}
      >
        <Radio.Group disabled={disabled}>
          {options.map((item) => (
            <Radio value={item.value} key={item.value}>
              {item.label}
            </Radio>
          ))}
        </Radio.Group>
      </Form.Item>
    );
  };

  const URIRewriteType: React.FC = () => {
    const options = [
      {
        value: URI_REWRITE_TYPE.KEEP,
        label: formatMessage({ id: 'page.route.radio.staySame' }),
      },
      {
        value: URI_REWRITE_TYPE.STATIC,
        label: formatMessage({ id: 'page.route.radio.static' }),
        dataCypress: 'uri-static',
      },
      {
        value: URI_REWRITE_TYPE.REGEXP,
        label: formatMessage({ id: 'page.route.radio.regex' }),
      },
    ];

    return (
      <React.Fragment>
        <Form.Item
          label={formatMessage({ id: 'page.route.form.itemLabel.URIRewriteType' })}
          name="URIRewriteType"
        >
          <Radio.Group disabled={disabled}>
            {options.map((item) => (
              <Radio data-cy={item.dataCypress} value={item.value} key={item.value}>
                {item.label}
              </Radio>
            ))}
          </Radio.Group>
        </Form.Item>
        <Form.Item
          shouldUpdate={(prevValues, curValues) =>
            prevValues.URIRewriteType !== curValues.URIRewriteType
          }
          noStyle
        >
          {() => {
            return getUriRewriteItems();
          }}
        </Form.Item>
      </React.Fragment>
    );
  };

  const HostRewriteType: React.FC = () => {
    const options = [
      {
        label: formatMessage({ id: 'page.route.radio.staySame' }),
        value: HOST_REWRITE_TYPE.KEEP,
        dataCypress: 'host-keep',
      },
      {
        label: formatMessage({ id: 'page.route.radio.static' }),
        value: HOST_REWRITE_TYPE.REWRITE,
        dataCypress: 'host-static',
      },
    ];
    return (
      <React.Fragment>
        <Form.Item
          label={formatMessage({ id: 'page.route.form.itemLabel.hostRewriteType' })}
          name="hostRewriteType"
        >
          <Radio.Group disabled={disabled}>
            {options.map((item) => (
              <Radio data-cy={item.dataCypress} value={item.value} key={item.value}>
                {item.label}
              </Radio>
            ))}
          </Radio.Group>
        </Form.Item>
        <Form.Item
          shouldUpdate={(prevValues, curValues) =>
            prevValues.hostRewriteType !== curValues.hostRewriteType
          }
          noStyle
        >
          {() => {
            return getHostRewriteItems();
          }}
        </Form.Item>
      </React.Fragment>
    );
  };

  const MethodRewriteType: React.FC = () => {
    const methods = [
      'GET',
      'POST',
      'PUT',
      'DELETE',
      'PATCH',
      'OPTIONS',
      'HEAD',
      'TRACE',
      'MKCOL',
      'COPY',
      'MOVE',
      'PROPFIND',
      'LOCK',
      'UNLOCK',
    ];
    return (
      <React.Fragment>
        <Form.Item
          label={formatMessage({ id: 'page.route.form.itemLabel.methodRewrite' })}
          name={['proxyRewrite', 'method']}
          wrapperCol={{ span: 3 }}
          initialValue=""
        >
          <Select data-cy="proxyRewrite-method">
            <Select.Option value="">
              {formatMessage({ id: 'page.route.select.option.methodRewriteNone' })}
            </Select.Option>
            {methods.map((method) => (
              <Select.Option value={method} key={method}>
                {method}
              </Select.Option>
            ))}
          </Select>
        </Form.Item>
      </React.Fragment>
    );
  };

  const Headers: React.FC = () => {
    return (
      <Form.List name={['proxyRewrite', 'kvHeaders']} initialValue={[{}]}>
        {(fields, { add, remove }) => (
          <>
            <Form.Item
              label={formatMessage({ id: 'page.route.form.itemLabel.headerRewrite' })}
              style={{ marginBottom: 0 }}
            >
              {fields.map((field, index) => (
                <Row gutter={12} key={index} style={{ marginBottom: 10 }}>
                  <Col span={5}>
                    <Form.Item
                      name={[field.name, 'key']}
                      fieldKey={[field.fieldKey, 'key']}
                      noStyle
                    >
                      <Input
                        placeholder={`${formatMessage({
                          id: 'component.global.pleaseEnter',
                        })} ${formatMessage({ id: 'page.route.parameterName' })}`}
                        disabled={disabled}
                      />
                    </Form.Item>
                  </Col>
                  <Col span={5}>
                    <Form.Item
                      name={[field.name, 'value']}
                      fieldKey={[field.fieldKey, 'value']}
                      noStyle
                    >
                      <Input
                        placeholder={`${formatMessage({
                          id: 'component.global.pleaseEnter',
                        })} ${formatMessage({ id: 'page.route.value' })}`}
                        disabled={disabled}
                      />
                    </Form.Item>
                  </Col>
                  {!disabled && fields.length > 1 && (
                    <Col style={{ ...removeBtnStyle, marginLeft: -5 }}>
                      <MinusCircleOutlined
                        className="dynamic-delete-button"
                        onClick={() => remove(field.name)}
                      />
                    </Col>
                  )}
                </Row>
              ))}
            </Form.Item>
            <Form.Item {...FORM_ITEM_WITHOUT_LABEL}>
              <Button
                data-cy="create-new-rewrite-header"
                type="dashed"
                disabled={disabled}
                onClick={() => add()}
                icon={<PlusOutlined />}
              >
                {formatMessage({ id: 'component.global.add' })}
              </Button>
            </Form.Item>
          </>
        )}
      </Form.List>
    );
  };

  return (
    <PanelSection title={formatMessage({ id: 'page.route.panelSection.title.requestOverride' })}>
      <SchemeComponent />
      <URIRewriteType />
      <HostRewriteType />
      <MethodRewriteType />
      <Headers />
    </PanelSection>
  );
};

export default ProxyRewrite;
