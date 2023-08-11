import { Model, belongsTo, hasMany } from 'miragejs';
import { ModelDefinition } from 'miragejs/-types';

import { AccountObject } from '../../models/AccountModel';
import { TransferModel } from '../../models/TransferModel';
import { IToken } from '@/models/AssetModel';

const accountModel: ModelDefinition<AccountObject> = Model.extend({
  assets: hasMany(),
});
const transferModel: ModelDefinition<TransferModel> = Model.extend({});
const assetModel: ModelDefinition<IToken> = Model.extend({
  account: belongsTo(),
});

export const models = {
  account: accountModel,
  transfer: transferModel,
  asset: assetModel,
};
