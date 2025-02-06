import { accountFactory } from './account';
import { assetFactory } from './asset';
import { signRuleFactory } from './signRule';
import { transferFactory } from './transfer';

export const factories = {
  account: accountFactory,
  transfer: transferFactory,
  asset: assetFactory,
  signRule: signRuleFactory,
};
