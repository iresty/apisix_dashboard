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
import { PageContainer } from '@ant-design/pro-layout';
import { Card, Form, notification, Steps } from 'antd';
import React, { useEffect, useRef, useState } from 'react';
import { history, useIntl } from 'umi';

import ActionBar from '@/components/ActionBar';

import Step1 from './components/Step1';
import { create, fetchOne, update } from './service';

const Page: React.FC = (props) => {
  const [step, setStep] = useState(1);
  const [form1] = Form.useForm();
  const { formatMessage } = useIntl();
  const upstreamRef = useRef<any>();
  const [submitLoading, setSubmitLoading] = useState(false);

  useEffect(() => {
    const { id } = (props as any).match.params;

    if (id) {
      fetchOne(id).then((data) => {
        const newData = data;
        if (newData) {
          if ('nodes' in newData) {
            newData.upstream_type = 'node';
          } else {
            newData.upstream_type = 'service_discovery';
          }
        }

        form1.setFieldsValue(newData);
      });
    }
  }, []);

  const onSubmit = () => {
    setSubmitLoading(true);
    form1.validateFields().then(() => {
      const data = upstreamRef.current?.getData();
      if (!data) {
        notification.error({
          message: formatMessage({ id: 'page.upstream.other.configuration.invalid' }),
        });
        setSubmitLoading(false);
        return;
      }

      const { id } = (props as any).match.params;
      (id ? update(id, data) : create(data))
        .then(() => {
          notification.success({
            message: `${
              id
                ? formatMessage({ id: 'page.upstream.edit.upstream.successfully' })
                : formatMessage({ id: 'page.upstream.create.upstream.successfully' })
            }`,
          });
          history.replace('/upstream/list');
        })
        .finally(() => {
          setSubmitLoading(false);
        });
    });
  };

  const onStepChange = (nextStep: number) => {
    if (step === 1) {
      form1.validateFields().then(() => {
        setStep(nextStep);
      });
    } else if (nextStep === 3) {
      onSubmit();
    } else {
      setStep(nextStep);
    }
  };

  return (
    <>
      <PageContainer
        title={
          (props as any).match.params.id
            ? formatMessage({ id: 'page.upstream.configure' })
            : formatMessage({ id: 'page.upstream.create' })
        }
      >
        <Card bordered={false}>
          <Steps current={step - 1} style={{ marginBottom: 30 }}>
            <Steps.Step title={formatMessage({ id: 'page.upstream.create.basic.info' })} />
            <Steps.Step title={formatMessage({ id: 'page.upstream.create.preview' })} />
          </Steps>

          {step === 1 && <Step1 form={form1} upstreamRef={upstreamRef} neverReadonly />}
          {step === 2 && <Step1 form={form1} upstreamRef={upstreamRef} disabled />}
        </Card>
      </PageContainer>
      <ActionBar loading={submitLoading} step={step} lastStep={2} onChange={onStepChange} />
    </>
  );
};

export default Page;
