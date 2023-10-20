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

    it('should have loading state', () => {
      const account = mockedAccounts.at(2);

      cy.visit('/');

      cy.get('[data-testid="account-2"]').click();

      cy.url()
        .should('eq', `${baseUrl}/${account.nickname}/home`)
        .then(() => {
          cy.get('[data-testid="loading"]')
            .should('be.visible')
            .then(() => {
              cy.get('[data-testid="loading"]').should('not.exist');
            });
        });

      compareSnapshot(cy, 'wallet-loading');
    });

    it('should land on send page when send CTA is clicked', () => {
      const account = mockedAccounts.at(2);

      cy.visit('/');

      cy.get('[data-testid="account-2"]').click();
      cy.get('[data-testid="send-button"]').click();
      cy.url().should('eq', `${baseUrl}/${account.nickname}/transfer-coins`);

      cy.get('[data-testid="send-coins"]').should('be.visible');
      cy.get('[data-testid="receive-coins"]').should('not.be.visible');
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

    // it('should copy wallet address when clipboard field is clicked', () => {
    //   // we are adding the permission to chrome on the fly
    //   cy.wrap(
    //     Cypress.automation('remote:debugger:protocol', {
    //       command: 'Browser.grantPermissions',
    //       params: {
    //         permissions: ['clipboardReadWrite', 'clipboardSanitizedWrite'],
    //         origin: window.location.origin,
    //       },
    //     }),
    //   );

    //   const account = mockedAccounts.at(2);

    //   cy.visit('/');

    //   cy.get('[data-testid="account-2"]').click();
    //   cy.url().should('eq', `${baseUrl}/Mario/home`);

    //   compareSnapshot(cy, 'wallet-home');

    //   cy.get('[data-testid="clipboard-field"]').click();
    //   cy.assertValueCopiedFromClipboard(account.address);
    // });
  });
});
