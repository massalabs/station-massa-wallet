import { mockServer } from '../../../../src/mirage';
import { compareSnapshot } from '../../../compareSnapshot';

describe('E2E | Acceptance | Account', () => {
  let server: any;
  const baseUrl = Cypress.config().baseUrl;

  beforeEach(() => {
    server = mockServer('test');
  });

  afterEach(() => {
    server.shutdown();
  });

  const nickname = 'Mario';

  describe('/account-select', () => {
    beforeEach(() => {
      const accounts = server.createList('account', 2);
      accounts.push(server.create('account', { nickname }));
    });
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
      cy.visit('/');
      cy.url().should('eq', `${baseUrl}/index`);

      cy.get('[data-testid="accounts-list"]')
        .find('[data-testid="selector"]')
        .should('have.length', 3);
      cy.get('[data-testid="button"]').should('exist');

      compareSnapshot(cy, 'accounts-list');
    });

    it('should land in create new account when click add account button', () => {
      cy.visit('/');
      cy.get('[data-testid="button"]').click();

      cy.url().should('eq', `${baseUrl}/account-create`);
    });

    it('should land in the account home when click on account selector', () => {
      cy.visit('/');
      cy.get('[data-testid="account-2"]').click();

      cy.url().should('eq', `${baseUrl}/${nickname}/home`);
    });
  });
});
