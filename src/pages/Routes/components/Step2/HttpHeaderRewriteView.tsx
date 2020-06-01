import React, { useState } from 'react';
import Form from 'antd/es/form';
import { Button, Table, Space, Modal, Input } from 'antd';
import TextArea from 'antd/lib/input/TextArea';

import PanelSection from '../PanelSection';

interface Props extends RouteModule.Data {}

const HttpHeaderRewriteView: React.FC<Props> = ({ data, disabled, onChange }) => {
  const { upstream_header } = data.step2Data;
  const [visible, setVisible] = useState(false);
  const [modalForm] = Form.useForm();
  const [mode, setMode] = useState<RouteModule.ModalType>('create');

  const handleEdit = (record: RouteModule.UpstreamHeader) => {
    setMode('edit');
    setVisible(true);
    modalForm.setFieldsValue(record);
  };
  const handleRemove = (key: string) => {
    onChange({ upstream_header: upstream_header.filter((item) => item.key !== key) });
  };

  const columns = [
    {
      title: '请求头',
      dataIndex: 'header_name',
      key: 'header_name',
    },
    {
      title: '值',
      dataIndex: 'header_value',
      key: 'header_value',
    },
    {
      title: '描述',
      dataIndex: 'header_desc',
      key: 'header_desc',
    },
    disabled
      ? {}
      : {
          title: '操作',
          key: 'action',
          render: (_: any, record: RouteModule.UpstreamHeader) => (
            <Space size="middle">
              <a
                onClick={() => {
                  handleEdit(record);
                }}
              >
                编辑
              </a>
              <a
                onClick={() => {
                  handleRemove(record.key);
                }}
              >
                移除
              </a>
            </Space>
          ),
        },
  ];

  const renderModal = () => {
    const handleOk = () => {
      modalForm.validateFields().then((value) => {
        if (mode === 'edit') {
          const key = modalForm.getFieldValue('key');
          const newUpstreamHeader = upstream_header.concat();
          const findIndex = newUpstreamHeader.findIndex((item) => item.key === key);
          newUpstreamHeader[findIndex] = { ...(value as RouteModule.UpstreamHeader), key };
          onChange({ upstream_header: newUpstreamHeader, key });
        } else {
          onChange({
            upstream_header: upstream_header.concat({
              ...(value as RouteModule.UpstreamHeader),
              key: Math.random().toString(36).slice(2),
            }),
          });
        }
        modalForm.resetFields();
        setVisible(false);
      });
    };

    return (
      <Modal
        title={mode === 'edit' ? '编辑' : '新增'}
        centered
        visible={visible}
        onOk={handleOk}
        onCancel={() => {
          setVisible(false);
          modalForm.resetFields();
        }}
        destroyOnClose
      >
        <Form form={modalForm} labelCol={{ span: 4 }} wrapperCol={{ span: 20 }}>
          <Form.Item
            label="请求头"
            name="header_name"
            rules={[{ required: true, message: '请输入请求头' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            label="值"
            name="header_value"
            rules={[{ required: true, message: '请输入值' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item label="描述" name="header_desc">
            <TextArea />
          </Form.Item>
        </Form>
      </Modal>
    );
  };

  return (
    <PanelSection title="HTTP 头改写">
      {!disabled && (
        <Button
          onClick={() => {
            setMode('create');
            setVisible(true);
          }}
          type="primary"
          style={{
            marginBottom: 16,
          }}
        >
          新增
        </Button>
      )}
      <Table key="table" bordered dataSource={upstream_header} columns={columns} />
      {renderModal()}
    </PanelSection>
  );
};

export default HttpHeaderRewriteView;
