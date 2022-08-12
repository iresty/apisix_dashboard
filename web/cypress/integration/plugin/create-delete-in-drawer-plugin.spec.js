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

context('Delete Plugin List with the Drawer', () => {
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
    authPluginCategory: '#plugin-category-authentication',
    pluginStateSelector: '[data-cy-plugin-state]',
  };

  const data = {
    basicAuthPlugin: 'basic-auth',
  };

  beforeEach(() => {
    cy.login();
  });

  it('should visit plugin market and enable plugin then show plugin in order by state', function () {
    cy.visit('/');
    cy.contains('Plugin').click();
    cy.contains('Enable').click();

    cy.contains(data.basicAuthPlugin)
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
          cy.get(selector.monacoMode).click();
          cy.get(selector.selectDropdown).should('be.visible');
          cy.get(selector.selectJSON).click();
        }
      });

    cy.get(selector.drawer)
      .should('be.visible')
      .within(() => {
        cy.get(selector.disabledSwitcher).click();
        cy.get(selector.checkedSwitcher).should('exist');
      });

    cy.contains('button', 'Submit').click();

    cy.log('**sort by state**').wait(1000);
    cy.get(selector.authPluginCategory).within(() => {
      cy.get(selector.pluginStateSelector).then((button$) => {
        const justState = [...button$].map(
          (item) =>
            item.attributes[
              selector.pluginStateSelector.slice(1, selector.pluginStateSelector.length - 1)
            ].value,
        );
        const sortedState = [...justState].sort((a, b) => b - a);
        cy.log(justState);
        cy.log(sortedState);
        expect(sortedState, 'plugin are sorted').to.deep.equal(justState);
      });
    });

    cy.get(selector.drawer, {
      timeout,
    }).should('not.exist');
  });

  it('should verify openid-connect and authz-keycloak plugin need to configure', function () {
    cy.visit('/');
    cy.contains('Plugin').click();
    cy.contains('Enable').click();
    ['openid-connect', 'authz-keycloak'].forEach((plugin) => {
      cy.contains(plugin)
        .parents(selector.pluginCardBordered)
        .within(() => {
          cy.get('button').click({
            force: true,
          });
        });
      cy.get(selector.drawer)
        .should('be.visible')
        .within(() => {
          cy.contains('.ant-alert-warning', "Doesn't need configuration").should('not.exist');
          cy.contains('button', 'Cancel').click();
        });
    });
  });

  it('should delete the plugin with the drawer', function () {
    cy.visit('/plugin/list');
    cy.get(selector.refresh).click();
    cy.contains('button', 'Configure').click();
    cy.get(selector.drawerFooter).contains('button', 'Delete').click({
      force: true,
    });
    cy.contains('button', 'Confirm').click({
      force: true,
    });
    cy.get(selector.empty).should('be.visible');
  });
});
