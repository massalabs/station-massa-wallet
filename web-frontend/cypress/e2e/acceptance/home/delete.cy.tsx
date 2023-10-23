/* eslint-disable @typescript-eslint/ban-ts-comment */
// @ts-nocheck
import { Server } from 'miragejs';

import { mockServer } from '../../../../src/mirage';
import { compareSnapshot } from '../../../compareSnapshot';

describe('E2E | Acceptance | Home', () => {
  let server: Server;
  const baseUrl = Cypress.config().baseUrl;

  beforeEach(() => {
    server = mockServer('test');
  });

  afterEach(() => {
    server.shutdown();
  });

  describe('delete', () => {
    // accounts used for test
    let mockedAccounts;

    function mockAccounts() {
      const mockedAccounts = server.createList('account', 1);
      mockedAccounts.forEach((account) => {
        server.createList('asset', 1, { account });
      });
      // reference account
      const account = {
        nickname: 'Mario',
        address: 'AUHdadXyJZUeINwiUVMtXZXJRTFXtYdihRWitUcAJSBwAHgcKAjtxx',
      };
      mockedAccounts.push(server.create('account', { ...account }));

      return mockedAccounts;
    }

    beforeEach(() => {
      mockedAccounts = mockAccounts();
    });

    it('should direct me to /nickname/home', () => {
      // position of reference account in mockedAccounts
      const account = mockedAccounts.at(1);

      cy.visit('/');

      cy.get('[data-testid="account-1"]').click();

      cy.url().should('eq', `${baseUrl}/${account.nickname}/home`);

      cy.get('[data-testid="side-menu"]').should('exist').and('be.visible');
      compareSnapshot(cy, 'delete-home');
    });

    it('should access delete page', () => {
      const account = mockedAccounts.at(1);

      cy.visit('/');

      cy.get('[data-testid="account-1"]').click();

      cy.get('[data-testid="side-menu"]').click();

      cy.get('[data-testid="side-menu-settings-icon"]')
        .should('exist')
        .and('be.visible')
        .click();

      cy.url().should('eq', `${baseUrl}/${account.nickname}/settings`);
      compareSnapshot(cy, 'delete-settings-page');
    });

    it('should delete account', () => {
      const account = mockedAccounts.at(1);

      cy.visit('/');

      cy.get('[data-testid="account-1"]').click();

      cy.get('[data-testid="side-menu"]').click();

      cy.get('[data-testid="side-menu-settings-icon"]').click();

      cy.get('[data-testid="button"]')
        .contains('Delete account')
        .should('exist')
        .and('be.visible')
        .click();

      cy.url().should('eq', `${baseUrl}/account-select`);

      cy.get('[data-testid="account-1"]').should('not.exist');
      compareSnapshot(cy, 'delete-account-select');
    });
  });
});
