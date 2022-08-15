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

defaultSettings.overwrite(Cypress.env());

Cypress.Commands.add('login', () => {
  const { SERVE_ENV = 'dev' } = Cypress.env();
  const serveUrl = defaultSettings.serveUrlMap[SERVE_ENV];

  cy.request('POST', `${serveUrl}/apisix/admin/user/login`, {
    username: 'user',
    password: 'user',
  }).then((res) => {
    expect(res.body.code).to.equal(0);
    localStorage.setItem('token', res.body.data.token);
    // set default language
    localStorage.setItem('umi_locale', 'en-US');
    cy.log(res.body.data.token);
  });
});

const timeout = 1000;
const domSelector = {
  nameGen: (name) => `[data-cy-plugin-name="${name}"]`,
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

Cypress.Commands.add('configurePlugin', ({ name, cases }) => {
  const { shouldValid, data, type } = cases;

  cy.get('main.ant-layout-content', { timeout })
    .get(domSelector.nameGen(name), { timeout })
    .then(function (card) {
      if (name !== card.innerText) {
        return;
      }

      card.parents(domSelector.parents).within(() => {
        cy.find('button').click({
          force: true,
        });
      });

      // NOTE: wait for the Drawer to appear on the DOM
      cy.focused(domSelector.drawer).should('exist');

      cy.get(domSelector.monacoMode)
        .as('monacoMode')
        .invoke('text')
        .then((text) => {
          if (text === 'Form') {
            cy.wait(1000);
            cy.get(domSelector.monacoMode).should('be.visible').click();
            cy.find(domSelector.selectDropdown).should('be.visible');
            cy.find(domSelector.selectJSON).click();
          }
        });

      cy.get(domSelector.switch, { timeout, withinSubject: domSelector.drawer }).click({
        force: true,
      });

      cy.get(domSelector.monacoMode)
        .invoke('text')
        .then((text) => {
          if (text === 'Form') {
            // FIXME: https://github.com/cypress-io/cypress/issues/7306
            cy.wait(1000);
            cy.find(domSelector.monacoMode).should('be.visible').click();
            cy.find(domSelector.selectDropdown).should('be.visible');
            cy.find(domSelector.selectJSON).click();
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

      if (shouldValid) {
        cy.get(domSelector.drawer).should('not.exist');
      } else {
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
    expect(res.body.code).to.equal(0);
    return res;
  });
});
