/* eslint-disable @typescript-eslint/ban-ts-comment */
// @ts-nocheck
import { mockServer } from '../../../../src/mirage';
import { AccountObject } from '@/models/AccountModel';

describe('E2E | Acceptance | Home | Receive', () => {
  let server: any;
  const baseUrl = Cypress.config().baseUrl;

  beforeEach(() => {
    server = mockServer('test');
  });

  afterEach(() => {
    server.shutdown();
  });

  describe('receive', () => {
    let mockedAccounts;

    function mockAccounts() {
      const mockedAccounts = server.createList('account', 2);
      mockedAccounts.forEach((account) => {
        server.createList('asset', 3, { account });
      });
      const account = {
        nickname: 'Mario',
        address: 'AU1ZMBZeARHYMFfV4uvbyCB85DAUPr2BJXzU1kSYdwCKKrY5crWY',
      };
      mockedAccounts.push(server.create('account', { ...account }));

      return mockedAccounts;
    }

    beforeEach(() => {
      mockedAccounts = mockAccounts();
    });

    function navigateToReceivePage() {
      cy.visit('/');

      cy.get('[data-testid="account-2"]').click();
      cy.get('[data-testid="receive-button"]').click();
    }

    function generateLink(amount) {
      cy.get('[data-testid="button"]')
        .should('exist')
        .contains('Generate a link')
        .click();

      cy.get('[data-testid="amount-to-send"]').type(amount);

      cy.get('[data-testid="generate-link-button"]').click();

      cy.on('window:confirm', () => true);

      cy.window().then((win) => {
        cy.stub(win, 'prompt').returns(win.prompt).as('copyToClipboardPrompt');
      });

      cy.get('[data-testid="clipboard-link"]').click();
    }

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

    it('should land on receive page when receive CTA is clicked', () => {
      const account = mockedAccounts.at(2);

      cy.visit('/');

      cy.get('[data-testid="account-2"]').click();

      cy.get('[data-testid="receive-button"]').click();
      cy.url().should('eq', `${baseUrl}/${account.nickname}/transfer-coins`);

      cy.get('[data-testid="send-coins"]').should('not.be.visible');
      cy.get('[data-testid="receive-coins"]').should('be.visible');
    });

    it('clipboard field should contain user address', () => {
      const account = mockedAccounts.at(2);

      navigateToReceivePage();
      cy.get('[data-testid="clipboard-field"]').contains(account.address);
    });

    it('should copy from clipboard', () => {
      const account = mockedAccounts.at(2);
      navigateToReceivePage();
      cy.on('window:confirm', () => true);

      cy.window().then((win) => {
        cy.stub(win, 'prompt').returns(win.prompt).as('copyToClipboardPrompt');
      });

      cy.get('[data-testid="clipboard-field"]').click();
      cy.get('@copyToClipboardPrompt').should('be.called');

      cy.get('@copyToClipboardPrompt').should((prompt) => {
        expect(prompt.args[0][1]).to.equal(account.address);
      });
    });

    it('should generate and copy link', () => {
      const amount = 5000;
      const account = mockedAccounts.at(2);

      const generatedLink = `http://localhost:8080/send-redirect/?to=${account.address}&amount=${amount}`;

      navigateToReceivePage();
      cy.get('[data-testid="button"]').contains('Generate a link').click();

      cy.get('[data-testid="popup-modal"]').should('be.visible');

      cy.get('[data-testid="amount-to-send"]')
        .should('have.attr', 'placeholder', 'Amount to ask')
        .type(amount);

      cy.get('[data-testid="clipboard-field"]').should('contain', '');

      cy.get('[data-testid="generate-link-button"]').click();

      cy.on('window:confirm', () => true);

      cy.window().then((win) => {
        cy.stub(win, 'prompt').returns(win.prompt).as('copyToClipboardPrompt');
      });

      cy.get('[data-testid="clipboard-link"]')
        .should('contain', generatedLink)
        .click();

      cy.get('@copyToClipboardPrompt').should('be.called');

      cy.get('@copyToClipboardPrompt').should((prompt) => {
        expect(prompt.args[0][1]).to.equal(generatedLink);
      });
    });

    it('should redirect to generated link and fill input fields', () => {
      const account = mockedAccounts.at(2);
      const amount = 5000;

      const generatedLink = `http://localhost:8080/send-redirect/?to=${account.address}&amount=${amount}`;

      navigateToReceivePage();

      generateLink(amount);

      cy.visit(generatedLink);

      cy.get('[data-testid="money-field"').should(
        'have.value',
        customFormatNumber(amount),
      );

      cy.get('[data-testid="input-field"').should(
        'have.value',
        account.address,
      );
    });
  });
});
