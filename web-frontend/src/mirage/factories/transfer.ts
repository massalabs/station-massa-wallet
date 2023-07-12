import { Factory } from 'miragejs';
import { faker } from '@faker-js/faker';
import { TransferModel } from '../../models/TransferModel';

export const transferFactory = Factory.extend<TransferModel>({
  operationId() {
    return faker.number.int().toString();
  },
});
