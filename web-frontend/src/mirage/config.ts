import { createServer } from 'miragejs';

import { factories } from './factories';
import { handlers } from './handlers';
import { models } from './models';
import { ENV } from '../const/env/env';
import { AccountObject } from '@/models/AccountModel';

export function mockServer(environment = ENV.DEV) {
  const server = createServer({
    environment,
    models,
    factories,
    seeds(server) {
      const accounts = server.createList('account', 5);
      accounts.forEach((account: AccountObject) => {
        server.createList('asset', 3, { account });
      });
    },
  });

  for (const namespace of Object.keys(handlers)) {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    handlers[namespace](server);
  }

  return server;
}
