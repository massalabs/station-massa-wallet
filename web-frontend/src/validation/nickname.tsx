import { AccountObject } from '../models/AccountModel';

export function isAlreadyExists(
  nickname: string,
  accounts: AccountObject[],
): boolean {
  return accounts.some((account) => account.nickname === nickname);
}

/**
 * This function validates the nickname using the following rules:
 * - must contain only alphanumeric characters, underscores and dashes
 * - must be between 1 and 32 characters long
 *
 * @param nickname nickname to validate
 * @returns true if nickname is valid, false otherwise
 */
export function isNicknameValid(nickname: string): boolean {
  return /^[a-zA-Z0-9_-]+$/.test(nickname) && nickname.length <= 32;
}
