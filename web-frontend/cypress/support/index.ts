export {};

declare global {
  namespace Cypress {
    interface Chainable {
      /**
       * Custom command to select DOM element by data-cy attribute.
       * @example cy.dataCy('greeting')
       */
      waitForRequest(
        server: string,
        urlPattern: string,
        method: string,
      ): Chainable<JQuery<HTMLElement>>;
    }
  }
}
