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
import defaultSettings from '../../config/defaultSettings';
import 'cypress-file-upload';
import '@4tw/cypress-drag-drop';

Cypress.Commands.add('login', () => {
  const { SERVE_ENV = 'dev' } = Cypress.env();

  cy.request('POST', `${defaultSettings.serveUrlMap[SERVE_ENV]}/apisix/admin/user/login`, {
    username: 'user',
    password: 'user',
  }).then((res) => {
    expect(res.body.success).is.true;
    localStorage.setItem('token', res.body.data.token);
    // set default language
    localStorage.setItem('umi_locale', 'en-US');
  });
});

Cypress.Commands.add('configurePlugin', ({ name, cases }) => {
  const timeout = 300;
  const domSelector = {
    name: '[data-cy-plugin-name]',
    parents: '.ant-card-bordered',
    drawer_wrap: '.ant-drawer-content-wrapper',
    drawer: '.ant-drawer-content',
    switch: '#disable',
    close: '.anticon-close',
    selectDropdown: '.ant-select-dropdown',
    monacoMode: '[data-cy="monaco-mode"]',
    selectJSON: '.ant-select-dropdown [label=JSON]',
    monacoViewZones: '.view-zones',
    notification: '.ant-notification-notice-message',
  };

  const shouldValid = cases.shouldValid;
  const data = cases.data;
  const type = cases.type;

  if (type === 'consumer') {
    cy.log('consumer schema case, skipping');
    return;
  }

  cy.get(domSelector.name, { timeout }).then(function (cards) {
    let needCheck = false;
    [...cards].forEach((item) => {
      if (name === item.innerText) needCheck = true;
    });

    if (!needCheck) {
      cy.log('non global plugin, skipping');
      return;
    }

    cy.contains(name)
      .parents(domSelector.parents)
      .within(() => {
        cy.get('button').click({
          force: true,
        });
      });

    // NOTE: wait for the Drawer to appear on the DOM
    cy.focused(domSelector.drawer).should('exist');

    cy.get(domSelector.monacoMode)
      .invoke('text')
      .then((text) => {
        if (text === 'Form') {
          cy.wait(1000);
          cy.get(domSelector.monacoMode).should('be.visible');
          cy.get(domSelector.monacoMode).click();
          cy.get(domSelector.selectDropdown).should('be.visible');
          cy.get(domSelector.selectJSON).click();
        }
      });

    cy.get(domSelector.drawer, { timeout }).within(() => {
      cy.get(domSelector.switch).click({
        force: true,
      });
    });

    cy.get(domSelector.monacoMode)
      .invoke('text')
      .then((text) => {
        if (text === 'Form') {
          // FIXME: https://github.com/cypress-io/cypress/issues/7306
          cy.wait(1000);
          cy.get(domSelector.monacoMode).should('be.visible');
          cy.get(domSelector.monacoMode).click();
          cy.get(domSelector.selectDropdown).should('be.visible');
          cy.get(domSelector.selectJSON).click();
        }
      });
    // edit monaco
    cy.get(domSelector.monacoViewZones).should('exist').click({ force: true });
    cy.window().then((window) => {
      window.monacoEditor.setValue(JSON.stringify(data));

      cy.get(domSelector.drawer, { timeout }).within(() => {
        cy.contains('Submit').click({
          force: true,
        });
        cy.get(domSelector.drawer).should('not.exist');
      });
    });

    if (shouldValid === true) {
      cy.get(domSelector.drawer).should('not.exist');
    } else if (shouldValid === false) {
      cy.get(domSelector.notification).should('contain', 'Invalid plugin data');

      cy.get(domSelector.close).should('be.visible').click({
        force: true,
        multiple: true,
      });

      cy.get(domSelector.drawer, { timeout })
        .invoke('show')
        .within(() => {
          cy.contains('Cancel').click({
            force: true,
          });
        });
    }
  });
});

Cypress.Commands.add('requestWithToken', ({ method, url, payload }) => {
  const { SERVE_ENV = 'dev' } = Cypress.env();
  // Make sure the request is synchronous
  cy.request({
    method,
    url: defaultSettings.serveUrlMap[SERVE_ENV] + url,
    body: payload,
    headers: { Authorization: localStorage.getItem('token') },
  }).then((res) => {
    expect(res.body.success).to.true;
    return res;
  });
});
