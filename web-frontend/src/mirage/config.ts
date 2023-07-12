import { createServer } from 'miragejs';
import { handlers } from './handlers';
import { models } from './models';
import { factories } from './factories';
import { ENV } from '../const/env/env';

export function mockServer(environment = ENV.DEV) {
  const server = createServer({
    environment,
    models,
    factories,
    seeds(server) {
      server.createList('account', 5);
    },
  });

  for (const namespace of Object.keys(handlers)) {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    handlers[namespace](server);
  }

  return server;
}
