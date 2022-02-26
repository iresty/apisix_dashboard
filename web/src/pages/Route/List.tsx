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
import type { ReactNode } from 'react';
import React, { useRef, useEffect, useState } from 'react';
import { PageHeaderWrapper } from '@ant-design/pro-layout';
import ProTable from '@ant-design/pro-table';
import type { ProColumns, ActionType } from '@ant-design/pro-table';
import {
  Button,
  Popconfirm,
  notification,
  Tag,
  Space,
  Select,
  Radio,
  Form,
  Upload,
  Modal,
  Divider,
  Menu,
  Dropdown,
  Tooltip,
} from 'antd';
import { history, useIntl } from 'umi';
import {
  PlusOutlined,
  BugOutlined,
  ExportOutlined,
  ImportOutlined,
  DownOutlined,
} from '@ant-design/icons';
import { js_beautify } from 'js-beautify';
import yaml from 'js-yaml';
import moment from 'moment';
import { saveAs } from 'file-saver';
import querystring from 'query-string';
import { omit } from 'lodash';

import { DELETE_FIELDS } from '@/constants';
import { timestampToLocaleString } from '@/helpers';
import type { RcFile } from 'antd/lib/upload';
import {
  update,
  create,
  fetchList,
  remove,
  fetchLabelList,
  updateRouteStatus,
  exportRoutes,
  importRoutes,
} from './service';
import { DebugDrawView } from './components/DebugViews';
import { RawDataEditor } from '@/components/RawDataEditor';
import { EXPORT_FILE_MIME_TYPE_SUPPORTED } from './constants';

const { OptGroup, Option } = Select;

