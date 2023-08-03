import { Model } from 'miragejs';
import { ModelDefinition } from 'miragejs/-types';

import { AccountObject } from '../../models/AccountModel';
import { TransferModel } from '../../models/TransferModel';

const accountModel: ModelDefinition<AccountObject> = Model.extend({});
const transferModel: ModelDefinition<TransferModel> = Model.extend({});

export const models = {
  account: accountModel,
  transfer: transferModel,
};
