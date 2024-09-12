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
/* eslint-disable */

context('Create and Search Route', () => {
  const timeout = 2000;

  const selector = {
    name: '#name',
    description: '#desc',
    hosts_0: '#hosts_0',
    hosts_1: '#hosts_1',
    add_host: 'button[data-cy="addHost"]',
    uris_0: '#uris_0',
    labels_0_labelKey: '#labels_0_labelKey',
    labels_0_labelValue: '#labels_0_labelValue',
    nodes_0_host: '#submitNodes_0_host',
    nodes_0_port: '#submitNodes_0_port',
    nodes_0_weight: '#submitNodes_0_weight',
    nameSearchInput: '#name',
    hostSearchInput: '#host',
    pathSearchInput: '#uri',
    labelSelect_0: '.ant-select-selection-overflow',
    dropdown: '.rc-virtual-list',
    disabledSwitcher: '#disable',
    checkedSwitcher: '.ant-switch-checked',
    deleteAlert: '.ant-modal-body',
    drawerBody: '.ant-drawer-wrapper-body',
    notification: '.ant-notification-notice-message',
    notificationClose: '.anticon-close',
    expandSearch: '.ant-pro-form-collapse-button',
  };

  const data = {
    host1: '11.11.11.11',
    host2: '12.12.12.12',
    hostx: 'aa.com',
    port: '80',
    weight: 1,
    uris: '/get',
    uris0: '/get0',
    uris1: '/get1',
    uris2: '/get2',
    urisx: '/getx',
    submitSuccess: 'Submit Successfully',
    deleteRouteSuccess: 'Delete Route Successfully',
    test: 'test',
    test0: 'test0',
    test1: 'test1',
    test2: 'test2',
    testx: 'testx',
    desc0: 'desc0',
    desc1: 'desc1',
    desc2: 'desc2',
    value0: 'value0',
    label0_value0: 'label0:value0',
  };

  before(() => {
    cy.clearLocalStorageSnapshot();
    cy.login();
    cy.saveLocalStorage();
  });

  beforeEach(() => {
    cy.restoreLocalStorage();
  });

  it('should create route test0, test1, test2', function () {
    cy.visit('/');
    cy.contains('Route').click();
    cy.wait(timeout);
    for (let i = 0; i < 3; i += 1) {
      cy.contains('Create').click();
      cy.contains('Next').click().click();
      cy.get(selector.name).type(`test${i}`);
      cy.get(selector.description).type(`desc${i}`);
      cy.get(selector.hosts_0).type(data.host1);
      if (i == 2) {
        cy.get(selector.add_host).click();
        cy.get(selector.hosts_1).type(data.hostx);
      }

      cy.get(selector.uris_0).clear().type(`/get${i}`);

      // config label
      cy.contains('Manage').click();

      // eslint-disable-next-line no-loop-func
      cy.get(selector.drawerBody).within(() => {
        cy.contains('button', 'Add')
          .should('not.be.disabled')
          .click()
          .then(() => {
            cy.get(selector.labels_0_labelKey).type(`label${i}`);
            cy.get(selector.labels_0_labelValue).type(`value${i}`);
            cy.contains('Confirm').click();
          });
      });

      cy.contains('button', 'Next').should('not.be.disabled').click();
      cy.get(selector.nodes_0_host).type(data.host2, {
        timeout,
      });
      cy.wait(timeout);
      cy.get(selector.nodes_0_port).type(data.port);
      cy.get(selector.nodes_0_weight).type(data.weight);
      cy.contains('Next').click();
      cy.contains('Next').click();
      cy.contains('Submit').click();
      cy.contains(data.submitSuccess);
      cy.contains('Goto List').click();
      cy.url().should('contains', 'routes/list');
    }
  });

  it('should search the route with name', function () {
    cy.visit('/');
    cy.contains('Route').click();
    cy.wait(timeout);
    // full match
    cy.get(selector.nameSearchInput).type(data.test1);
    cy.contains('Search').click();
    cy.contains(data.test1).siblings().should('contain', data.desc1);
    cy.contains(data.test0).should('not.exist');
    cy.contains(data.test2).should('not.exist');
    // partial match
    cy.get(selector.nameSearchInput).clear().type(data.test);
    cy.contains('Search').click();
    cy.contains(data.test0).siblings().should('contain', data.desc0);
    cy.contains(data.test1).siblings().should('contain', data.desc1);
    cy.contains(data.test2).siblings().should('contain', data.desc2);
    // no match
    cy.get(selector.nameSearchInput).clear().type(data.testx);
    cy.contains('Search').click();
    cy.contains(data.test0).should('not.exist');
    cy.contains(data.test1).should('not.exist');
    cy.contains(data.test2).should('not.exist');
  });

  it('should search the route with host', function () {
    cy.visit('/');
    cy.contains('Route').click();
    cy.wait(timeout);
    // only one host
    cy.get(selector.hostSearchInput).type(data.hostx);
    cy.contains('Search').click();
    cy.contains(data.test2).siblings().should('contain', data.hostx);
    cy.contains(data.test0).should('not.exist');
    cy.contains(data.test1).should('not.exist');
    // search hosts
    cy.get(selector.hostSearchInput).clear().type(data.host1);
    cy.contains('Search').click();
    cy.contains(data.test0).siblings().should('contain', data.desc0);
    cy.contains(data.test1).siblings().should('contain', data.desc1);
    cy.contains(data.test2).siblings().should('contain', data.desc2);
    // no match host
    cy.get(selector.hostSearchInput).clear().type(data.host2);
    cy.contains('Search').click();
    cy.contains(data.test0).should('not.exist');
    cy.contains(data.test1).should('not.exist');
    cy.contains(data.test2).should('not.exist');
  });

  it('should search the route with path', function () {
    cy.visit('/');
    cy.contains('Route').click();
    cy.wait(timeout);
    // full match
    // expand search
    cy.get(selector.expandSearch).click();
    cy.get(selector.pathSearchInput).type(data.uris1);
    cy.contains('Search').click();
    cy.contains(data.uris1).should('contain', data.uris1);
    cy.contains(data.uris0).should('not.exist');
    cy.contains(data.uris2).should('not.exist');
    // partial match
    cy.get(selector.pathSearchInput).clear().type(data.uris);
    cy.contains('Search').click();
    cy.contains(data.uris0).should('contain', data.uris);
    cy.contains(data.uris1).should('contain', data.uris);
    cy.contains(data.uris2).should('contain', data.uris);
    // no match
    cy.get(selector.pathSearchInput).clear().type(data.urisx);
    cy.contains('Search').click();
    cy.contains(data.uris0).should('not.exist');
    cy.contains(data.uris1).should('not.exist');
    cy.contains(data.uris2).should('not.exist');
  });

  it('should search the route with labels', function () {
    cy.visit('/');
    cy.contains('Route').click();
    cy.wait(timeout);
    // search one label
    cy.get(selector.expandSearch).click();
    cy.get(selector.labelSelect_0).click({ timeout });
    cy.get(selector.dropdown).contains(data.value0).should('be.visible').click();
    cy.contains('Search').click();
    cy.contains(data.test0).siblings().should('contain', data.label0_value0);
    cy.contains(data.test1).should('not.exist');
    cy.contains(data.test2).should('not.exist');
  });

  it('should delete the route', function () {
    cy.visit('/routes/list');
    cy.wait(timeout);
    for (let i = 0; i < 3; i += 1) {
      cy.contains(`test${i}`).siblings().contains('More').click({ timeout });
      cy.contains('Delete').should('be.visible').click({ timeout });
      cy.get(selector.deleteAlert)
        .should('be.visible')
        .within(() => {
          cy.contains('OK').click();
        });
      cy.get(selector.notification).should('contain', data.deleteRouteSuccess);
      cy.get(selector.notificationClose).should('be.visible').click({
        force: true,
        multiple: true,
      });
      cy.wait(timeout);
    }
  });
});
