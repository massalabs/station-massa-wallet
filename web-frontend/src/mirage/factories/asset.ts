import { faker } from '@faker-js/faker';
import { Factory } from 'miragejs';

export const assetFactory = Factory.extend<any>({
  name(i: number) {
    return faker.word.sample(5) + 'Token' + i++;
  },
  assetAddress: faker.string.alpha({ length: 8 }),
  symbol: faker.word.sample(5).slice(0, 3).toUpperCase(),
  decimals: 9,
  balance: faker.number.int().toString(),
});
