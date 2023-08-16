import { Model, belongsTo, hasMany } from 'miragejs';
import { ModelDefinition } from 'miragejs/-types';

import { AccountObject } from '../../models/AccountModel';
import { TransferModel } from '../../models/TransferModel';
import { Token } from '@/models/AssetModel';

const accountModel: ModelDefinition<AccountObject> = Model.extend({
  assets: hasMany(),
});
const transferModel: ModelDefinition<TransferModel> = Model.extend({});
const assetModel: ModelDefinition<Token> = Model.extend({
  account: belongsTo(),
});

export const models = {
  account: accountModel,
  transfer: transferModel,
  asset: assetModel,
};
