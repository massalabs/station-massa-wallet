import { Factory } from 'miragejs';
import { faker } from '@faker-js/faker';
import { AccountObject } from '../../models/AccountModel';

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
});
