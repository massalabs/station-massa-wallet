import { AccountObject } from '../models/AccountModel';

export function isAlreadyExists(
  nickname: string,
  accounts: AccountObject[],
): boolean {
  return accounts.some((account) => account.nickname === nickname);
}

export function isNicknameValid(nickname: string): boolean {
  return /^[a-z0-9_-]$/.test(nickname);
}
