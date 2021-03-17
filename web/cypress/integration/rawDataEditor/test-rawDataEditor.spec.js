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

context('Test RawDataEditor', () => {
  const timeout = 1000;

  beforeEach(() => {
    cy.login();

    cy.fixture('rawDataEditor-dataset.json').as('dataset');
    cy.fixture('selector.json').as('domSelectors');
  });

  it('should create and update with rawDataEditor', function () {
    const menuList = ['Route', 'Service', 'Upstream', 'Consumer'];
    const dateset = this.dataset;
    const domSelectors = this.domSelectors;
    console.log('dateset: ', dateset);
    menuList.forEach(function (item) {
      cy.visit('/');
      cy.contains(item).click();
      cy.contains('Create with Editor').click();
      const data = dateset[item];
      // create with editor
      cy.window().then(({ codemirror }) => {
        if (codemirror) {
          codemirror.setValue(JSON.stringify(data));
        }
        cy.get(domSelectors.drawer).should('exist');
        cy.get(domSelectors.drawer, { timeout }).within(() => {
          cy.contains('Submit').click({
            force: true,
          });
          cy.get(domSelectors.drawer).should('not.exist');
        });
      });

      cy.reload();

      // update with editor
      cy.contains(item === 'Consumer' ? data.username : data.name)
        .siblings()
        .contains('View')
        .click();

      cy.window().then(({ codemirror }) => {
        if (codemirror) {
          if (item === 'Consumer') {
            codemirror.setValue(JSON.stringify({ ...data, desc: 'newDesc' }));
          } else {
            codemirror.setValue(JSON.stringify({ ...data, name: 'newName' }));
          }
        }
        cy.get(domSelectors.drawer).should('exist');
        cy.get(domSelectors.drawer, { timeout }).within(() => {
          cy.contains('Submit').click({
            force: true,
          });
          cy.get(domSelectors.drawer).should('not.exist');
        });
      });

      cy.reload();
      cy.get('.ant-table-tbody').should('contain', item === 'Consumer' ? 'newDesc' : 'newName');
    });
  });
});
