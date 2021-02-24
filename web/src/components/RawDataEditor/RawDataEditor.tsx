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
import { Button, Drawer, PageHeader } from 'antd';
import { useIntl } from 'umi';
import { LinkOutlined } from '@ant-design/icons';

type Props = {
  visible: boolean
};

const RawDataEditor: React.FC<Props> = ({ visible }) => {
  const { formatMessage } = useIntl();

  return (
    <div>
      <Drawer
        title='Raw Data Editor'
        placement="right"
        width={700}
        visible={true}
      >
        <PageHeader
          title=""
          ghost={false}
          extra={[
            <Button
              type="default"
              icon={<LinkOutlined />}
              onClick={() => { }}
              key={1}
            >
              Document
            </Button>,
            <Button type="primary" onClick={() => { }} key={2}>
              Format
            </Button>,
            <Button type="primary" onClick={() => { }} key={2}>
              Copy
           </Button>,
          ]}
        />
      </Drawer>

    </div>
  );
};

export default RawDataEditor;
