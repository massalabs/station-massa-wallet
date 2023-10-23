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
