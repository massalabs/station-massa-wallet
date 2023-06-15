import { createServer, Model, Factory, Response } from 'miragejs';
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
        address() {
          return 'AU' + faker.string.alpha({ length: { min: 128, max: 256 } });
        },
      }),
    },

    seeds(server) {
      server.createList('account', 5);
    },

    routes() {
      this.namespace = import.meta.env.VITE_BASE_API;

      this.get(
        'accounts',
        (schema) => {
          let { models: accounts } = schema.accounts.all();

          return accounts;
        },
        { timing: 500 },
      );

      this.get(
        'accounts/:nickname',
        (schema, request) => {
          let nickname = request.params.nickname;

          let { attrs: account } = schema.findBy('account', { nickname });

          return { ...account };
        },
        { timing: 500 },
      );

      this.put('accounts', (schema) => {
        return schema.create('account', {
          nickname: 'imported',
          candidateBalance: '117',
        });
      });

      this.put('accounts/:nickname', (schema, request) => {
        let { nickname: newNickname } = JSON.parse(request.requestBody);
        let nickname = request.params.nickname;
        let account = schema.accounts.findBy({ nickname });

        return account.update({ nickname: newNickname });
      });

      this.post('accounts/:nickname', (schema, request) => {
        let nickname = request.params.nickname;

        return schema.accounts.create({ nickname });
      });

      this.post('accounts/:nickname/backup', {});

      this.post(
        'accounts/:nickname/transfer',
        (schema, request) => {
          let { amount, fee, recipientAddress } = JSON.parse(
            request.requestBody,
          );

          if (!amount || !fee || !recipientAddress) {
            return new Response(400, {
              code: '0001',
              message: 'missing fields',
            });
          }

          return schema.create('transfer');
        },
        { timing: 500 },
      );
    },
  });

  return mockedServer;
}

export default mockServer;
