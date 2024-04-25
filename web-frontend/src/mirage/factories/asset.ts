import { faker } from '@faker-js/faker';
import { Factory } from 'miragejs';

import { MAS } from '../../const/assets/assets';

function getRandomElement(arr: string[] | number[]) {
  if (arr.length === 0) {
    throw new Error('Array is empty');
  }
  const index = Math.floor(Math.random() * arr.length);
  return arr[index];
}

export const assetFactory = Factory.extend<any>({
  name(i: number) {
    return faker.word.sample(5) + 'Token' + i++;
  },
  assetAddress: faker.string.alpha({ length: 8 }),
  symbol: getRandomElement(['WETH', MAS]),
  decimals: getRandomElement([9, 18, 6]),
  balance: faker.number.int().toString(),
});
