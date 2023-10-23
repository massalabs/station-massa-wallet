import { mockServer } from '../../../../src/mirage';
import { compareSnapshot } from '../../../compareSnapshot';
import { AccountObject } from '@/models/AccountModel';

describe('E2E | Acceptance | Account', () => {
  let server: any;
  const baseUrl = Cypress.config().baseUrl;

  beforeEach(() => {
    server = mockServer('test');
  });

  afterEach(() => {
    server.shutdown();
  });

  describe('/account-select', () => {
    it('should have loading state', () => {
      cy.visit('/');

      cy.url()
        .should('eq', `${baseUrl}/index`)
        .then(() => {
          cy.get('[data-testid="loading"]')
            .should('be.visible')
            .then(() => {
              cy.get('[data-testid="loading"]').should('not.exist');
            });
        });

      compareSnapshot(cy, 'account-select-loading');
    });

    it('should have accounts loaded', () => {
      const accounts = server.createList('account', 2);
      accounts.forEach((account: AccountObject) => {
        server.createList('asset', 3, { account });
      });

      cy.visit('/');
      cy.url().should('eq', `${baseUrl}/index`);

      cy.get('[data-testid="accounts-list"]')
        .find('[data-testid="selector"]')
        .should('have.length', 2);
      cy.get('[data-testid="button"]').should('exist');

      compareSnapshot(cy, 'accounts-list');
    });

    it('should land in create new account when click add account button', () => {
      const accounts = server.createList('account', 2);
      accounts.forEach((account: AccountObject) => {
        server.createList('asset', 3, { account });
      });

      cy.visit('/');
      cy.get('[data-testid="button"]').click();

      cy.url().should('eq', `${baseUrl}/account-create`);
    });

    it('should land in the account home when click on account selector', () => {
      const accounts = server.createList('account', 2);
      const nickname = 'Mario';
      accounts.push(server.create('account', { nickname }));
      accounts.forEach((account: AccountObject) => {
        server.createList('asset', 3, { account });
      });

      cy.visit('/');
      cy.get('[data-testid="account-2"]').click();

      cy.url().should('eq', `${baseUrl}/${nickname}/home`);
    });
  });
});
