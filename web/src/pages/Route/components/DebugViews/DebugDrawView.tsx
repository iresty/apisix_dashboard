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
import React, { useEffect, useState } from 'react';
import {
  Button,
  Card,
  Divider,
  Drawer,
  Form,
  Input,
  notification,
  Radio,
  Select,
  Spin,
  Tabs,
} from 'antd';
import { useIntl } from 'umi';
import queryString from 'query-string';
import Base64 from 'base-64';
import CopyToClipboard from 'react-copy-to-clipboard';
import { CopyOutlined } from '@ant-design/icons';

import PanelSection from '@/components/PanelSection';

import {
  DEBUG_BODY_MODE_SUPPORTED,
  DEBUG_BODY_TYPE_SUPPORTED,
  DEBUG_RESPONSE_BODY_MODE_SUPPORTED,
  DebugBodyFormDataValueType,
  DEFAULT_DEBUG_AUTH_FORM_DATA,
  DEFAULT_DEBUG_PARAM_FORM_DATA,
  HTTP_METHOD_OPTION_LIST,
  PROTOCOL_SUPPORTED,
} from '../../constants';
import { AuthenticationView, DebugFormDataView, DebugParamsView } from '.';
import { debugRoute } from '../../service';
import styles from './index.less';
import type { Monaco } from '@monaco-editor/react';
import Editor from '@monaco-editor/react';

const { Option } = Select;
const { TabPane } = Tabs;

