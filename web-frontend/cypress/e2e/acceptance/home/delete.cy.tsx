import { mockServer } from '../../../../src/mirage';
import { compareSnapshot } from '../../../compareSnapshot';
import { AccountObject } from '@/models/AccountModel';

describe('E2E | Acceptance | Home', () => {
  let server: any;
  const baseUrl = Cypress.config().baseUrl;

  beforeEach(() => {
    server = mockServer('test');
  });

  afterEach(() => {
    server.shutdown();
  });

  describe('delete', () => {
    let mockedAccounts: AccountObject[];

    function mockAccounts() {
      const mockedAccounts = server.createList('account', 1);
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
    function selectAccount(index: number) {
      cy.visit('/');
      cy.get(`[data-testid="account-${index}"]`).click();
    }

    function navigateToSettings() {
      cy.get('[data-testid="side-menu"]').click();
      cy.get('[data-testid="side-menu-settings-icon"]').click();
    }

    it('should direct me to /nickname/home', () => {
      const account = mockedAccounts.at(1);
      if (account === undefined) throw new Error('Account not found');

      selectAccount(1);

      cy.url().should('eq', `${baseUrl}/${account.nickname}/home`);

      cy.get('[data-testid="side-menu"]').should('exist').and('be.visible');
      compareSnapshot(cy, 'delete-home');
    });

    it('should access delete page', () => {
      const account = mockedAccounts.at(1);
      if (account === undefined) throw new Error('Account not found');

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
