import { Model, belongsTo, hasMany } from 'miragejs';
import { ModelDefinition } from 'miragejs/-types';

import { AccountObject } from '../../models/AccountModel';
import { TransferModel } from '../../models/TransferModel';
import { Asset } from '@/models/AssetModel';
import { SignRule } from '@/models/ConfigModel';

const accountModel: ModelDefinition<AccountObject> = Model.extend({
  assets: hasMany(),
});
const transferModel: ModelDefinition<TransferModel> = Model.extend({});
const assetModel: ModelDefinition<Asset> = Model.extend({
  account: belongsTo(),
});

interface SignRuleModel extends SignRule {
  accountNickname: string;
}

const signRuleModel: ModelDefinition<SignRuleModel> = Model.extend({
  account: belongsTo(),
});

export const models = {
  account: accountModel,
  transfer: transferModel,
  asset: assetModel,
  signRule: signRuleModel,
};
