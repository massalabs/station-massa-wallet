/* eslint-disable @typescript-eslint/ban-ts-comment */
// @ts-nocheck
import { Server } from 'miragejs';

import { mockServer } from '../../../../src/mirage';
// import { compareSnapshot } from '../../../compareSnapshot';

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
      const mockedAccounts = server.createList('account', 2);
      mockedAccounts.forEach((account) => {
        server.createList('asset', 3, { account });
      });
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
      const account = mockedAccounts.at(2);

      cy.visit('/');

      cy.get('[data-testid="account-2"]').click();

      cy.url().should('eq', `${baseUrl}/${account.nickname}/home`);
    });
  });
});
