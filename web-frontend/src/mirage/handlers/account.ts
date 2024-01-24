import { faker } from '@faker-js/faker';
import { Server, Response } from 'miragejs';

import { AppSchema } from '../types';

function accountObject(nickname: string) {
  return {
    nickname,
    candidateBalance: faker.number.int().toString(),
    balance: faker.number.int().toString(),
    // TODO: Generate valid address
    address: 'AU1ZMBZeARHYMFfV4uvbyCB85DAUPr2BJXzU1kSYdwCKKrY5crWY',
    keyPair: {
      // Standard info: https://github.com/massalabs/massa-standards/blob/main/wallet/file-format.md
      publicKey: 'P' + faker.string.alpha({ length: 50 }),
      privateKey: 'S' + faker.string.alpha({ length: 50 }),
      nonce: faker.string.alpha({ length: 30 }),
      salt: faker.string.alpha({ length: 30 }),
    },
  };
}

export function routesForAccounts(server: Server) {
  server.get(
    'accounts',
    (schema: AppSchema) => {
      const { models: accounts } = schema.all('account');

      return accounts;
    },
    { timing: 500 },
  );

  server.get(
    'accounts/:nickname',
    (schema: AppSchema, request) => {
      const { nickname } = request.params;
      const data = schema.findBy('account', { nickname });

      if (!data)
        return new Response(
          404,
          {},
          { code: '404', error: 'Account not found' },
        );

      const account = data.attrs;

      return { ...account };
    },
    { timing: 500 },
  );

  server.delete(
    'accounts/:nickname',
    (schema: AppSchema, request) => {
      const { nickname } = request.params;
      const accountName = schema.findBy('account', { nickname });

      if (!accountName)
        return new Response(
          404,
          {},
          { code: '404', error: 'Account not found' },
        );

      accountName.destroy();

      return new Response(200);
    },
    { timing: 500 },
  );

  server.get(
    'accounts/:nickname/assets',
    (schema: AppSchema, request) => {
      const { nickname } = request.params;
      const accounts: any = schema.findBy('account', { nickname });

      if (!accounts)
        return new Response(
          404,
          {},
          { code: '404', error: 'Failed to retrieve assets' },
        );

      return accounts.assets.models;
    },
    { timing: 500 },
  );

  server.delete(
    'accounts/:nickname/assets',
    (schema: AppSchema, request) => {
      let assetAddress = request.queryParams.assetAddress;
      const { nickname } = request.params;
      const accounts: any = schema.findBy('account', { nickname });

      if (!nickname)
        return new Response(
          402,
          {},
          { code: '402', error: 'Provide nickname in query' },
        );

      if (!accounts)
        return new Response(
          404,
          {},
          { code: '404', error: 'Failed to find account' },
        );

      if (!assetAddress)
        return new Response(
          422,
          {},
          { code: '422', error: 'Address provided Invalid' },
        );

      accounts.assets.models.at(0).destroy();

      return new Response(204);
    },
    { timing: 500 },
  );

  server.post(
    'accounts/:nickname/assets',
    (schema: AppSchema, request) => {
      let assetAddress = request.queryParams.assetAddress;
      const { nickname } = request.params;
      const account: any = schema.findBy('account', { nickname });

      if (!account)
        return new Response(
          400,
          {},
          { code: '400', error: 'Failed to find account' },
        );

      if (!nickname) {
        return new Response(401, {
          code: '0001',
          message: 'missing nickname',
        });
      }

      if (!assetAddress)
        return new Response(
          402,
          {},
          { code: '402', error: 'Provide an asset address in query' },
        );

      // const newAsset = schema.create('asset', { account });

      return new Response(201);
    },
    { timing: 500 },
  );

  server.put('accounts', (schema: AppSchema) => {
    return schema.create(
      'account',
      accountObject(
        faker.internet.userName({
          firstName: 'Imported',
          lastName: 'Imported',
        }),
      ),
    );
  });

  server.put('accounts/:nickname', (schema: AppSchema, request) => {
    const { nickname: newNickname } = JSON.parse(request.requestBody);
    const { nickname } = request.params;
    const account = schema.findBy('account', { nickname });

    if (!account)
      return new Response(404, {}, { code: '404', error: 'Account not found' });

    account.update({ nickname: newNickname });

    return new Response(
      200,
      {},
      { code: '200', message: 'Account edited successfully' },
    );
  });

  server.post('accounts/:nickname', (schema: AppSchema, request) => {
    const { nickname } = request.params;

    return schema.create('account', accountObject(nickname));
  });

  server.post('accounts/:nickname/backup', (schema: AppSchema, request) => {
    const { nickname } = request.params;
    const account = schema.findBy('account', { nickname });

    if (!account)
      return new Response(404, {}, { code: '404', error: 'Account not found' });

    return {
      privateKey: account.keyPair.privateKey,
      publicKey: account.keyPair.publicKey,
    };
  });

  server.post(
    'accounts/:nickname/transfer',
    (schema: AppSchema, request) => {
      const { nickname } = request.params;
      const { amount, fee, recipientAddress } = JSON.parse(request.requestBody);

      if (!amount || !fee || !recipientAddress) {
        return new Response(400, {
          code: '0001',
          message: 'missing fields',
        });
      }

      const account = schema.findBy('account', { nickname });

      if (!account) {
        return new Response(
          404,
          {},
          { code: '404', error: 'Account not found' },
        );
      }

      return schema.create('transfer');
    },
    { timing: 500 },
  );
}
