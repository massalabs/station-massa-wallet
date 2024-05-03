import { faker } from '@faker-js/faker';
import { Factory } from 'miragejs';

export const accountFactory = Factory.extend({
  nickname: (i: number) => faker.internet.userName() + i,
  candidateBalance: (i: number) =>
    Math.trunc(faker.number.int() * (i > 0 ? i / 10 : 0)).toString(),
  balance: faker.number.int().toString(),
  // TODO: Generate valid address
  address: 'AU1ZMBZeARHYMFfV4uvbyCB85DAUPr2BJXzU1kSYdwCKKrY5crWY',
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
