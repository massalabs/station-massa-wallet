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
