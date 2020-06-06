import React, { useRef } from 'react';
import { PageHeaderWrapper } from '@ant-design/pro-layout';
import ProTable, { ProColumns, ActionType } from '@ant-design/pro-table';
import { Button, Switch, Popconfirm, notification } from 'antd';
import { history, useIntl } from 'umi';
import { PlusOutlined } from '@ant-design/icons';

import { fetchList as fetchSSLList, remove as removeSSL } from './service';

const List: React.FC = () => {
  const tableRef = useRef<ActionType>();
  const { formatMessage } = useIntl();

  const columns: ProColumns<SSLModule.ResSSL>[] = [
    {
      title: 'SNI',
      dataIndex: 'sni',
    },
    {
      title: '过期时间',
      dataIndex: 'validity_end',
      hideInSearch: true,
      render: (text) => `${new Date(Number(text) * 1000).toLocaleString()}`,
    },
    {
      title: '是否启用',
      valueType: 'option',
      render: () => (
        <>
          <Switch defaultChecked />
        </>
      ),
    },
    {
      title: formatMessage({ id: 'component.global.action' }),
      valueType: 'option',
      render: (_, record) => (
        <Popconfirm
          title="删除"
          // TODO: 确认按钮应为红色警告
          onConfirm={() =>
            removeSSL(record.id).then(() => {
              notification.success({
                message: formatMessage({ id: 'component.ssl.removeSSLSuccess' }),
              });
              /* eslint-disable no-unused-expressions */
              requestAnimationFrame(() => tableRef.current?.reload());
            })
          }
        >
          <Button type="primary" danger>
            {formatMessage({ id: 'component.global.remove' })}
          </Button>
        </Popconfirm>
      ),
    },
  ];

  return (
    <PageHeaderWrapper>
      <ProTable<SSLModule.ResSSL>
        request={(params) => fetchSSLList(params)}
        search={false}
        rowKey="id"
        columns={columns}
        actionRef={tableRef}
        toolBarRender={() => [
          <Button type="primary" onClick={() => history.push(`/ssl/create`)}>
            <PlusOutlined />
            {formatMessage({ id: 'component.global.create' })}
          </Button>,
        ]}
      />
    </PageHeaderWrapper>
  );
};

export default List;
