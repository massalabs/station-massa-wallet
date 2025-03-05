/* eslint-disable @typescript-eslint/ban-ts-comment */
// @ts-nocheck
import { faker } from '@faker-js/faker';
import { createServer, Response } from 'miragejs';

import { factories } from './factories';
import { handlers, otherDomainHandlers } from './handlers';
import { models } from './models';
import { ENV } from '../const/env/env';
import { AccountObject } from '@/models/AccountModel';

const serverData = {
  trackRequests: true,
  models,
  factories,
  seeds(server) {
    const accounts: AccountObject[] = server.createList('account', 5);
    accounts.forEach((account: AccountObject) => {
      server.createList('asset', 3, { account });

      const numRules = faker.number.int({ min: 0, max: 8 });
      server.createList('signRule', numRules, {
        accountNickname: account.nickname,
      });
    });
  },
};

export function mockServer(environment = ENV.DEV) {
  const server = createServer({ ...serverData, environment });

  for (const namespace of Object.keys(handlers)) {
    handlers[namespace](server);
  }

  otherDomainHandlers(server);

  return server;
}

export function mockServerWithCypress() {
  if (window.Cypress) {
    // If your app makes requests to domains other than / (the current domain), add them
    // here so that they are also proxied from your app to the handleFromCypress function.
    // For example: let otherDomains = ["https://my-backend.herokuapp.com/"]
    let otherDomains = [];
    let methods = ['get', 'put', 'patch', 'post', 'delete'];

    const server = createServer({
      ...serverData,
      environment: 'test',
      routes() {
        for (const domain of ['/', ...otherDomains]) {
          for (const method of methods) {
            this[method](`${domain}*`, async (schema, request) => {
              let [status, headers, body] = await window.handleFromCypress(
                request,
              );
              return new Response(status, headers, body);
            });
          }
        }
        otherDomainHandlers(this);
      },
    });
  }
}