const DebugDrawView: React.FC<RouteModule.DebugDrawProps> = (props) => {
  const { formatMessage } = useIntl();
  const [httpMethod, setHttpMethod] = useState(
    props.data.method !== undefined && props.data.method.length !== 0
      ? props.data.methods[0]
      : HTTP_METHOD_OPTION_LIST[0],
  );
  const [requestProtocol, setRequestProtocol] = useState(PROTOCOL_SUPPORTED[0]);
  const [pathWildcardRewriteVisible, setPathWildcardRewriteVisible] = useState(false);
  const [pathWildcardRewrite, setPathWildcardRewrite] = useState('');
  const [pathWildcardRewritePrefix, setPathWildcardRewritePrefix] = useState('');
  const [showBodyTab, setShowBodyTab] = useState(false);
  const [queryForm] = Form.useForm();
  const [urlencodedForm] = Form.useForm();
  const [formDataForm] = Form.useForm();
  const [authForm] = Form.useForm();
  const [headerForm] = Form.useForm();
  const [response, setResponse] = useState<RouteModule.debugResponse | null>();
  const [loading, setLoading] = useState(false);
  const [body, setBody] = useState('');
  const [height, setHeight] = useState(50);
  const [bodyType, setBodyType] = useState('none');
  const methodWithoutBody = ['GET', 'HEAD'];
  const [responseBodyMode, setResponseBodyMode] = useState(
    DEBUG_RESPONSE_BODY_MODE_SUPPORTED[0].mode,
  );
  const [bodyMode, setBodyCodeMode] = useState(DEBUG_BODY_MODE_SUPPORTED[0].mode);

  enum DebugBodyType {
    None = 0,
    FormUrlencoded,
    FormData,
    RawInput,
  }

  const resetForms = () => {
    queryForm.setFieldsValue(DEFAULT_DEBUG_PARAM_FORM_DATA);
    urlencodedForm.setFieldsValue(DEFAULT_DEBUG_PARAM_FORM_DATA);
    formDataForm.setFieldsValue(DEFAULT_DEBUG_PARAM_FORM_DATA);
    headerForm.setFieldsValue(DEFAULT_DEBUG_PARAM_FORM_DATA);
    authForm.setFieldsValue(DEFAULT_DEBUG_AUTH_FORM_DATA);
    setResponse(null);
    setBodyType(DEBUG_BODY_TYPE_SUPPORTED[DebugBodyType.None]);
  };

  useEffect(() => {
    resetForms();
  }, []);

  const transformBodyParamsFormData = () => {
    if (methodWithoutBody.includes(httpMethod)) {
      return {
        bodyFormData: undefined,
      };
    }

    switch (bodyType) {
      case DEBUG_BODY_TYPE_SUPPORTED[DebugBodyType.FormUrlencoded]: {
        let transformFormUrlencoded: string[] = [];
        const FormUrlencodedData: RouteModule.debugRequestParamsFormData[] = urlencodedForm.getFieldsValue()
          .params;
        transformFormUrlencoded = (FormUrlencodedData || [])
          .filter((data) => data && data.check)
          .map((data) => {
            return `${data.key}=${data.value}`;
          });

        return {
          bodyFormData: transformFormUrlencoded.join('&'),
          header: {
            'Content-type': ['application/x-www-form-urlencoded;charset=UTF-8'],
          },
        };
      }
      case DEBUG_BODY_TYPE_SUPPORTED[DebugBodyType.RawInput]: {
        let contentType = [''];
        switch (bodyMode) {
          case DEBUG_BODY_MODE_SUPPORTED[0].mode:
            contentType = ['application/json;charset=UTF-8'];
            break;
          case DEBUG_BODY_MODE_SUPPORTED[1].mode:
            contentType = ['text/plain;charset=UTF-8'];
            break;
          case DEBUG_BODY_MODE_SUPPORTED[2].mode:
            contentType = ['application/xml;charset=UTF-8'];
            break;
          default:
            break;
        }

        return {
          bodyFormData: body,
          header: {
            'Content-type': contentType,
          },
        };
      }
      case DEBUG_BODY_TYPE_SUPPORTED[DebugBodyType.FormData]: {
        const transformFormData = new FormData();
        const formDataData: RouteModule.debugRequestParamsFormData[] = formDataForm.getFieldsValue()
          .params;

        (formDataData || [])
          .filter((data) => data && data.check)
          .forEach((data) => {
            if (data.type === DebugBodyFormDataValueType.File) {
              transformFormData.append(data.key, data.value.originFileObj);
            } else {
              transformFormData.append(data.key, data.value);
            }
          });
        return {
          bodyFormData: transformFormData,
        };
      }
      case DEBUG_BODY_TYPE_SUPPORTED[DebugBodyType.None]:
      default:
        return {
          bodyFormData: undefined,
        };
    }
  };

  const transformHeaderAndQueryParamsFormData = (
    formData: RouteModule.debugRequestParamsFormData[],
  ) => {
    let transformData = {};
    (formData || [])
      .filter((data) => data && data.check)
      .forEach((data) => {
        transformData = {
          ...transformData,
          [data.key]: [...(transformData[data.key] || []), data.value],
        };
      });

    return transformData;
  };

  const transformAuthFormData = (
    formData: RouteModule.authData,
    userHeaderData: any,
    formHeaderData: any,
  ) => {
    const { authType } = formData;

    switch (authType) {
      case 'basic-auth':
        return {
          ...formHeaderData,
          ...userHeaderData,
          Authorization: [`Basic ${Base64.encode(`${formData.username}:${formData.password}`)}`],
        };
      case 'jwt-auth':
        return {
          ...formHeaderData,
          ...userHeaderData,
          Authorization: [formData.Authorization],
        };
      case 'key-auth':
        return {
          ...formHeaderData,
          ...userHeaderData,
          apikey: [formData.apikey],
        };
      default:
        break;
    }

    return {
      ...formHeaderData,
      ...userHeaderData,
    };
  };

  const handleDebug = (url: string) => {
    const queryFormData = transformHeaderAndQueryParamsFormData(queryForm.getFieldsValue().params);
    const bodyFormRelateData = transformBodyParamsFormData();
    const { bodyFormData, header: bodyFormHeader } = bodyFormRelateData;
    const pureHeaderFormData = transformHeaderAndQueryParamsFormData(
      headerForm.getFieldsValue().params,
    );
    const headerFormData = transformAuthFormData(
      authForm.getFieldsValue(),
      pureHeaderFormData,
      bodyFormHeader,
    );
    const urlQueryString =
      url.indexOf('?') === -1
        ? `?${queryString.stringify(queryFormData)}`
        : `&${queryString.stringify(queryFormData)}`;

    setLoading(true);
    // TODO: grpc and websocket
    debugRoute(
      {
        online_debug_route_id: props.data.id,
        online_debug_header_params: JSON.stringify(headerFormData),
        online_debug_path: `${url}${urlQueryString}`,
        online_debug_request_protocol: requestProtocol,
        online_debug_method: httpMethod,
      },
      bodyFormData,
    )
      .then((req) => {
        setLoading(false);
        const resp: RouteModule.debugResponse = req.data;
        if (typeof resp.data !== 'string') {
          resp.data = JSON.stringify(resp.data, null, 2);
        }
        setResponse(resp);
        const contentType = resp.header['Content-Type'];
        if (contentType == null || contentType.length !== 1) {
          setResponseBodyMode('TEXT');
        } else if (contentType[0].toLowerCase().indexOf('json') !== -1) {
          setResponseBodyMode('JSON');
        } else if (contentType[0].toLowerCase().indexOf('xml') !== -1) {
          setResponseBodyMode('XML');
        } else if (contentType[0].toLowerCase().indexOf('html') !== -1) {
          setResponseBodyMode('HTML');
        } else {
          setResponseBodyMode('TEXT');
        }
      })
      .catch(() => {
        setLoading(false);
      });
  };

  const handleEditorMount = (editor: monacoEditor.editor.IStandaloneCodeEditor, monaco: Monaco) => {
    editor.onDidChangeModelDecorations(() => {
      if (!editor.getDomNode()) {
        return;
      }
      const padding = 40;
      const lineHeight = editor.getOption(monaco.editor.EditorOption.lineHeight);
      const lineCount = editor.getModel()?.getLineCount() || 1;
      setHeight(editor.getTopForLineNumber(lineCount + 1) + lineHeight + padding);
    });
  };

  return (
    <Drawer
      title={formatMessage({ id: 'page.route.onlineDebug' })}
      visible={props.visible}
      width={650}
      onClose={() => {
        props.onClose();
        setPathWildcardRewriteVisible(false);
        setPathWildcardRewrite('');
        setResponse(null);
        queryForm.setFieldsValue(DEFAULT_DEBUG_PARAM_FORM_DATA);
        urlencodedForm.setFieldsValue(DEFAULT_DEBUG_PARAM_FORM_DATA);
        formDataForm.setFieldsValue(DEFAULT_DEBUG_PARAM_FORM_DATA);
        headerForm.setFieldsValue(DEFAULT_DEBUG_PARAM_FORM_DATA);
        authForm.setFieldsValue(DEFAULT_DEBUG_AUTH_FORM_DATA);
      }}
      className={styles.routeDebugDraw}
      data-cy="debug-draw"
      destroyOnClose={true}
    >
      <Card bordered={false}>
        <PanelSection title={formatMessage({ id: 'page.route.PanelSection.title.requestLine' })}>
          <Input.Group compact>
            <Select
              placeholder={formatMessage({ id: 'page.route.Select.method' })}
              style={{ width: '20%' }}
              onChange={(value) => {
                setHttpMethod(value);
                setShowBodyTab(!(methodWithoutBody.indexOf(value) > -1));
              }}
              size="large"
              data-cy="debug-method"
            >
              {HTTP_METHOD_OPTION_LIST.map((method) => {
                return (
                  <Option
                    key={method}
                    value={method}
                    disabled={
                      props.data.methods !== undefined && !props.data.methods.includes(method)
                    }
                  >
                    {method}
                  </Option>
                );
              })}
            </Select>
            <Select
              placeholder={formatMessage({ id: 'page.route.Select.protocol' })}
              style={{ width: '20%' }}
              onChange={(value) => {
                setRequestProtocol(value);
              }}
              size="large"
              data-cy="debug-protocol"
            >
              {PROTOCOL_SUPPORTED.map((protocol) => {
                return (
                  <Option key={protocol} value={protocol}>
                    {`${protocol}://`}
                  </Option>
                );
              })}
            </Select>
            <Select
              placeholder={formatMessage({ id: 'page.route.Select.requestPath' })}
              style={{ width: '60%' }}
              onChange={(value: string) => {
                if (value.endsWith('*')) {
                  setPathWildcardRewrite('');
                  setPathWildcardRewriteVisible(true);
                  setPathWildcardRewritePrefix(value.slice(0, value.length - 1));
                } else {
                  setPathWildcardRewriteVisible(false);
                  setPathWildcardRewritePrefix(value);
                  setPathWildcardRewrite('');
                }
              }}
              size="large"
              data-cy="debug-path"
              options={
                props.data.uris
                  ? props.data.uris.map((path: string) => {
                      return {
                        label: path,
                        value: path,
                      };
                    })
                  : [
                      {
                        label: props.data.uri,
                        value: props.data.uri,
                      },
                    ]
              }
            />
            <Divider
              orientation="left"
              style={{ display: pathWildcardRewriteVisible ? 'flex' : 'none' }}
            >
              {formatMessage({ id: 'page.route.Model.inputPath' })}
            </Divider>
            <Input
              size="large"
              style={{ display: pathWildcardRewriteVisible ? 'flex' : 'none' }}
              addonBefore={pathWildcardRewritePrefix}
              onInput={(e) => {
                setPathWildcardRewrite(e.target.value);
              }}
            />
          </Input.Group>
        </PanelSection>
        <PanelSection
          title={formatMessage({ id: 'page.route.PanelSection.title.defineRequestParams' })}
          style={{ textAlign: 'left' }}
        >
          <Tabs>
            <TabPane
              data-cy="query"
              tab={formatMessage({ id: 'page.route.TabPane.queryParams' })}
              key="query"
            >
              <DebugParamsView form={queryForm} name="queryForm" />
            </TabPane>
            <TabPane
              data-cy="auth"
              tab={formatMessage({ id: 'page.route.TabPane.authentication' })}
              key="auth"
            >
              <AuthenticationView form={authForm} />
            </TabPane>
            <TabPane
              data-cy="header"
              tab={formatMessage({ id: 'page.route.TabPane.headerParams' })}
              key="header"
            >
              <DebugParamsView form={headerForm} name="headerForm" inputType="header" />
            </TabPane>
            {showBodyTab && (
              <TabPane
                data-cy="body"
                tab={formatMessage({ id: 'page.route.TabPane.bodyParams' })}
                key="body"
              >
                <Radio.Group
                  onChange={(e) => {
                    setBodyType(e.target.value);
                  }}
                  value={bodyType}
                >
                  {DEBUG_BODY_TYPE_SUPPORTED.map((type) => (
                    <Radio value={type} key={type}>
                      {type}
                    </Radio>
                  ))}
                </Radio.Group>
                {bodyType === DEBUG_BODY_TYPE_SUPPORTED[DebugBodyType.RawInput] && (
                  <Select
                    size="small"
                    onChange={(value) => {
                      setBodyCodeMode(value);
                    }}
                    style={{ width: 100 }}
                    defaultValue={bodyMode}
                  >
                    {DEBUG_BODY_MODE_SUPPORTED.map((modeObj) => (
                      <Option key={modeObj.name} value={modeObj.mode}>
                        {modeObj.name}
                      </Option>
                    ))}
                  </Select>
                )}
                <div style={{ marginTop: 16 }}>
                  {bodyType === DEBUG_BODY_TYPE_SUPPORTED[DebugBodyType.FormUrlencoded] && (
                    <DebugParamsView form={urlencodedForm} name="urlencodedForm" />
                  )}

                  {bodyType === DEBUG_BODY_TYPE_SUPPORTED[DebugBodyType.FormData] && (
                    <DebugFormDataView form={formDataForm} />
                  )}

                  {bodyType === DEBUG_BODY_TYPE_SUPPORTED[DebugBodyType.RawInput] && (
                    <Form>
                      <Form.Item>
                        <Editor
                          value={body}
                          language={bodyMode.toLowerCase()}
                          onChange={(text) => {
                            if (text) {
                              setBody(text);
                            } else {
                              setBody('');
                            }
                          }}
                          height={250}
                          beforeMount={(monaco) => {
                            monaco?.languages.json.jsonDefaults.setDiagnosticsOptions({
                              validate: false,
                            });
                          }}
                          options={{
                            scrollbar: {
                              vertical: 'hidden',
                              horizontal: 'hidden',
                            },
                            wordWrap: 'on',
                            minimap: { enabled: false },
                          }}
                        />
                      </Form.Item>
                    </Form>
                  )}
                </div>
              </TabPane>
            )}
          </Tabs>
        </PanelSection>
        <Divider></Divider>
        <Button
          type="primary"
          size="large"
          onClick={() => {
            handleDebug(pathWildcardRewritePrefix + pathWildcardRewrite);
          }}
        >
          {formatMessage({ id: 'page.route.button.send' })}
        </Button>
        <PanelSection title={formatMessage({ id: 'page.route.PanelSection.title.responseResult' })}>
          <Spin tip="Loading..." spinning={loading}>
            <Tabs
              tabBarExtraContent={
                response
                  ? response.message
                  : formatMessage({ id: 'page.route.debug.showResultAfterSendRequest' })
              }
            >
              <TabPane tab={formatMessage({ id: 'page.route.TabPane.response' })} key="response">
                <Select
                  disabled={response == null}
                  value={responseBodyMode}
                  onSelect={(mode) => setResponseBodyMode(mode as string)}
                >
                  {DEBUG_RESPONSE_BODY_MODE_SUPPORTED.map((mode) => {
                    return (
                      <Option value={mode.mode} key={mode.mode}>
                        {mode.name}
                      </Option>
                    );
                  })}
                </Select>
                <CopyToClipboard
                  text={response ? response.data : ''}
                  onCopy={(_: string, result: boolean) => {
                    if (!result) {
                      notification.error({
                        message: formatMessage({ id: 'component.global.copyFail' }),
                      });
                      return;
                    }
                    notification.success({
                      message: formatMessage({ id: 'component.global.copySuccess' }),
                    });
                  }}
                >
                  <Button type="text" disabled={!response}>
                    <CopyOutlined />
                  </Button>
                </CopyToClipboard>
                <div id="monaco-response" style={{ marginTop: 16 }}>
                  <Editor
                    value={response ? response.data : ''}
                    height={height}
                    language={responseBodyMode.toLowerCase()}
                    onMount={handleEditorMount}
                    beforeMount={(monaco) => {
                      monaco?.languages.json.jsonDefaults.setDiagnosticsOptions({
                        validate: false,
                      });
                    }}
                    options={{
                      automaticLayout: true,
                      scrollbar: {
                        vertical: 'hidden',
                        horizontal: 'hidden',
                      },
                      wordWrap: 'on',
                      minimap: { enabled: false },
                      readOnly: true,
                    }}
                  />
                </div>
              </TabPane>
              <TabPane tab={formatMessage({ id: 'page.route.TabPane.header' })} key="header">
                {response &&
                  Object.keys(response.header).map((header) => {
                    return response.header[header].map((value) => {
                      return (
                        <div>
                          <b>{header}</b>: {value}
                        </div>
                      );
                    });
                  })}
              </TabPane>
            </Tabs>
          </Spin>
        </PanelSection>
      </Card>
    </Drawer>
  );
};

export default DebugDrawView;
