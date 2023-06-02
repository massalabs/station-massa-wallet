import { toMAS } from '@massalabs/massa-web3';

/**
 * Enumeration for unit options.
 */
export enum Unit {
  MAS = 'MAS',
  NanoMAS = 'NanoMAS',
}

/**
 * Formats a number according to the specified unit and formatting options.
 * @param num - The number to format.
 * @param unit - The unit to use for formatting. Defaults to `Unit.MAS`.
 * @param maximumFractionDigits - The maximum number of fraction digits to display. Defaults to `2`.
 * @returns The formatted number as a string.
 */
export function formatStandard(
  num: number,
  unit = Unit.MAS,
  maximumFractionDigits = 2,
): string {
  const numInMas = unit === Unit.MAS ? num : toMAS(num).toNumber();
  const locale = localStorage.getItem('locale') || 'en-US';
  return numInMas.toLocaleString(locale, {
    minimumFractionDigits: 2,
    maximumFractionDigits,
  });
}

/**
 * Reverses the formatting of a number string and converts it to a number.
 * @param str - The formatted number string to reverse.
 * @returns The reversed and parsed number.
 */
export function reverseFormatStandard(str: string): number {
  const formattedString = str.replace(/[^0-9.-]/g, ''); // Remove non-numeric characters
  return parseFloat(formattedString);
}

/**
 * Checks if a recipient address is in the correct format.
 * @param recipient - The recipient address to check.
 * @returns `true` if the address is in the correct format, `false` otherwise.
 */
export function checkRecipientFormat(recipient: string): boolean {
  const regex = /^AU[a-zA-Z0-9]{51}$/;
  return regex.test(recipient);
}
