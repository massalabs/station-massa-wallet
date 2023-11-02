/* eslint-disable @typescript-eslint/no-namespace */
// ***********************************************************
// This example support/e2e.ts is processed and
// loaded automatically before your test files.
//
// This is a great place to put global configuration and
// behavior that modifies Cypress.
//
// You can change the location of this file or turn off
// automatically serving support files with the
// 'supportFile' configuration option.
//
// You can read more here:
// https://on.cypress.io/configuration
// ***********************************************************

// Import commands.js using ES2015 syntax:
import './commands';

// Alternatively you can use CommonJS syntax:
// require('./commands')

Cypress.on('window:before:load', (win) => {
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-ignore
  win.handleFromCypress = (request: any) => {
    return fetch(request.url, {
      method: request.method,
      headers: request.requestHeaders,
      body: request.requestBody,
    }).then((res) => {
      if (res && res.headers) {
        let contentType = res.headers.get('content-type');
        if (contentType && contentType.includes('application/json')) {
          return res.json().then((body) => [res.status, res.headers, body]);
        } else {
          return res.text().then((body) => [res.status, res.headers, body]);
        }
      } else {
        throw new Error('Response or headers is null or undefined');
      }
    });
  };
});

// we are intercepting miragejs calls
// https://github.com/miragejs/miragejs/issues/999#issuecomment-1495410040

import 'cypress-wait-until';

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
Cypress.Commands.add('waitForRequest', (server, urlPattern, method = '') => {
  const wildcardRegex = /\*/g;
  const questionMarkRegex = /\?/g;

  const replacedWildCard = urlPattern
    .replace(wildcardRegex, '[^ ]*')
    .replace(questionMarkRegex, '\\?');
  const regexUrl = new RegExp(replacedWildCard, 'g');
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-ignore
  const checkMatches = (r) => {
    let isMatched = r.url.match(regexUrl);
    if (method) {
      isMatched = isMatched && r.method.toLowerCase() === method.toLowerCase();
    }

    return isMatched;
  };

  const requestFound = () =>
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    server.pretender.handledRequests.some(checkMatches);
  const errorMsg = `[Request [${method}] to ${urlPattern} didn't happen`;

  return cy
    .waitUntil(requestFound, { errorMsg })
    .should('be.true')
    .log(`Request [${method}] to ${urlPattern} has been finished`);
});
