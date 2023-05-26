import { createServer, Model, Factory } from 'miragejs';
import { ENV } from '../const/env/env';
import { AccountObject } from '../models/AccountModel';
import { faker } from '@faker-js/faker';

function mockServer(environment = ENV.DEV) {
  let mockedServer = createServer({
    environment,

    models: {
      account: Model.extend<Partial<AccountObject>>({}),
    },

    factories: {
      account: Factory.extend<Partial<AccountObject>>({
        nickname() {
          return faker.internet.userName();
        },
        candidateBalance() {
          return faker.number.int().toString();
        },
      }),
    },

    seeds(server) {
      server.createList('account', 2);
    },

    routes() {
      this.namespace = import.meta.env.VITE_BASE_API;

      this.get('accounts', (schema) => {
        return schema.all('account').models;
      });

      this.get('accounts/:nickname', (schema, request) => {
        let nickname = request.params.nickname;

        return schema.find('account', nickname);
      });

      this.put('accounts', (schema) => {
        return schema.create('account', {
          nickname: 'imported',
          candidateBalance: '117',
        });
      });
    },
  });

  return mockedServer;
}

export default mockServer;
