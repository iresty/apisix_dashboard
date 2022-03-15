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
/* eslint-disable no-undef */

context('Batch Create Service And Delete Service', () => {
  const selector = {
    name: '#name',
    nodes_0_host: '#submitNodes_0_host',
    page_item: '.ant-pagination-item-2',
    deleteAlert: '.ant-modal-body',
    notificationCloseIcon: '.ant-notification-close-icon',
    notification: '.ant-notification-notice-message',
    table_row: '.ant-table-row',
  };

  const data = {
    createServiceSuccess: 'Create Service Successfully',
    deleteServiceSuccess: 'Delete Service Successfully',
  };

  const deleteService = (serviceName) => {
    cy.contains(serviceName).siblings().contains('Delete').click();
    cy.contains('button', 'Confirm').click();
    cy.get(selector.notification).should('contain', data.deleteServiceSuccess);
    cy.get(selector.notificationCloseIcon).click();
  };

  beforeEach(() => {
    cy.login();
  });

  it('should batch create eleven service', () => {
    Array.from({ length: 11 }).forEach((value, key) => {
      const payload = {
        name: `serviceName${key}`,
        plugins: {},
        upstream: {
          type: 'roundrobin',
          pass_host: 'pass',
          scheme: 'http',
          timeout: {
            connect: 6,
            send: 6,
            read: 6,
          },
          keepalive_pool: {
            size: 320,
            idle_timeout: 60,
            requests: 1000,
          },
          nodes: {
            '127.0.0.1': 1,
          },
        },
      };
      cy.requestWithToken({ method: 'POST', payload, url: '/apisix/admin/services' });
    });
  });

  it('should delete the service', () => {
    cy.visit('/');
    cy.contains('Service').click();
    cy.wait(1000);
    cy.get(selector.page_item).click();
    deleteService('serviceName10');
    cy.url().should('contains', '/service/list?page=1&pageSize=10');
    cy.get(selector.table_row).should((service) => {
      expect(service).to.have.length(10);
    });
    Array.from({ length: 10 }).forEach((value, key) => {
      cy.contains(`serviceName${9 - key}`)
        .prev()
        .then(function ($elem) {
          cy.requestWithToken({ method: 'DELETE', url: `/apisix/admin/services/${$elem.text()}` });
        });
    });
  });
});
