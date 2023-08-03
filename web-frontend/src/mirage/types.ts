import { Registry } from 'miragejs';
import Schema from 'miragejs/orm/schema';

import { factories } from './factories';
import { models } from './models';

type AppRegistry = Registry<typeof models, typeof factories>;
export type AppSchema = Schema<AppRegistry>;
