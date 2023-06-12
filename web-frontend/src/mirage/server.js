import { createServer, Model, Factory } from 'miragejs';
import { ENV } from '../const/env/env';
import { faker } from '@faker-js/faker';

// eslint-disable-next-line jsdoc/require-jsdoc
function mockServer(environment = ENV.DEV) {
  let mockedServer = createServer({
    environment,

    models: {
      account: Model,
      transfer: Model,
    },

    factories: {
      transfer: Factory.extend({
        operationId() {
          return faker.number.int().toString();
        },
      }),
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

      this.get(
        'accounts',
        (schema) => {
          let { models: accounts } = schema.accounts.all();

          return accounts;
        },
        { timing: 1000 },
      );

      this.get('accounts/:nickname', (schema, request) => {
        let nickname = request.params.nickname;

        let { attrs: account } = schema.findBy('account', { nickname });

        return { ...account };
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

      this.post('accounts/:nickname', (schema, request) => {
        let nickname = request.params.nickname;

        return schema.accounts.create({ nickname });
      });

      this.post('accounts/:nickname/backup', {});

      this.post('accounts/:nickname/transfer', (schema, request) => {
        let { amount, fee, recipientAddress } = JSON.parse(request.requestBody);

        if (!amount || !fee || !recipientAddress) {
          // TODO
          // we must handle here the missing payload fields
        }

        return schema.create('transfer');
      });
    },
  });

  return mockedServer;
}

export default mockServer;
