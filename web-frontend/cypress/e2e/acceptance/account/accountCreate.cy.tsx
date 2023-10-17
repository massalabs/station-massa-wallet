import { useParams } from 'react-router-dom';

import { compareSnapshot } from '../../../compareSnapshot';

describe('E2E | Acceptance | create', function () {
  describe('create account flow', () => {
    it('it should render select account page', () => {
      cy.visit(`/`);
      cy.get('[data-testid="account-select-page"]').should('exist');

      cy.get('[data-testid="account-select-list"]').should('exist');

      cy.get('[data-testid="account-create"]').should('exist');

      cy.get('[data-testid="account-select-list-0"]').should('exist');
    });
  });
  //   it('should display accounts list after loading', () => {
  //     cy.visit(`/`).then(() => {
  //       cy.url().should('include', '/account-select');

  //       cy.intercept(
  //         {
  //           method: 'GET',
  //           url: `${Cypress.config().baseUrl}/accounts`,
  //         },
  //         { fixture: 'accounts.json' },
  //       ).as('getAccounts');

  //       cy.wait('@getAccounts');
  //     });
  //   });
});
