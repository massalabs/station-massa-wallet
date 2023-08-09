import { faker } from '@faker-js/faker';
import { Factory } from 'miragejs';

import { AccountObject, IToken } from '../../models/AccountModel';

export const accountFactory = Factory.extend<AccountObject>({
  nickname() {
    return faker.internet.userName();
  },
  candidateBalance() {
    return faker.number.int().toString();
  },
  balance() {
    return faker.number.int().toString();
  },
  address() {
    return 'AU' + faker.string.alpha({ length: { min: 50, max: 53 } });
  },
  keyPair() {
    return {
      publicKey: 'P' + faker.string.alpha({ length: 50 }),
      privateKey: 'S' + faker.string.alpha({ length: 50 }),
      nonce: faker.string.alpha({ length: 30 }),
      salt: faker.string.alpha({ length: 30 }),
    };
  },
  assets() {
    const initialTokens: IToken[] = [
      {
        name: faker.word.sample(5) + 'Token',
        assetAddress: faker.string.alpha({ length: 8 }),
        symbol: faker.word.sample(5).slice(0, 3).toUpperCase(),
        decimals: 9,
        balance: faker.number.int().toString(),
      },
      {
        name: faker.word.sample(5) + 'Token',
        assetAddress: faker.string.alpha({ length: 8 }),
        symbol: faker.word.sample(5).slice(0, 3).toUpperCase(),
        decimals: 9,
        balance: faker.number.int().toString(),
      },
    ];
    return initialTokens;
  },
  status() {
    return 'ok';
  },
});