const Page: React.FC = () => {
  const ref = useRef<ActionType>();
  const { formatMessage } = useIntl();
  const [exportFileTypeForm] = Form.useForm();

  enum RouteStatus {
    Offline = 0,
    Publish,
  }

  enum ExportFileType {
    JSON = 0,
    YAML,
  }

  const [labelList, setLabelList] = useState<LabelList>({});
  const [selectedRowKeys, setSelectedRowKeys] = useState<string[]>([]);
  const [uploadFileList, setUploadFileList] = useState<RcFile[]>([]);
  const [showImportModal, setShowImportModal] = useState(false);
  const [visible, setVisible] = useState(false);
  const [rawData, setRawData] = useState<Record<string, any>>({});
  const [id, setId] = useState('');
  const [editorMode, setEditorMode] = useState<'create' | 'update'>('create');
  const [paginationConfig, setPaginationConfig] = useState({ pageSize: 10, current: 1 });
  const [debugDrawVisible, setDebugDrawVisible] = useState(false);

  const savePageList = (page = 1, pageSize = 10) => {
    history.replace(`/routes/list?page=${page}&pageSize=${pageSize}`);
  };

  useEffect(() => {
    fetchLabelList().then(setLabelList);
  }, []);

  useEffect(() => {
    const { page = 1, pageSize = 10 } = querystring.parse(window.location.search);
    setPaginationConfig({ pageSize: Number(pageSize), current: Number(page) });
  }, [window.location.search]);

  const rowSelection: any = {
    selectedRowKeys,
    onChange: (currentSelectKeys: string[]) => {
      setSelectedRowKeys(currentSelectKeys);
    },
    preserveSelectedRowKeys: true,
  };

  const handleTableActionSuccessResponse = (msgTip: string) => {
    notification.success({
      message: msgTip,
    });

    ref.current?.reload();
  };

  const handlePublishOffline = (rid: string, status: RouteModule.RouteStatus) => {
    updateRouteStatus(rid, status).then(() => {
      const actionName = status
        ? formatMessage({ id: 'page.route.publish' })
        : formatMessage({ id: 'page.route.offline' });
      handleTableActionSuccessResponse(
        `${actionName} ${formatMessage({
          id: 'menu.routes',
        })} ${formatMessage({ id: 'component.status.success' })}`,
      );
    });
  };

  const handleExport = (exportFileType: ExportFileType) => {
    exportRoutes(selectedRowKeys.join(',')).then((resp) => {
      let exportFile: string;
      let exportFileName = `APISIX_routes_${moment().format('YYYYMMDDHHmmss')}`;

      switch (exportFileType) {
        case ExportFileType.YAML:
          exportFile = yaml.dump(resp.data);
          exportFileName = `${exportFileName}.${ExportFileType[
            ExportFileType.YAML
          ].toLocaleLowerCase()}`;
          break;
        case ExportFileType.JSON:
        default:
          exportFile = js_beautify(JSON.stringify(resp.data), {
            indent_size: 2,
          });
          exportFileName = `${exportFileName}.${ExportFileType[
            ExportFileType.JSON
          ].toLocaleLowerCase()}`;
          break;
      }

      const blob = new Blob([exportFile], {
        type: EXPORT_FILE_MIME_TYPE_SUPPORTED[exportFileType],
      });

      saveAs(window.URL.createObjectURL(blob), exportFileName);
    });
  };

  const handleImport = () => {
    const formData = new FormData();
    if (!uploadFileList[0]) {
      notification.warn({
        message: formatMessage({ id: 'page.route.button.selectFile' }),
      });
      return;
    }
    formData.append('file', uploadFileList[0]);

    importRoutes(formData).then(() => {
      handleTableActionSuccessResponse(
        `${formatMessage({ id: 'page.route.button.importOpenApi' })} ${formatMessage({
          id: 'component.status.success',
        })}`,
      );
      setShowImportModal(false);
    });
  };

  const ListToolbar = () => {
    const tools = [
      {
        name: formatMessage({ id: 'page.route.pluginTemplateConfig' }),
        icon: <PlusOutlined />,
        onClick: () => {
          history.push('/plugin-template/list');
        },
      },
      {
        name: formatMessage({ id: 'component.global.data.editor' }),
        icon: <PlusOutlined />,
        onClick: () => {
          setVisible(true);
          setEditorMode('create');
          setRawData({});
        },
      },
      {
        name: formatMessage({ id: 'page.route.button.importOpenApi' }),
        icon: <ImportOutlined />,
        onClick: () => {
          setUploadFileList([]);
          setShowImportModal(true);
        },
      },
      {
        name: formatMessage({ id: 'page.route.onlineDebug' }),
        icon: <BugOutlined />,
        onClick: () => {
          setDebugDrawVisible(true);
        },
      },
    ];

    return (
      <Dropdown
        overlay={
          <Menu>
            {tools.map((item) => (
              <Menu.Item key={item.name} onClick={item.onClick}>
                {item.icon}
                {item.name}
              </Menu.Item>
            ))}
          </Menu>
        }
      >
        <Button type="dashed">
          <DownOutlined /> {formatMessage({ id: 'menu.advanced-feature' })}
        </Button>
      </Dropdown>
    );
  };

  const RecordActionDropdown: React.FC<{ record: any }> = ({ record }) => {
    const tools: {
      name: string;
      onClick: () => void;
      icon?: ReactNode;
    }[] = [
      {
        name: formatMessage({ id: 'component.global.view' }),
        onClick: () => {
          setId(record.id);
          setRawData(omit(record, DELETE_FIELDS));
          setVisible(true);
          setEditorMode('update');
        },
      },
      {
        name: formatMessage({ id: 'component.global.duplicate' }),
        onClick: () => {
          history.push(`/routes/${record.id}/duplicate`);
        },
      },
      {
        name: formatMessage({ id: 'component.global.delete' }),
        onClick: () => {
          Modal.confirm({
            type: 'warning',
            title: formatMessage({ id: 'component.global.popconfirm.title.delete' }),
            content: (
              <>
                {formatMessage({ id: 'component.global.name' })} - {record.name}
                <br />
                ID - {record.id}
              </>
            ),
            onOk: () => {
              return remove(record.id!).then(() => {
                handleTableActionSuccessResponse(
                  `${formatMessage({ id: 'component.global.delete' })} ${formatMessage({
                    id: 'menu.routes',
                  })} ${formatMessage({ id: 'component.status.success' })}`,
                );
              });
            },
          });
        },
      },
      {
        name: formatMessage({ id: 'component.global.debug' }),
        onClick: () => {
          setId(record.id);
          setRawData(record);
          setDebugDrawVisible(true);
        },
      },
    ];

    return (
      <Dropdown
        overlay={
          <Menu>
            {tools.map((item) => (
              <Menu.Item key={item.name} onClick={item.onClick}>
                {item.icon && item.icon}
                {item.name}
              </Menu.Item>
            ))}
          </Menu>
        }
      >
        <Button type="dashed">
          <DownOutlined />
          {formatMessage({ id: 'menu.more' })}
        </Button>
      </Dropdown>
    );
  };

  const ListFooter: React.FC = () => {
    return (
      <Popconfirm
        title={
          <Form form={exportFileTypeForm} initialValues={{ fileType: ExportFileType.JSON }}>
            <div style={{ marginBottom: 8 }}>
              {formatMessage({ id: 'page.route.exportRoutesTips' })}
            </div>
            <Form.Item name="fileType" noStyle>
              <Radio.Group>
                <Radio value={ExportFileType.JSON}>Json</Radio>
                <Radio value={ExportFileType.YAML}>Yaml</Radio>
              </Radio.Group>
            </Form.Item>
          </Form>
        }
        onConfirm={() => {
          handleExport(exportFileTypeForm.getFieldValue('fileType'));
        }}
        okText={formatMessage({ id: 'component.global.confirm' })}
        cancelText={formatMessage({ id: 'component.global.cancel' })}
        disabled={selectedRowKeys.length === 0}
      >
        <Button type="primary" disabled={selectedRowKeys.length === 0}>
          <ExportOutlined />
          {formatMessage({ id: 'page.route.button.exportOpenApi' })}
        </Button>
      </Popconfirm>
    );
  };
  const tagStyle = {
    maxWidth: '200px',
    overflow: 'hidden',
    WhiteSpace: 'nowrap',
    textOverflow: 'ellipsis',
  };
  const columns: ProColumns<RouteModule.ResponseBody>[] = [
    {
      title: formatMessage({ id: 'component.global.name' }),
      dataIndex: 'name',
      fixed: 'left',
    },
    {
      title: formatMessage({ id: 'component.global.id' }),
      hideInSearch: true,
      dataIndex: 'id',
      width: 200,
    },
    {
      title: formatMessage({ id: 'page.route.host' }),
      hideInSearch: true,
      width: 224,
      render: (_, record) => {
        const list = record.hosts || (record.host && [record.host]) || [];

        return list.map((item) => (
          <Tooltip placement="topLeft" title={item}>
            <Tag key={item} color="geekblue" style={tagStyle}>
              {item}
            </Tag>
          </Tooltip>
        ));
      },
    },
    {
      title: formatMessage({ id: 'page.route.path' }),
      dataIndex: 'uri',
      width: 224,
      render: (_, record) => {
        const list = record.uris || (record.uri && [record.uri]) || [];

        return list.map((item) => (
          <Tooltip placement="topLeft" title={item}>
            <Tag key={item} color="geekblue" style={tagStyle}>
              {item}
            </Tag>
          </Tooltip>
        ));
      },
    },
    {
      title: formatMessage({ id: 'component.global.description' }),
      dataIndex: 'desc',
      hideInSearch: true,
      ellipsis: true,
      width: 200,
    },
    {
      title: formatMessage({ id: 'component.global.labels' }),
      dataIndex: 'labels',
      render: (_, record) => {
        return Object.keys(record.labels || {})
          .filter((item) => item !== 'API_VERSION')
          .map((item) => (
            <Tag key={Math.random().toString(36).slice(2)}>
              {item}:{record.labels[item]}
            </Tag>
          ));
      },
      renderFormItem: (_, { type }) => {
        if (type === 'form') {
          return null;
        }

        return (
          <Select
            mode="tags"
            style={{ width: '100%' }}
            placeholder={formatMessage({ id: 'component.global.pleaseChoose' })}
            tagRender={(props) => {
              const { value, closable, onClose } = props;
              return (
                <Tag closable={closable} onClose={onClose} style={{ marginRight: 3 }}>
                  {value}
                </Tag>
              );
            }}
          >
            {Object.keys(labelList)
              .filter((item) => item !== 'API_VERSION')
              .map((key) => {
                return (
                  <OptGroup label={key} key={Math.random().toString(36).slice(2)}>
                    {(labelList[key] || []).map((value: string) => (
                      <Option key={Math.random().toString(36).slice(2)} value={`${key}:${value}`}>
                        {' '}
                        {value}{' '}
                      </Option>
                    ))}
                  </OptGroup>
                );
              })}
          </Select>
        );
      },
    },
    {
      title: formatMessage({ id: 'component.global.version' }),
      dataIndex: 'API_VERSION',
      width: 100,
      render: (_, record) => {
        return Object.keys(record.labels || {})
          .filter((item) => item === 'API_VERSION')
          .map((item) => record.labels[item]);
      },
      renderFormItem: (_, { type }) => {
        if (type === 'form') {
          return null;
        }

        return (
          <Select
            style={{ width: '100%' }}
            placeholder={formatMessage({ id: 'component.global.pleaseChoose' })}
            allowClear
          >
            {Object.keys(labelList)
              .filter((item) => item === 'API_VERSION')
              .map((key) => {
                return (
                  <OptGroup label={key} key={Math.random().toString(36).slice(2)}>
                    {(labelList[key] || []).map((value: string) => (
                      <Option key={Math.random().toString(36).slice(2)} value={`${key}:${value}`}>
                        {' '}
                        {value}{' '}
                      </Option>
                    ))}
                  </OptGroup>
                );
              })}
          </Select>
        );
      },
    },
    {
      title: formatMessage({ id: 'page.route.status' }),
      dataIndex: 'status',
      width: 100,
      render: (_, record) => (
        <>
          {record.status ? (
            <Tag color="green">{formatMessage({ id: 'page.route.published' })}</Tag>
          ) : (
            <Tag color="red">{formatMessage({ id: 'page.route.unpublished' })}</Tag>
          )}
        </>
      ),
      renderFormItem: (_, { type }) => {
        if (type === 'form') {
          return null;
        }

        return (
          <Select
            style={{ width: '100%' }}
            placeholder={`${formatMessage({ id: 'page.route.unpublished' })}/${formatMessage({
              id: 'page.route.published',
            })}`}
            allowClear
          >
            <Option key={RouteStatus.Offline} value={RouteStatus.Offline}>
              {formatMessage({ id: 'page.route.unpublished' })}
            </Option>
            <Option key={RouteStatus.Publish} value={RouteStatus.Publish}>
              {formatMessage({ id: 'page.route.published' })}
            </Option>
          </Select>
        );
      },
    },
    {
      title: formatMessage({ id: 'component.global.updateTime' }),
      dataIndex: 'update_time',
      hideInSearch: true,
      width: 200,
      render: (text) => timestampToLocaleString(text as number),
    },
    {
      title: formatMessage({ id: 'component.global.operation' }),
      valueType: 'option',
      fixed: 'right',
      hideInSearch: true,
      width: 240,
      render: (_, record) => (
        <>
          <Space align="baseline">
            {!record.status ? (
              <Button
                type="primary"
                onClick={() => {
                  handlePublishOffline(record.id, RouteStatus.Publish);
                }}
              >
                {formatMessage({ id: 'page.route.publish' })}
              </Button>
            ) : null}
            {record.status ? (
              <Popconfirm
                title={formatMessage({ id: 'page.route.popconfirm.title.offline' })}
                onConfirm={() => {
                  handlePublishOffline(record.id, RouteStatus.Offline);
                }}
                okButtonProps={{
                  danger: true,
                }}
                okText={formatMessage({ id: 'component.global.confirm' })}
                cancelText={formatMessage({ id: 'component.global.cancel' })}
              >
                <Button type="primary" danger disabled={Boolean(!record.status)}>
                  {formatMessage({ id: 'page.route.offline' })}
                </Button>
              </Popconfirm>
            ) : null}
            <Button type="primary" onClick={() => history.push(`/routes/${record.id}/edit`)}>
              {formatMessage({ id: 'component.global.edit' })}
            </Button>
            <RecordActionDropdown record={record} />
          </Space>
        </>
      ),
    },
  ];

  return (
    <PageHeaderWrapper
      title={formatMessage({ id: 'page.route.list' })}
      content={formatMessage({ id: 'page.route.list.description' })}
    >
      <ProTable<RouteModule.ResponseBody>
        actionRef={ref}
        rowKey="id"
        columns={columns}
        request={fetchList}
        pagination={{
          onChange: (page, pageSize?) => savePageList(page, pageSize),
          pageSize: paginationConfig.pageSize,
          current: paginationConfig.current,
        }}
        search={{
          searchText: formatMessage({ id: 'component.global.search' }),
          resetText: formatMessage({ id: 'component.global.reset' }),
        }}
        toolBarRender={() => [
          <Button type="primary" onClick={() => history.push(`/routes/create`)}>
            <PlusOutlined />
            {formatMessage({ id: 'component.global.create' })}
          </Button>,
          <ListToolbar />,
        ]}
        rowSelection={rowSelection}
        footer={() => <ListFooter />}
        tableAlertRender={false}
        scroll={{ x: 1300 }}
      />
      <DebugDrawView
        visible={debugDrawVisible}
        data={rawData}
        onClose={() => {
          setDebugDrawVisible(false);
        }}
      />
      <RawDataEditor
        visible={visible}
        type="route"
        readonly={false}
        data={rawData}
        onClose={() => {
          setVisible(false);
        }}
        onSubmit={(data: any) => {
          (editorMode === 'create' ? create(data, 'RawData') : update(id, data, 'RawData')).then(
            () => {
              setVisible(false);
              ref.current?.reload();
            },
          );
        }}
      />
      <Modal
        title={formatMessage({ id: 'page.route.button.importOpenApi' })}
        visible={showImportModal}
        okText={formatMessage({ id: 'component.global.confirm' })}
        onOk={handleImport}
        onCancel={() => {
          setShowImportModal(false);
        }}
      >
        <Upload
          fileList={uploadFileList as any}
          beforeUpload={(file) => {
            setUploadFileList([file]);
            return false;
          }}
          onRemove={() => {
            setUploadFileList([]);
          }}
        >
          <Button type="primary" icon={<ImportOutlined />}>
            {formatMessage({ id: 'page.route.button.selectFile' })}
          </Button>
        </Upload>
        <Divider />
        <div>
          <p>{formatMessage({ id: 'page.route.instructions' })}:</p>
          <p>
            <a
              href="https://apisix.apache.org/docs/dashboard/IMPORT_OPENAPI_USER_GUIDE"
              target="_blank"
            >
              1.{' '}
              {`${formatMessage({ id: 'page.route.import' })} ${formatMessage({
                id: 'page.route.instructions',
              })}`}
            </a>
          </p>
        </div>
      </Modal>
    </PageHeaderWrapper>
  );
};

export default Page;
