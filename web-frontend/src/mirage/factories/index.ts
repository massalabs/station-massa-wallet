import { accountFactory } from './account';
import { assetFactory } from './asset';
import { transferFactory } from './transfer';

export const factories = {
  account: accountFactory,
  transfer: transferFactory,
  asset: assetFactory,
};
