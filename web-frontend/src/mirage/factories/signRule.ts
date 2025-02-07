import { faker } from '@faker-js/faker';
import { Factory } from 'miragejs';

export const signRuleFactory = Factory.extend({
  id() {
    return faker.string.uuid();
  },
  ruleType() {
    return faker.helpers.arrayElement(['disable_password_prompt', 'auto_sign']);
  },
  contract() {
    return 'AS' + faker.string.alphanumeric(50);
  },
  enabled() {
    return faker.datatype.boolean();
  },
  name() {
    return faker.helpers.maybe(() => faker.word.words(2), { probability: 0.9 });
  },
});
