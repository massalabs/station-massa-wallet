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
  num: string,
  unit = Unit.MAS,
  maximumFractionDigits = 2,
): string {
  let numInMas = parseInt(num);
  numInMas = unit === Unit.MAS ? numInMas : toMAS(num).toNumber();
  const locale = localStorage.getItem('locale') || 'en-US';
  return numInMas.toLocaleString(locale, {
    minimumFractionDigits: 2,
    maximumFractionDigits,
  });
}
