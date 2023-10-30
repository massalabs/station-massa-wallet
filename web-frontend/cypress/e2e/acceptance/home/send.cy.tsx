/* eslint-disable @typescript-eslint/ban-ts-comment */
// @ts-nocheck
import { wait } from '@massalabs/massa-web3/dist/utils/time';
import { Server } from 'miragejs';

import { mockServer } from '../../../../src/mirage';
import {
  toMASS,
  formatStandard,
  maskAddress,
} from '../../../../src/utils/massaFormat';
import { handlePercent } from '../../../../src/utils/math';

describe('E2E | Acceptance | Home | Send', () => {
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
        nickname: 'MarioX',
        candidateBalance: '1000000000000',
        balance: '10000',
        address: 'AUHdadXyJZUeINwiUVMtXZXJRTFXtYdihRWitUcAJSBwAHgcKAjtxx',
      };
      mockedAccounts.push(server.create('account', { ...account }));
      return mockedAccounts;
    }

    beforeEach(() => {
      mockedAccounts = mockAccounts();
    });

    function navigateToHome() {
      cy.visit('/');
      cy.get('[data-testid="account-2"]').click();
    }

    function navigateToTransfercoinsOfAccountIndex(index) {
      cy.visit('/');

      cy.get(`[data-testid="account-${index}"]`).click();
      cy.get('[data-testid="side-menu"]').click();
      cy.get('[data-testid="side-menu-sendreceive-icon"]').click();
    }

    function setAccountBalance(account) {
      const balance = Number(account.candidateBalance);
      const formattedBalance = formatStandard(toMASS(balance));
      return formattedBalance;
    }

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
    });

    it('should land on send page when send CTA is clicked', () => {
      const account = mockedAccounts.at(2);
      navigateToHome();
      cy.get('[data-testid="send-button"]').click();
      cy.url().should('eq', `${baseUrl}/${account.nickname}/transfer-coins`);

      cy.get('[data-testid="send-coins"]').should('be.visible');
      cy.get('[data-testid="receive-coins"]').should('not.be.visible');
    });

    it('should navigate to /transfer-coins when accessing it by the side-menu', () => {
      const account = mockedAccounts.at(2);

      navigateToHome();
      cy.get('[data-testid="side-menu"]').click();
      cy.get('[data-testid="side-menu-sendreceive-icon"]').click();
      cy.url().should('eq', `${baseUrl}/${account.nickname}/transfer-coins`);
    });

    it('should render balance and amount should equal account balance', () => {
      const account = mockedAccounts.at(2);

      navigateToTransfercoinsOfAccountIndex(2);

      cy.get('[data-testid="balance').should('exist');

      cy.get('[data-testid="balance-amount"]').contains(
        setAccountBalance(account),
      );
    });

    it('should send coins to address', () => {
      const account = mockedAccounts.at(2);
      const recipientAccount = mockedAccounts.at(1);
      const amount = 550.1234;
      const standardFees = '1000';

      navigateToTransfercoinsOfAccountIndex(2);
      cy.get('[data-testid="currency-field"')
        .type(amount)
        .should('have.value', '550.1234 MAS');
      cy.get('[data-testid="input-field"').type(recipientAccount.address);

      cy.get('[data-testid="button"]').contains('Send').click();
      cy.url().should('eq', `${baseUrl}/${account.nickname}/transfer-coins`);
      cy.get('[data-testid="send-confirmation"]').should('be.visible');

      cy.get('[data-testid="send-confirmation-info"]').should(
        'contain',
        `Amount (550.1234 MAS) + fees (${standardFees} nMAS)`,
      );

      cy.get('[data-testid="balance-amount"]').contains(formatStandard(amount));

      cy.get('[data-testid="send-confirmation-recipient"]').contains(
        maskAddress(recipientAccount.address),
      );

      cy.get('[data-testid="button"]')
        .contains('Confirm and sign with password')
        .click();
      cy.url().should('eq', `${baseUrl}/${account.nickname}/home`);
    });

    it('should process any % as input', async () => {
      server.trackRequest = true;
      const account = mockedAccounts.at(2);

      cy.visit('/');

      cy.get('[data-testid="account-2"]')
        .click()
        .then(() => {
          cy.url().should('eq', `${baseUrl}/${account.nickname}/home`);
          cy.waitForRequest(server, '/accounts/MarioX', 'GET').then(() => {
            cy.get('[data-testid="send-button"]')
              .click()
              .then(() => {
                cy.waitForRequest(server, '/accounts/MarioX', 'GET');
                cy.url().should(
                  'eq',
                  `${baseUrl}/${account.nickname}/transfer-coins`,
                );
                cy.get(`[data-testid="send-percent-25"]`).click();
                cy.get('[data-testid="currency-field"').should(
                  'have.value',
                  '250 MAS',
                );

                cy.get(`[data-testid="send-percent-50"]`).click();
                cy.get('[data-testid="currency-field"').should(
                  'have.value',
                  '500 MAS',
                );

                cy.get(`[data-testid="send-percent-75"]`).click();
                cy.get('[data-testid="currency-field"').should(
                  'have.value',
                  '750 MAS',
                );

                cy.get(`[data-testid="send-percent-100"]`).click();
                cy.get('[data-testid="currency-field"').should(
                  'have.value',
                  '999.999999 MAS',
                );
              });
          });
        });
      server.trackRequest = false;
    });

    it('should transfer to accounts', () => {
      const randomIndex = Math.floor(Math.random() * mockedAccounts.length);
      const selectedAccount = mockedAccounts.at(randomIndex);

      navigateToTransfercoinsOfAccountIndex(2);

      cy.get('[data-testid="transfer-between-accounts"]').click();

      cy.get('[data-testid="popup-modal-content"').should('be.visible');

      for (let i = 0; i < mockedAccounts.length; i++) {
        cy.get(`[data-testid="selector-account-${i}"]`).should('be.visible');
      }

      cy.get(`[data-testid="selector-account-${randomIndex}"]`).click();
      cy.get(`[data-testid="input-field"]`).should(
        'have.value',
        selectedAccount.address,
      );
    });

    it('should refuse wrong currency input', () => {
      const account = mockedAccounts.at(2);
      const recipientAccount = mockedAccounts.at(1);

      const amount = 1000.12311132;
      const invalidAmount = 'things';
      const tooMuch = Number(account.candidateBalance) + 1000;
      const tooLow = 0;
      const notEnoughForFees = Number(account.candidateBalance);

      const standardFees = '1000';
      navigateToTransfercoinsOfAccountIndex(2);

      cy.get('[data-testid="currency-field"')
        .should('exist')
        .type(invalidAmount);

      cy.get('[data-testid="input-field"').type(recipientAccount.address);
      cy.get('[data-testid="button"]').contains('Send').click();

      cy.get('[data-testid="input-field-message"]').should(
        'contain.text',
        'Invalid amount',
      );

      cy.get('[data-testid="currency-field"').clear().type(tooMuch);
      cy.get('[data-testid="button"]').contains('Send').click();

      cy.get('[data-testid="input-field-message"]').should(
        'contain.text',
        'Insufficient funds',
      );

      cy.get('[data-testid="currency-field"').clear().type(tooLow);
      cy.get('[data-testid="button"]').contains('Send').click();

      cy.get('[data-testid="input-field-message"]').should(
        'contain.text',
        'Amount must be greater than zero',
      );

      cy.get('[data-testid="currency-field"').clear().type(notEnoughForFees);
      cy.get('[data-testid="button"]').contains('Send').click();

      cy.get('[data-testid="input-field-message"]').should(
        'contain.text',
        'Insufficient funds',
      );
    });

    it('should refuse wrong addrress input', () => {
      const account = mockedAccounts.at(2);
      const wrongAddress = 'wrong address';
      const amount = 42;

      navigateToTransfercoinsOfAccountIndex(2);
      cy.get('[data-testid="currency-field"')
        .type(amount)
        .should('have.value', '42 MAS');

      cy.get('[data-testid="input-field"').type(wrongAddress);
      cy.get('[data-testid="button"]').contains('Send').click();

      cy.get('[data-testid="input-field-message"]').should(
        'contain.text',
        'Recipient address is not a valid Massa address',
      );

      cy.get('[data-testid="input-field"').clear();
      cy.get('[data-testid="button"]').contains('Send').click();

      cy.get('[data-testid="input-field-message"]').should(
        'contain.text',
        'Recipient is required',
      );
    });
  });
});
