import { Address } from '@massalabs/massa-web3/';

/**
 * Checks if a recipient address is in the correct format.
 * @param recipient - The recipient address to check.
 * @returns `true` if the address is in the correct format, `false` otherwise.
 */

export function checkAddressFormat(recipient: string): boolean {
  try {
    // eslint-disable-next-line no-new
    new Address(recipient);
  } catch (error) {
    return false;
  }

  return true;
}

/**
 * Masks the middle of an address with a specified character.
 * @param str - The address to mask.
 * @param mask - The character to use for masking. Defaults to `.`.
 * @returns The masked address.
 */

export function maskAddress(str: string, length = 4, mask = '. . .'): string {
  const start = length;
  const end = str?.length - length;

  return str ? str?.substring(0, start) + mask + str?.substring(end) : '';
}

export function maskNickname(str: string, length = 32): string {
  if (!str) return '';

  if (str.length <= length) return str;

  return str?.substring(0, length) + '...';
}
