import { toMAS } from '@massalabs/massa-web3';
import { formatValue } from 'react-currency-input-field';

export const masToken = 'MAS';

/**
 * Enumeration for unit options.
 */
export enum Unit {
  MAS = masToken,
  NanoMAS = 'NanoMAS',
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

export function maskNickname(str: string, length = 10): string {
  if (!str) return '';

  if (str.length <= length) return str;

  return str?.substring(0, length) + '...';
}
