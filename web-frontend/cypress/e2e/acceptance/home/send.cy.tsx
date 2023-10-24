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

    function navigateToTransfercoins(index) {
      cy.visit('/');

      cy.get(`[data-testid="account-${index}"]`).click();
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
    // custom fn because currency field input component uses a specific component node modules
    function customFormatNumber(number) {
      let numberString = number.toString();
      let [integerPart, decimalPart] = numberString.split('.');
      integerPart = integerPart.replace(/\B(?=(\d{3})+(?!\d))/g, ',');
      if (decimalPart !== undefined) {
        return `${integerPart}.${decimalPart} MAS`;
      } else {
        return `${integerPart} MAS`;
      }
    }

    // mimics handlepercent fn in sendForm.tsx
    // TODO: refactor handlePercent to be a util function
    function handlePercent(amount = 0, percent: number, fees, balance) {
      let newAmount = amount * percent;
      const feesAsNumber = Number(fees);

      if (newAmount > balance - feesAsNumber)
        newAmount = Math.max(balance - feesAsNumber, 0);

      return toMASS(newAmount);
    }

    // standardize percent testing
    function performPercentAction(percentValue, account, fees) {
      const balance = Number(account?.candidateBalance);
      cy.get(`[data-testid="send-percent-${percentValue * 100}"]`)
        .should('exist')
        .trigger('mouseover')
        .click();

      cy.get('[data-testid="currency-field"')
        .should(
          'have.value',
          customFormatNumber(
            handlePercent(
              balance,
              percentValue,
              fees,
              account?.candidateBalance,
            ),
          ),
        )
        .clear();
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

      navigateToTransfercoins(2);

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

      navigateToTransfercoins(2);
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

    it('should process any % as input', () => {
      const account = mockedAccounts.at(1);
      const recipientAccount = mockedAccounts.at(2);
      const standardFees = '1000';

      navigateToTransfercoins(1);

      // we have to wait or else balance is not initialized properly
      cy.wait(1000);

      performPercentAction(0.25, account, standardFees);
      performPercentAction(0.5, account, standardFees);
      performPercentAction(0.75, account, standardFees);
      performPercentAction(1, account, standardFees);
    });

    it('should transfer to accounts', () => {
      const randomIndex = Math.floor(Math.random() * mockedAccounts.length);
      const selectedAccount = mockedAccounts.at(randomIndex);

      console.log(mockedAccounts);

      navigateToTransfercoins(2);

      cy.get('[data-testid="transfer-between-accounts"]')
        .should('exist')
        .click();

      cy.get('[data-testid="popup-modal-content"').should('exist');

      cy.get('[data-testid="selector-account-list"]').should('exist');

      for (let i = 0; i < mockedAccounts.length; i++) {
        cy.get(`[data-testid="selector-account-${i}"]`).should('exist');
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
      const tooMuch = Number(account?.candidateBalance) + 1000;
      const tooLow = 0;
      const notEnoughForFees = Number(account?.candidateBalance);

      const standardFees = '1000';
      navigateToTransfercoins(2);

      cy.get('[data-testid="currency-field"')
        .should('exist')
        .type(invalidAmount);

      cy.get('[data-testid="input-field"').type(recipientAccount.address);

      cy.get('[data-testid="button"]').contains('Send').click();

      cy.get('[data-testid="input-field-message"]').should(
        'contain.text',
        'Invalid amount',
      );

      cy.get('[data-testid="currency-field"')
        .clear()
        .should('exist')
        .type(tooMuch);

      cy.get('[data-testid="button"]').contains('Send').click();

      cy.get('[data-testid="input-field-message"]').should(
        'contain.text',
        'Insufficient funds',
      );

      cy.get('[data-testid="currency-field"')
        .clear()
        .should('exist')
        .type(tooLow);

      cy.get('[data-testid="button"]').contains('Send').click();

      cy.get('[data-testid="input-field-message"]').should(
        'contain.text',
        'Amount must be greater than zero',
      );

      // testing amount + fees to hugh seems tricky
      // amount + fees to high also triggers insufficient funds
      // there mayeb we can refacto it ?
      // See sendForm.tsx line 84

      cy.get('[data-testid="currency-field"')
        .clear()
        .should('exist')
        .type(notEnoughForFees);

      cy.get('[data-testid="button"]').contains('Send').click();

      cy.get('[data-testid="input-field-message"]').should(
        'contain.text',
        'Insufficient funds',
      );
    });

    it('should refuse wrong addrress input', () => {
      const account = mockedAccounts.at(2);
      const wrongAddress = 'wrong address';

      const amount = 1000.12311132;

      navigateToTransfercoins(2);
      cy.get('[data-testid="currency-field"')
        .should('exist')
        .type(amount)
        .should('have.value', customFormatNumber(amount));

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
