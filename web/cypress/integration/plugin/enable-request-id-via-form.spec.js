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

context('Enable request id via form submission', () => {
  const timeout = 5000;

  const selector = {
    pluginCardBordered: '.ant-card-bordered',
    drawer: '.ant-drawer-content',
    selectDropdown: '.ant-select-dropdown',
    monacoMode: "[data-cy='monaco-mode']",
    selectJSON: '.ant-select-dropdown [label=JSON]',
    disabledSwitcher: '#disable',
    checkedSwitcher: '.ant-switch-checked',
    includeInResponse: '[title="include_in_response"]',
    headerName: '[title="header_name"]',
    lineContent: '.lines-content',
    headerNameInput: '#header_name',
  };

  const data = {
    requestIdPlugin: 'request-id',
    defaultHeaderName: 'X-Request-Id',
  };

  beforeEach(() => {
    cy.login();
  });

  it('should visit plugin market and enable plugin', function () {
    cy.visit('/');
    cy.contains('Plugin').click();
    cy.contains('Enable').click();

    cy.contains(data.requestIdPlugin)
      .parents(selector.pluginCardBordered)
      .within(() => {
        cy.get('button').click({
          force: true,
        });
      });

    cy.get(selector.monacoMode)
      .invoke('text')
      .then((text) => {
        if (text === 'Form') {
          // check requestId form fields
          cy.get(selector.headerName).should('be.visible');
          cy.get(selector.headerNameInput).should('have.value', data.defaultHeaderName);
          cy.get(selector.includeInResponse).should('be.visible');

          // using form to add values
          cy.get(selector.includeInResponse).click();
          cy.get(selector.checkedSwitcher).should('exist');
          cy.get(selector.monacoMode).click();
          cy.get(selector.selectDropdown).should('be.visible');
          cy.get(selector.selectJSON).click();
          cy.get(selector.lineContent).should('contain.text', 'true');
        }
      });

    cy.get(selector.drawer)
      .should('be.visible')
      .within(() => {
        cy.get(selector.disabledSwitcher).click();
        cy.get(selector.checkedSwitcher).should('exist');
      });

    cy.contains('button', 'Submit').click();
    cy.get(selector.drawer, {
      timeout,
    }).should('not.exist');
  });
});
