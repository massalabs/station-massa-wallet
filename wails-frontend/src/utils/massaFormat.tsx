import { fromMAS, toMAS } from '@massalabs/massa-web3';
import { formatValue } from 'react-currency-input-field';

/**
 * Enumeration for unit options.
 */
export enum Unit {
  MAS = 'MAS',
  NanoMAS = 'NanoMAS',
}

export const masToken = 'MAS';

export const presetFees: { [key: string]: string } = {
  low: '1',
  standard: '1000',
  high: '5000',
};

export function toMASS(num: string | number): number {
  return Number(toMAS(num));
}

export function toNanoMASS(str: string): number {
  const formattedString = str?.replace(/[^0-9.-]/g, ''); // Remove non-numeric characters

  // fromMAS -> MassaToNano
  // toMASS -> NanotoMASSsa
  return Number(fromMAS(formattedString));
}

export function removeTrailingZeros(numStr: string): string {
  return numStr.replace(/\.?0+$/, '');
}

/**
 * Formats a number according to the specified unit and formatting options.
 * @param value - The number to format, which can be either a string or a bigint.
 * @param unit - The unit to use for formatting. Defaults to `Unit.MAS`.
 * @returns The formatted number as a string.
 */
export function formatStandard(
  value: string | bigint,
  unit = Unit.MAS,
): string {
  if (typeof value === 'bigint') {
    value = value.toString();
  }
  const numInMas = unit === Unit.MAS ? value : toMAS(value).toString();

  // Format the number as a string with separators
  const formattedNum = formatValue({
    value: numInMas,
    groupSeparator: ',',
    decimalSeparator: '.',
    decimalScale: 9,
  });
  return removeTrailingZeros(formattedNum);
}

/**
 * Reverses the formatting of a number string and converts it to a number.
 * @param str - The formatted number string to reverse.
 * @returns The reversed and parsed number.
 */

export function reverseFormatStandard(str: string): number {
  const formattedString = str?.replace(/[^0-9.-]/g, ''); // Remove non-numeric characters
  return Number(formattedString);
}

/**
 * Checks if a recipient address is in the correct format.
 * @param recipient - The recipient address to check.
 * @returns `true` if the address is in the correct format, `false` otherwise.
 */

export function checkAddressFormat(recipient: string): boolean {
  return /^AU[a-zA-Z0-9]{4,}$/.test(recipient);
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
