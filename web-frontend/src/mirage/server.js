import { createServer, Model, Factory } from 'miragejs';
import { ENV } from '../const/env/env';
import { faker } from '@faker-js/faker';

// eslint-disable-next-line jsdoc/require-jsdoc
function mockServer(environment = ENV.DEV) {
  let mockedServer = createServer({
    environment,

    models: {
      account: Model,
    },

    factories: {
      account: Factory.extend({
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

        return schema.findBy('account', { nickname });
      });

      this.put('accounts', (schema) => {
        return schema.create('account', {
          nickname: 'imported',
          candidateBalance: '117',
        });
      });

      this.put('accounts/:nickname', (schema, request) => {
        let { newNickname } = JSON.parse(request.requestBody);
        let nickname = request.params.nickname;
        let account = schema.accounts.findBy({ nickname });

        return account.update({ nickname: newNickname });
      });
    },
  });

  return mockedServer;
}

export default mockServer;
