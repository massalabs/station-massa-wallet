import { mockServer } from '../../../../src/mirage';
import { compareSnapshot } from '../../../compareSnapshot';

describe('E2E | Acceptance | Account | Create', () => {
  let server: any;
  const baseUrl = Cypress.config().baseUrl;

  beforeEach(() => {
    server = mockServer('test');
  });

  afterEach(() => {
    server.shutdown();
  });

  describe('/account-create', () => {
    it('should have no accounts', () => {
      cy.visit('/');
      cy.url().should('eq', `${baseUrl}/index`);

      cy.get('[data-testid="accounts-list"]').should('not.exist');
      cy.get('[data-testid="button"]').should('exist');
      compareSnapshot(cy, 'account-none');
    });

    it('should direct to /create-account when click create account button', () => {
      cy.visit('/');
      cy.get('[data-testid="button"]').contains('Create an account').click();

      cy.url().should('eq', `${baseUrl}/account-create-step-one`);
      compareSnapshot(cy, 'account-create-step-one');
    });

    it('should create a new account', () => {
      cy.visit('/account-create-step-one');

      cy.url().should('eq', `${baseUrl}/account-create-step-one`);

      cy.get('[data-testid="input-field"]').type('testAccount');
      cy.get('[data-testid="button"]').contains('Next').click();

      cy.url().should('eq', `${baseUrl}/account-create-step-two`);
      compareSnapshot(cy, 'account-create-step-two');

      cy.get('[data-testid="button"]').contains('Define a password').click();

      cy.url().should('eq', `${baseUrl}/account-create-step-three`);
      compareSnapshot(cy, 'account-create-step-three');

      cy.get('[data-testid="button"]').contains('Skip').click();

      server.create('account');

      cy.waitForRequest(server, '/accounts/testAccount', 'POST').then(() => {
        cy.url().should('eq', `${baseUrl}/testAccount/home`);
        compareSnapshot(cy, 'account-create-home');
      });
    });

    it('should prevent me from creating an account with wrong nickname', () => {
      cy.visit('/account-create-step-one');

      cy.get('[data-testid="input-field"]').type('testA ccount');
      cy.get('[data-testid="button"]').contains('Next').click();

      cy.get('[data-testid="input-field-message"]')
        .should('exist')
        .contains(
          "The account name can't contain any special characters, and 32 characters max.",
        );
      cy.get('[data-testid="button"]').contains('Next').click();
      cy.url().should('eq', `${baseUrl}/account-create-step-one`);

      compareSnapshot(cy, 'wrong-account-name-format');
    });
  });
});
