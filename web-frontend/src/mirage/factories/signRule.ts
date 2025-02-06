import { faker } from '@faker-js/faker';
import { Factory } from 'miragejs';

import { RuleType } from '../../models/ConfigModel';

export const signRuleFactory = Factory.extend({
  id() {
    return faker.string.uuid();
  },
  ruleType() {
    return faker.helpers.arrayElement([
      RuleType.DisablePasswordPrompt,
      RuleType.AutoSign,
    ]);
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
