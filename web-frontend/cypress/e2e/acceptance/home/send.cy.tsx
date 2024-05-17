import { mockServer } from '../../../../src/mirage';
import { AccountObject } from '@/models/AccountModel';

describe('E2E | Acceptance | Home | Send', () => {
  let server: any;
  const baseUrl = Cypress.config().baseUrl;

  beforeEach(() => {
    server = mockServer('test');
  });

  afterEach(() => {
    server.shutdown();
  });

  describe('send', () => {
    let mockedAccounts: AccountObject[];

    function mockAccounts() {
      const mockedAccounts = server.createList('account', 2);
      mockedAccounts.forEach((account: AccountObject) => {
        server.createList('asset', 3, { account });
      });

      const account = server.create('account', {
        nickname: 'MarioX',
        candidateBalance: '1000000000000',
        balance: '10000',
        address: 'AU1ZMBZeARHYMFfV4uvbyCB85DAUPr2BJXzU1kSYdwCKKrY5crWY',
      });
      const assets = server.createList('asset', 3, { account });
      account.assets = assets;

      mockedAccounts.push(account);
      return mockedAccounts;
    }

    function navigateToSendForm() {
      server.trackRequest = true;
      const account = mockedAccounts.at(2);
      if (account === undefined) throw new Error('Account not found');

      cy.visit('/');

      cy.get('[data-testid="account-2"]').click();

      cy.url().should('eq', `${baseUrl}/${account.nickname}/home`);

      cy.get('[data-testid="send-button"]').click();
      cy.url().should('eq', `${baseUrl}/${account.nickname}/transfer-coins`);

      cy.get('[data-testid="spinner"]').should('not.exist');
    }

    beforeEach(() => {
      mockedAccounts = mockAccounts();
    });

    function navigateToHome() {
      cy.visit('/');
      cy.get('[data-testid="account-2"]').should('exist').click();
    }

    function navigateToTransferCoinsOfAccountIndex(index: string) {
      cy.visit('/');

      cy.get(`[data-testid="account-${index}"]`).should('exist').click();
      cy.get('[data-testid="side-menu"]').should('exist').click();
      cy.get('[data-testid="side-menu-sendreceive-icon"]')
        .should('exist')
        .click();
    }

    it('should have loading state', () => {
      const account = mockedAccounts.at(2);
      if (account === undefined) throw new Error('Account not found');

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
      if (account === undefined) throw new Error('Account not found');

      navigateToHome();
      cy.get('[data-testid="send-button"]').click();
      cy.url().should('eq', `${baseUrl}/${account.nickname}/transfer-coins`);

      cy.get('[data-testid="send-coins"]').should('be.visible');
      cy.get('[data-testid="receive-coins"]').should('not.be.visible');
    });

    it('should navigate to /transfer-coins when accessing it by the side-menu', () => {
      const account = mockedAccounts.at(2);
      if (account === undefined) throw new Error('Account not found');

      navigateToHome();
      cy.get('[data-testid="side-menu"]').click();
      cy.get('[data-testid="side-menu-sendreceive-icon"]').click();
      cy.url().should('eq', `${baseUrl}/${account.nickname}/transfer-coins`);
    });

    it('should render balance and amount should equal account balance', () => {
      navigateToSendForm();
      cy.get('[data-testid="available-amount"]').contains('1,000.000000000');
    });

    it('should send coins to address', () => {
      const account = mockedAccounts.at(2);
      if (account === undefined) throw new Error('Account not found');
      const recipientAccount = mockedAccounts.at(1);
      if (recipientAccount === undefined)
        throw new Error('Recipient account not found');

      navigateToSendForm();

      cy.get('[data-testid="money-field"]')
        .type('550.1234')
        .should('have.value', '550.1234');
      cy.get('[data-testid="input-field"]').type(recipientAccount.address);

      cy.get('[data-testid="button"]').contains('Send').click();
      cy.url().should('eq', `${baseUrl}/${account.nickname}/transfer-coins`);
      cy.get('[data-testid="send-confirmation"]').should('be.visible');

      cy.get('[data-testid="send-confirmation"]').should(
        'contain',
        '550.123400000 MAS',
      );

      cy.get('[data-testid="send-confirmation-recipient"]').contains(
        'AU1Z...crWY',
      );

      cy.get('[data-testid="button"]')
        .contains('Confirm and sign with password')
        .click();
      cy.url().should('eq', `${baseUrl}/${account.nickname}/home`);
    });

    it('should process any % as input', () => {
      navigateToSendForm();

      cy.get('[data-testid="send-percent-25"]').click();
      cy.get('[data-testid="money-field"]').should(
        'have.value',
        '250.000000000',
      );

      cy.get('[data-testid="send-percent-50"]').click();
      cy.get('[data-testid="money-field"]').should(
        'have.value',
        '500.000000000',
      );

      cy.get('[data-testid="send-percent-75"]').click();
      cy.get('[data-testid="money-field"]').should(
        'have.value',
        '750.000000000',
      );

      cy.get('[data-testid="send-percent-100"]').click();
      cy.get('[data-testid="money-field"]').should(
        'have.value',
        '999.990000000',
      );
      server.trackRequest = false;
    });

    it('should transfer to accounts', () => {
      const selectedAccount = mockedAccounts.at(1);
      if (selectedAccount === undefined) throw new Error('Account not found');

      navigateToTransferCoinsOfAccountIndex('0');

      cy.get('[data-testid="transfer-between-accounts"]')
        .should('exist')
        .click()
        .then(() => {
          cy.get('[data-testid="popup-modal-content"]').should('be.visible');

          cy.get('[data-testid="selector-account-0"]').should('exist').click();

          cy.get('[data-testid="input-field"]').should(
            'have.value',
            selectedAccount.address,
          );
        });
    });

    it('should refuse wrong currency input', () => {
      const account = mockedAccounts.at(2);
      if (account === undefined) throw new Error('Account not found');
      const recipientAccount = mockedAccounts.at(1);
      if (recipientAccount === undefined) throw new Error('Account not found');

      const invalidAmount = 'things';
      const tooMuch = Number(account.candidateBalance) + 1000;
      const tooLow = 0;
      const notEnoughForFees = Number(account.candidateBalance);

      navigateToSendForm();

      cy.get('[data-testid="money-field"]').should('exist').type(invalidAmount);

      cy.get('[data-testid="input-field"]').type(recipientAccount.address);
      cy.get('[data-testid="button"]').contains('Send').click();

      cy.get('[data-testid="input-field-message"]').should(
        'contain.text',
        'Invalid amount',
      );

      cy.get('[data-testid="money-field"]').clear().type(String(tooMuch));
      cy.get('[data-testid="button"]').contains('Send').click();

      cy.get('[data-testid="input-field-message"]').should(
        'contain.text',
        'Insufficient funds',
      );

      cy.get('[data-testid="money-field"]').clear().type(String(tooLow));
      cy.get('[data-testid="button"]').contains('Send').click();

      cy.get('[data-testid="input-field-message"]').should(
        'contain.text',
        'Amount must be greater than zero',
      );

      cy.get('[data-testid="money-field"]')
        .clear()
        .type(String(notEnoughForFees));
      cy.get('[data-testid="button"]').contains('Send').click();

      cy.get('[data-testid="input-field-message"]').should(
        'contain.text',
        'Insufficient funds',
      );
    });

    it('should refuse wrong address input', () => {
      const wrongAddress = 'wrong address';
      const amount = '42';

      navigateToSendForm();

      cy.get('[data-testid="money-field"]')
        .type(amount)
        .should('have.value', amount);

      cy.get('[data-testid="input-field"]').type(wrongAddress);
      cy.get('[data-testid="button"]').contains('Send').click();

      cy.get('[data-testid="input-field-message"]').should(
        'contain.text',
        'Recipient address is not a valid Massa address',
      );

      cy.get('[data-testid="input-field"]').clear();
      cy.get('[data-testid="button"]').contains('Send').click();

      cy.get('[data-testid="input-field-message"]').should(
        'contain.text',
        'Recipient is required',
      );
    });
  });
});
