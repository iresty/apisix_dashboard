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

context('text authz casbin plugin', () => {
  const timeout = 5000;

  const selector = {
    pluginCardBordered: '.ant-card-bordered',
    drawer: '.ant-drawer-content',
    selectDropdown: '.ant-select-dropdown',
    monacoMode: "[data-cy='monaco-mode']",
    selectJSON: '.ant-select-dropdown [label=JSON]',
    drawerFooter: '.ant-drawer-footer',
    disabledSwitcher: '#disable',
    checkedSwitcher: '.ant-switch-checked',
    refresh: '.anticon-reload',
    empty: '.ant-empty-normal',
    conmodelPath: '#model_path',
    conpolicyPath: '#policy_path',
    conmodel: '#model',
    conpolicy: '#policy',
    username: '#username'
  };

  const data = {
    authzCasbin: 'authz-casbin',
    conmodelPath: '#/path/to/model.conf',
    conpolicyPath: '#/path/to/policy.csv',
    conmodel: '#/to/model.conf',
    conpolicy: '#/to/policy.csv',
    conusername: '#user',
  };

  beforeEach(() => {
    cy.login();
  });

  it('should visit plugin market and enabl plugin', function () {
    cy.visit('/');
    cy.contains('Plugin').click();
    cy.contains('Enable').click();
    cy.contains(data.authzCasbin)
      .parents(selector.pluginCardBordered)
      .within(() => {
        cy.get('button').click({
          force: true,
        })
      })
    cy.get(selector.drawer)
      .should('be.visible')
      .within(() => {
        cy.get(selector.disabledSwitcher).click();
        cy.get(selector.checkedSwitcher).should('exist');
      })
    cy.contains('Config Path').click();
    cy.get(selector.conmodelPath).type(data.conmodelPath);
    cy.get(selector.conpolicyPath).type(data.conpolicyPath);
    cy.get(selector.username).type(data.conusername);
    cy.contains('Button', 'Submit').click();
    cy.contains('Button', 'Disable').click();
    cy.contains('Custom text').click();
    cy.get(selector.conmodel).type(data.conmodel);
    cy.get(selector.conpolicy).type(data.conpolicy);
    cy.get(selector.username).type(data.conusername);
    cy.contains('Button', 'Submit').click();
    cy.contains('Plugin').click();
    cy.contains('Button','Delete').click();
    cy.contains('Confirm').click();
  })
});
