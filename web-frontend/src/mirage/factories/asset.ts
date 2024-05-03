import { faker } from '@faker-js/faker';
import { Factory } from 'miragejs';

import { MAS } from '../../const/assets/assets';

export const assetFactory = Factory.extend({
  name(i: number) {
    return faker.word.sample(5) + 'Token' + i++;
  },
  address: faker.string.alpha({ length: 8 }),
  symbol: MAS,
  decimals: 9,
  balance: '1000000000000',
});
