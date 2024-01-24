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
    let mockedAccounts;

    function mockAccounts() {
      const mockedAccounts = server.createList('account', 1);
      mockedAccounts.forEach((account) => {
        server.createList('asset', 1, { account });
      });
      const account = {
        nickname: 'Mario',
        address: 'AU1ZMBZeARHYMFfV4uvbyCB85DAUPr2BJXzU1kSYdwCKKrY5crWY',
      };
      mockedAccounts.push(server.create('account', { ...account }));

      return mockedAccounts;
    }

    beforeEach(() => {
      mockedAccounts = mockAccounts();
    });

    // util functions
    function selectAccount(index) {
      const account = mockedAccounts.at(index);
      cy.visit('/');
      cy.get(`[data-testid="account-${index}"]`).click();
    }

    function navigateToSettings() {
      cy.get('[data-testid="side-menu"]').click();
      cy.get('[data-testid="side-menu-settings-icon"]').click();
    }

    it('should direct me to /nickname/home', () => {
      const account = mockedAccounts.at(1);

      selectAccount(1);

      cy.url().should('eq', `${baseUrl}/${account.nickname}/home`);

      cy.get('[data-testid="side-menu"]').should('exist').and('be.visible');
      compareSnapshot(cy, 'delete-home');
    });

    it('should access delete page', () => {
      const account = mockedAccounts.at(1);

      selectAccount(1);

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

      selectAccount(1);
      navigateToSettings();

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
