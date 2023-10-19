/* eslint-disable @typescript-eslint/ban-ts-comment */
// @ts-nocheck
import {
  HttpResponseInterceptor,
  RouteMatcher,
  StaticResponse,
} from 'cypress/types/net-stubbing';

// how to use it
// it('should land in selected account when pick one account', () => {
//   const interception = interceptIndefinitely('/account-select');

//   cy.visit('/account-select');
//   cy.get('[data-testid="account-1"]').click().then(() => {
//     // release the answer
//     interception.sendResponse();
//     // do something
//   });
// });

export function interceptIndefinitely(
  requestMatcher: RouteMatcher,
  response?: StaticResponse | HttpResponseInterceptor,
): { sendResponse: () => void } {
  let sendResponse;

  const trigger = new Promise((resolve) => {
    sendResponse = resolve;
  });

  cy.intercept(requestMatcher, (request) => {
    return trigger.then(() => {
      request.reply(response);
    });
  });

  return { sendResponse };
}
