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
// import menuLocaleUS from '../../../src/locales/en-US/menu';
// import routeLocaleUS from '../../../src/pages/Route/locales/en-US';
import defaultSettings from '../../../config/defaultSettings';

context('Online debug', () => {
  const { SERVE_ENV = 'dev' } = Cypress.env();

  const selector = {
    name: '#name',
    path: '#uris_0',
    nodes_0_host: '#nodes_0_host',
    nodes_0_port: '#nodes_0_port',
    nodes_0_weight: '#nodes_0_weight',
    input: ':input',
    upstreamSelector: '[data-cy=upstream_selector]',
    nameSelector: '[title=Name]',
    refresh: '.anticon-reload',
    debugUri: '#debugUri',
    debugDraw: '[data-cy=debug-draw]',
    debugMethod: '[data-cy=debug-method]',
    debugProtocol: '[data-cy=debug-protocol]',
    debugFormDataKey0: '#dynamic_form_data_item_params_0_key',
    debugFormDataType0: '[data-cy=debug-formdata-type-0]',
    debugFormDataValue0: '#dynamic_form_data_item_params_0_value',
    debugFormDataFileButton0: '[data-cy=debug-upload-btn-0]',
    codeMirrorResp: '#codeMirror-response',
    headerDataKey0: '#headerForm_params_0_key',
    headerDataValue0: '#headerForm_params_0_value',
    queryDataKey0: '#queryForm_params_0_key',
    queryDataValue0: '#queryForm_params_0_value',
    queryTab: '#rc-tabs-0-tab-query',
    authTab: '#rc-tabs-0-tab-auth',
    headerTab: '#rc-tabs-0-tab-header',
    bodyTab: '#rc-tabs-0-tab-body',
    jwtAuth: '[data-cy=jwt-auth]',
    jwtTokenInput: '#authForm_Authorization',
    drawerClose: '.ant-drawer-close',
    notification: '.ant-notification-notice-message',
    notificationCloseIcon: '.ant-notification-close-icon',
    deleteAlert: '.ant-modal-body',
  };

  const data = {
    routeName1: 'test_route1',
    routeName2: 'test_route2',
    path1: '/hello',
    path2: '/hello/*',
    submitSuccess: 'Submit Successfully',
    createUpstreamSuccess: 'Create Upstream Successfully',
    upstreamName: 'test_upstream',

    input: ':input',
    ip1: '127.0.0.1',
    ip2: '127.0.0.2',
    port1: '80',
    port2: '81',
    weight1: 1,
    weight2: 2,
    deleteRouteSuccess: 'Delete Route Successfully',
    validUris: [
      'localhost:9000/get',
      '127.0.0.1:9000/get',
      'www.baidu.com/get',
      'baidu.com/get',
      'test-host.com/get',
      'test-host.com/aBCdsf/defsdf456eweWQE',
      'localhost:9000/deps-get',
      'localhost:9000/deps.get',
      'localhost:9000/deps_get',
      'localhost:9000/get?search=1',
      'localhost:9000/get?search=v1,sd&number=-1',
      'localhost:9000/get?search=1+1',
      'localhost:9000/api/commands/submit.html#Requirements?test=apisix.com',
      'localhost:9000/js6/main.jsp?sid=pARQZYHABxkSVdeMvXAAEtfJKbWQocOA&df=mail126_mailmaster#module=mbox.ListModule%7C%7B',
    ],
    invalidUrls: ['000'],
    postFileUrl: `${
      defaultSettings.serveUrlMap[SERVE_ENV].split('//')[1]
    }/apisix/admin/import/routes`,
    routeOptUrl: `${defaultSettings.serveUrlMap[SERVE_ENV].split('//')[1]}/apisix/admin/routes`,
    uploadFile: '../../../api/test/testdata/import/default.json',
    headerAuthorizationKey: 'Authorization',
    queryKey0: 'name',
  };

  beforeEach(() => {
    cy.login();
    cy.fixture('route-json-data.json').as('routeData');
    cy.intercept('/apisix/admin/debug-request-forwarding').as('DebugAPI');
  });

  it('should create route without upstream ', function () {
    cy.visit('/');
    cy.get('[role=menu]')
      .should('be.visible')
      .within(() => {
        cy.contains('Route').click();
      });
    cy.contains('Create').click();

    cy.contains('Next').click().click();
    cy.get(selector.name).type(data.routeName1);
    cy.get(selector.path).clear().type(data.path1);

    cy.contains('Next').click();

    cy.get(selector.nodes_0_host).type(data.ip1);
    cy.get(selector.nodes_0_port).type(data.port1);
    cy.get(selector.nodes_0_weight).type(data.weight1);
    cy.contains('Next').click();
    cy.contains('Next').click();
    cy.get(selector.input).should('be.disabled');
    cy.contains('Submit').click();
    cy.contains(data.submitSuccess).should('be.visible');
    cy.contains('Goto List').click();
    cy.url().should('contains', 'routes/list');
  });

  it('should create upstream', function () {
    cy.visit('/');
    cy.contains('Upstream').click();
    cy.contains('Create').click();

    cy.get(selector.name).type(data.upstreamName);

    cy.get(selector.nodes_0_host).type(data.ip2);
    cy.get(selector.nodes_0_port).clear().type(data.port2);
    cy.get(selector.nodes_0_weight).clear().type(data.weight2);

    cy.contains('Next').click();
    cy.contains('Submit').click();
    cy.get(selector.notification).should('contain', data.createUpstreamSuccess);
    cy.url().should('contains', 'upstream/list');
  });

  it('should create route with upstream ', function () {
    cy.visit('/');
    cy.get('[role=menu]')
      .should('be.visible')
      .within(() => {
        cy.contains('Route').click();
      });
    cy.contains('Create').click();

    cy.contains('Next').click().click();
    cy.get(selector.name).type(data.routeName2);
    cy.get(selector.path).clear().type(data.path2);
    cy.contains('Next').click();
    cy.get(selector.upstreamSelector).click();
    cy.contains(data.upstreamName).click();
    cy.get(selector.input).should('be.disabled');
    cy.contains('Next').click().click();
    cy.contains('Submit').click();
    cy.contains(data.submitSuccess).should('be.visible');
    cy.contains('Goto List').click();
    cy.url().should('contains', 'routes/list');
  });

  it('debug draw data should be clear when change route', function () {
    cy.visit('/');
    cy.contains('Route').click();
    cy.get(selector.nameSelector).type(data.routeName1);
    cy.contains('Search').click();
    cy.contains(data.routeName1).siblings().contains('More').click();
    cy.contains('Debug').click();
    cy.get(selector.debugDraw).should('be.visible');
    cy.get(selector.debugMethod).click();
    cy.get('[title=POST]').click();
    cy.get(selector.debugProtocol).click();
    cy.contains('http://').click();

    cy.visit('/');
    cy.contains('Route').click();
    cy.get(selector.nameSelector).type(data.routeName2);
    cy.contains('Search').click();
    cy.contains(data.routeName2).siblings().contains('More').click();
    cy.contains('Debug').click();
    cy.get(selector.debugDraw).should('be.visible');
    cy.get(selector.debugMethod).should('contain', 'method');
    cy.get(selector.debugProtocol).should('contain', 'protocol');
  });
  // it('method and path of debug draw should be set by route chosen', function () {
  //   cy.visit('/');
  //   cy.contains('Route').click();
  //   cy.get(selector.nameSelector).type(data.routeName1);
  //   cy.contains('Search').click();
  //   cy.contains(data.routeName1).siblings().contains('More').click();
  //   cy.contains('Debug').click();
  //   cy.get(selector.debugDraw).should('be.visible');
  //   cy.get(selector.debugMethod).click();
  //   cy.get('[title=POST]').click();
  //   cy.get(selector.debugProtocol).click();
  //   cy.contains('http://').click();
  //
  //   cy.visit('/');
  //   cy.contains('Route').click();
  //   cy.get(selector.nameSelector).type(data.routeName2);
  //   cy.contains('Search').click();
  //   cy.contains(data.routeName2).siblings().contains('More').click();
  //   cy.contains('Debug').click();
  //   cy.get(selector.debugDraw).should('be.visible');
  //   cy.get(selector.debugMethod).should('contain','method')
  //   cy.get(selector.debugProtocol).should('contain','protocol')
  //   })
  //
  // it('path with * should be rewrite ', function () {
  //   cy.visit('/');
  //   cy.contains('Route').click();
  //   cy.get(selector.nameSelector).type(data.routeName1);
  //   cy.contains('Search').click();
  //   cy.contains(data.routeName1).siblings().contains('More').click();
  //   cy.contains('Debug').click();
  //   cy.get(selector.debugDraw).should('be.visible');
  //   cy.get(selector.debugMethod).click();
  //   cy.get('[title=POST]').click();
  //   cy.get(selector.debugProtocol).click();
  //   cy.contains('http://').click();
  //
  //   cy.visit('/');
  //   cy.contains('Route').click();
  //   cy.get(selector.nameSelector).type(data.routeName2);
  //   cy.contains('Search').click();
  //   cy.contains(data.routeName2).siblings().contains('More').click();
  //   cy.contains('Debug').click();
  //   cy.get(selector.debugDraw).should('be.visible');
  //   cy.get(selector.debugMethod).should('contain','method')
  //   cy.get(selector.debugProtocol).should('contain','protocol')
  // })
  // it('should autocomplete header', function () {
  //   cy.visit('/');
  //   cy.contains(menuLocaleUS['menu.routes']).click();
  //   const currentToken = localStorage.getItem('token');
  //
  //   // show online debug draw
  //   cy.get(selector.refresh).click();
  //   cy.contains('Advanced').click();
  //   cy.contains(routeLocaleUS['page.route.onlineDebug']).click();
  //   cy.get(selector.debugDraw).should('be.visible');
  //   cy.get(selector.headerTab).should('be.visible').click();
  //
  //   // show autocomplete
  //   cy.get(selector.headerDataKey0).click({ force: true });
  //   cy.get('.ant-select-item-option-content').contains('Accept').click();
  //   cy.get('.anticon-minus-circle').click();
  //
  //   // autocomplete should ingore case
  //   cy.get(selector.headerDataKey0).type('auth').click({ force: true });
  //   cy.get('.ant-select-item-option-content').contains('Authorization').click();
  //   cy.get(selector.headerDataValue0).type(currentToken);
  // });
  //
  // it('should debug POST request with file successfully', function () {
  //   cy.visit('/');
  //   cy.contains(menuLocaleUS['menu.routes']).click();
  //   const currentToken = localStorage.getItem('token');
  //
  //   // show online debug draw
  //   cy.get(selector.refresh).click();
  //   cy.contains('Advanced').click();
  //   cy.contains(routeLocaleUS['page.route.onlineDebug']).click();
  //   cy.get(selector.debugDraw).should('be.visible');
  //
  //   // change request method POST
  //   cy.get(selector.debugMethod).click();
  //   cy.get('[title=POST]').click();
  //
  //   // change request protocol http
  //   cy.get(selector.debugProtocol).click();
  //   cy.contains('http://').click();
  //   // set debug uri
  //   cy.get(selector.debugUri).type(data.postFileUrl);
  //   // set request body
  //   cy.get(selector.bodyTab).should('be.visible').click();
  //
  //   cy.contains('form-data').should('be.visible').click();
  //   cy.get(selector.debugFormDataKey0).type('file');
  //
  //   // check change type
  //   cy.get(selector.debugFormDataType0).click();
  //   cy.contains('Text').click();
  //   // assert: text input dom should be visible
  //   cy.get(selector.debugFormDataValue0).should('be.visible');
  //   cy.get(selector.debugFormDataType0).click();
  //   cy.contains('File').click();
  //   // assert: upload file button should be visible
  //   cy.get(selector.debugFormDataFileButton0).should('be.visible');
  //   // attach file
  //   cy.get(selector.debugFormDataValue0).attachFile(data.uploadFile);
  //
  //   // set header Authorization
  //   cy.get(selector.headerTab).should('be.visible').click();
  //   cy.get(selector.headerDataKey0).type(data.headerAuthorizationKey);
  //   cy.get(selector.headerDataValue0).type(currentToken);
  //
  //   cy.contains(routeLocaleUS['page.route.button.send']).click();
  //
  //   cy.wait('@DebugAPI');
  //   // assert: send request return
  //   cy.get(selector.codeMirrorResp).contains('data').should('be.visible');
  //   cy.get(selector.codeMirrorResp).contains('routes').should('be.visible');
  //
  //   // close debug drawer
  //   cy.get(selector.drawerClose).click();
  // });

  // it('should debug GET request with query parammeters and jwt auth successfully', function () {
  //   cy.visit('/');
  //   cy.contains(menuLocaleUS['menu.routes']).click();
  //   const currentToken = localStorage.getItem('token');
  //
  //   // show online debug draw
  //   cy.get(selector.refresh).click();
  //   cy.contains('Advanced').click();
  //   cy.contains(routeLocaleUS['page.route.onlineDebug']).click();
  //   cy.get(selector.debugDraw).should('be.visible');
  //   // set debug uri
  //   cy.get(selector.debugUri).type(data.routeOptUrl);
  //   cy.get(selector.bodyTab).should('not.exist');
  //   // set query param
  //   cy.get(selector.queryDataKey0).type(data.queryKey0);
  //   cy.get(selector.queryDataValue0).type(data.routeName);
  //   // set Authentication
  //   cy.get(selector.authTab).should('be.visible').click();
  //   cy.get(selector.jwtAuth).click();
  //   cy.get(selector.jwtTokenInput).should('be.visible').type(currentToken);
  //
  //   cy.contains(routeLocaleUS['page.route.button.send']).click();
  //   cy.wait('@DebugAPI');
  //   cy.get(selector.codeMirrorResp).within(() => {
  //     cy.get('.cm-property').should(($property) => {
  //       // eslint-disable-next-line array-callback-return
  //       $property.map((i, el) => {
  //         if (Cypress.$(el).text() === '"name"') {
  //           const findRouteName = Cypress.$(el).next().text();
  //           expect(findRouteName).to.equal('"hello"');
  //         }
  //         if (Cypress.$(el).text() === '"total_size"') {
  //           const findTotalNumber = Cypress.$(el).next().text();
  //           expect(findTotalNumber).to.equal('1');
  //         }
  //       });
  //     });
  //   });
  // });
  //
  // it('should debug POST request with raw json successfully', function () {
  //   cy.visit('/');
  //   cy.contains(menuLocaleUS['menu.routes']).click();
  //   const currentToken = localStorage.getItem('token');
  //
  //   // show online debug draw
  //   cy.get(selector.refresh).click();
  //   cy.contains('Advanced').click();
  //   cy.contains(routeLocaleUS['page.route.onlineDebug']).click();
  //   cy.get(selector.debugDraw).should('be.visible');
  //
  //   // change request method POST
  //   cy.get(selector.debugMethod).click();
  //   cy.get('[title=POST]').click();
  //
  //   // change request protocol http
  //   cy.get(selector.debugProtocol).click();
  //   cy.contains('http://').click();
  //   // set debug uri
  //   cy.get(selector.debugUri).type(data.routeOptUrl);
  //   // set Authentication
  //   cy.get(selector.authTab).should('be.visible').click();
  //   cy.get(selector.jwtAuth).click();
  //   cy.get(selector.jwtTokenInput).should('be.visible').type(currentToken);
  //
  //   cy.get(selector.bodyTab).should('be.visible').click();
  //
  //   cy.contains('raw input').should('be.visible').click();
  //
  //   cy.window().then(({ codeMirrorBody }) => {
  //     if (codeMirrorBody) {
  //       // eslint-disable-next-line @typescript-eslint/no-invalid-this
  //       codeMirrorBody.setValue(JSON.stringify(this.routeData.debugPostJson));
  //     }
  //     cy.contains(routeLocaleUS['page.route.button.send']).click();
  //   });
  //   cy.wait('@DebugAPI');
  //   cy.get(selector.codeMirrorResp).contains('code').should('be.visible');
  //   cy.get(selector.codeMirrorResp).within(() => {
  //     cy.get('.cm-property').should(($property) => {
  //       // eslint-disable-next-line array-callback-return
  //       $property.map((i, el) => {
  //         if (Cypress.$(el).text() === '"name"') {
  //           const findRouteName = Cypress.$(el).next().text();
  //           // eslint-disable-next-line @typescript-eslint/no-invalid-this
  //           expect(findRouteName).to.equal(`"${this.routeData.debugPostJson.name}"`);
  //         }
  //         if (Cypress.$(el).text() === '"data"') {
  //           // eslint-disable-next-line @typescript-eslint/no-shadow
  //           const data = Cypress.$(el).next().text();
  //           expect(data).to.not.equal('null');
  //         }
  //       });
  //     });
  //   });
  // });
  //
  // it('should delete routes create for test cases successfully', function () {
  //   cy.visit('/');
  //   cy.contains(menuLocaleUS['menu.routes']).click();
  //
  //   // eslint-disable-next-line @typescript-eslint/no-invalid-this
  //   const testRouteNames = [data.routeName, this.routeData.debugPostJson.name];
  //   // eslint-disable-next-line @typescript-eslint/no-for-in-array,guard-for-in,no-restricted-syntax
  //   for (const routeName in testRouteNames) {
  //     cy.contains(`${testRouteNames[routeName]}`).siblings().contains('More').click();
  //     cy.contains('Delete').click({ force: true });
  //     cy.get(selector.deleteAlert)
  //       .should('be.visible')
  //       // eslint-disable-next-line @typescript-eslint/no-loop-func
  //       .within(() => {
  //         cy.contains('OK').click();
  //       });
  //     cy.get(selector.notification).should('contain', data.deleteRouteSuccess);
  //     cy.get(selector.notificationCloseIcon).click({ multiple: true });
  //     cy.reload();
  //   }
  // });
});
