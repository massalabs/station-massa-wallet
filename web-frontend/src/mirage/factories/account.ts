import { faker } from '@faker-js/faker';
import { Factory } from 'miragejs';

export const accountFactory = Factory.extend<any>({
  nickname: (i: number) => faker.internet.userName() + i,
  candidateBalance: (i: number) =>
    (faker.number.int() * (i > 0 ? i / 10 : 0)).toString(),
  balance: faker.number.int().toString(),
  address: 'AU' + faker.string.alpha({ length: { min: 50, max: 53 } }),
  keyPair() {
    return {
      publicKey: 'P' + faker.string.alpha({ length: 50 }),
      privateKey: 'S' + faker.string.alpha({ length: 50 }),
      nonce: faker.string.alpha({ length: 30 }),
      salt: faker.string.alpha({ length: 30 }),
    };
  },
  status: 'ok',
});
