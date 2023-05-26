import { AccountObject } from '../models/AccountModel';

export function isAlreadyExists(nickname: string, accounts: AccountObject[]) {
  return accounts.some((account) => account.nickname === nickname);
}
