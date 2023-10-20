/* eslint-disable @typescript-eslint/ban-ts-comment */
// @ts-nocheck
import { Server } from 'miragejs/server';

import { mockServer } from '../../../src/mirage';

describe('E2E | Acceptance | Redirection', () => {
  let server: Server;
  const baseUrl = Cypress.config().baseUrl;

  beforeEach(() => {
    server = mockServer('test');

    const accounts = server.createList('account', 2);
    accounts.forEach((account) => {
      server.createList('asset', 3, { account });
    });
  });

  afterEach(() => {
    server.shutdown();
  });

  it('should land in /account-select when hit /', () => {
    cy.visit('/');

    cy.url().should('eq', `${baseUrl}/index`);
    cy.url().should('eq', `${baseUrl}/account-select`);
  });

  it('should land in /account-select when hit /index', () => {
    cy.visit('/index');

    cy.url().should('eq', `${baseUrl}/index`);
    cy.url().should('eq', `${baseUrl}/account-select`);
  });
});
