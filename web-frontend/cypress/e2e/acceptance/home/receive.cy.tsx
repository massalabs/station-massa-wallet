/* eslint-disable @typescript-eslint/ban-ts-comment */
// @ts-nocheck
import { mockServer } from '../../../../src/mirage';
import { AccountObject } from '@/models/AccountModel';

describe('E2E | Acceptance | Account | Receive', () => {
  let server: any;
  const baseUrl = Cypress.config().baseUrl;

  beforeEach(() => {
    server = mockServer('test');
  });

  afterEach(() => {
    server.shutdown();
  });

  describe('send', () => {
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

    it('should land on receive page when receive CTA is clicked', () => {
      const account = mockedAccounts.at(2);

      cy.visit('/');

      cy.get('[data-testid="account-2"]').click();
      cy.get('[data-testid="receive-button"]').click();
      cy.url().should('eq', `${baseUrl}/${account.nickname}/transfer-coins`);

      cy.get('[data-testid="send-coins"]').should('not.be.visible');
      cy.get('[data-testid="receive-coins"]').should('be.visible');
    });
  });
});
