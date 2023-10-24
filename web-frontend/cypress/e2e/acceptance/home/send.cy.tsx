/* eslint-disable @typescript-eslint/ban-ts-comment */
// @ts-nocheck
import { Server } from 'miragejs';

import { mockServer } from '../../../../src/mirage';
import {
  toMASS,
  formatStandard,
  maskAddress,
} from '../../../../src/utils/massaFormat';
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

    // util functions
    function navigateToHome() {
      cy.visit('/');

      cy.get('[data-testid="account-2"]').click();
    }

    function navigateToTransfercoins() {
      cy.visit('/');

      cy.get('[data-testid="account-2"]').click();
      cy.get('[data-testid="side-menu"]').click();
      cy.get('[data-testid="side-menu-sendreceive-icon"]').click();
    }

    function getBalance(account) {
      // same logic we use in SendForm.tsx where balance is rendered
      const balance = Number(account?.candidateBalance || 0);
      const formattedBalance = formatStandard(toMASS(balance));
      return formattedBalance;
    }

    // mimics currency field formating
    function customFormatNumber(number) {
      // Convert the number to a string
      let numberString = number.toString();

      // Split the number into integer and decimal parts
      let [integerPart, decimalPart] = numberString.split('.');

      // Add commas for thousands separation to the integer part
      integerPart = integerPart.replace(/\B(?=(\d{3})+(?!\d))/g, ',');

      // If there is a decimal part, combine integer and decimal parts with a dot
      if (decimalPart !== undefined) {
        return `${integerPart}.${decimalPart} MAS`;
      } else {
        return `${integerPart} MAS`;
      }
    }

    it('should have loading state', () => {
      const account = mockedAccounts.at(2);
      navigateToHome();

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
      navigateToHome();
      cy.get('[data-testid="send-button"]').click();
      cy.url().should('eq', `${baseUrl}/${account.nickname}/transfer-coins`);

      cy.get('[data-testid="send-coins"]').should('be.visible');
      cy.get('[data-testid="receive-coins"]').should('not.be.visible');
      cy.get('[data-testid="send-coins"]').should('exist');

      compareSnapshot(cy, 'send-page');
    });

    it('should navigate to /transfer-coins when accessing it by the side-menu', () => {
      const account = mockedAccounts.at(2);

      navigateToHome();
      cy.get('[data-testid="side-menu"]').click();
      cy.get('[data-testid="side-menu-sendreceive-icon"]').click();
      cy.url().should('eq', `${baseUrl}/${account.nickname}/transfer-coins`);
      compareSnapshot(cy, 'send-page');
    });

    it('should render balance and amount should equal account balance', () => {
      const account = mockedAccounts.at(2);

      navigateToTransfercoins();
      // cy.get('[data-testid="balance-amount"]').contains(getBalance(account));
      cy.get('[data-testid="balance-amount"]').should(
        'contain',
        getBalance(account),
      );
      compareSnapshot(cy, 'send-page');
    });

    it('should send coins to address', () => {
      const account = mockedAccounts.at(2);
      const recipientAccount = mockedAccounts.at(1);

      const amount = 1000.12311132;

      const standardFees = '1000';

      navigateToTransfercoins();
      cy.get('[data-testid="currency-field"')
        .should('exist')
        .type(amount)
        .should('have.value', customFormatNumber(amount));

      cy.get('[data-testid="input-field"').type(recipientAccount.address);

      cy.get('[data-testid="button"]').contains('Send').click();
      cy.url().should('eq', `${baseUrl}/${account.nickname}/transfer-coins`);
      cy.get('[data-testid="send-confirmation"]').should('exist');

      cy.get('[data-testid="send-confirmation-info"]').should(
        'contain',
        `Amount (${customFormatNumber(amount)}) + fees (${standardFees} nMAS)`,
      );

      cy.get('[data-testid="balance-amount"]').contains(formatStandard(amount));

      cy.get('[data-testid="send-confirmation-recipient"]').contains(
        maskAddress(recipientAccount.address),
      );

      compareSnapshot(cy, 'send-confirmation-page');

      cy.get('[data-testid="button"]')
        .contains('Confirm and sign with password')
        .click();
      cy.url().should('eq', `${baseUrl}/${account.nickname}/home`);

      compareSnapshot(cy, 'home-page');
    });

    // it('should land on receive page when receive CTA is clicked', () => {
    //   const account = mockedAccounts.at(2);

    //   cy.visit('/');

    //   cy.get('[data-testid="account-2"]').click();
    //   cy.get('[data-testid="receive-button"]').click();
    //   cy.url().should('eq', `${baseUrl}/${account.nickname}/transfer-coins`);

    //   cy.get('[data-testid="send-coins"]').should('not.be.visible');
    //   cy.get('[data-testid="receive-coins"]').should('be.visible');
    // });

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
