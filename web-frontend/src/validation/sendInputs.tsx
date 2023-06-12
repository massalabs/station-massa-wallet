import { fromMAS } from '@massalabs/massa-web3';
import {
  Unit,
  checkAddressFormat,
  reverseFormatStandard,
} from '../utils/MassaFormating';

export interface SendInputsErrors {
  unexpectedError?: string;
  amount?: string;
  address?: string;
  fees?: string;
}

/**
 * Checks if the given amount is positive.
 * @param amount The amount to check.
 * @returns A boolean indicating if the amount is positive.
 */
function isPositiveAmount(amount: number): boolean {
  const amountIsPositive = amount > 0;
  return amountIsPositive;
}

/**
 * Checks if the given number is a valid amount.
 * @param num The number to check.
 * @returns A boolean indicating if the number is valid.
 */
function isValidAmount(num: number): boolean {
  const numberIsNotNaN = !isNaN(num);
  return numberIsNotNaN;
}

/**
 * Checks if the given amount is an integer.
 * @param amount The amount to check.
 * @returns A boolean indicating if the amount is an integer.
 */
function isAmountAnInteger(amount: number): boolean {
  return amount % 1 === 0;
}

/**
 * Checks if the balance is higher than the given amount.
 * @param balance The balance of the account.
 * @param amount The amount to check.
 * @returns A boolean indicating if the balance is higher than the amount.
 */
function isBalanceHigherThanAmount(balance: number, amount: number): boolean {
  return balance > amount;
}

/**
 * Gets the error message for the amount input.
 * @param amount The amount to check.
 * @param unit The unit of the amount.
 * @returns The error message or null if there is no error.
 */
export const getAmountFormatError = (
  amount: string,
  unit: Unit = Unit.NanoMAS,
): string | null => {
  const unformattedAmount = reverseFormatStandard(amount);
  if (!isValidAmount(unformattedAmount)) {
    return 'errors.send.invalid-amount';
  }
  if (!isPositiveAmount(unformattedAmount)) {
    return 'errors.send.amount-to-low';
  }
  const amountInNMAS =
    unit === Unit.MAS
      ? fromMAS(unformattedAmount).toString()
      : unformattedAmount;
  if (!isAmountAnInteger(+amountInNMAS)) {
    return 'errors.send.invalid-amount-decimals';
  }
  return null;
};

/**
 * Checks if the amount is higher than the balance.
 * @param amount The amount to check.
 * @param balance The balance of the account.
 * @param unit The unit of the amount.
 * @returns The error message or null if there is no error.
 */
export function getAmountTooHighError(
  amount: string,
  balance: number,
  unit: Unit = Unit.NanoMAS,
): string | null {
  const unformattedAmount = reverseFormatStandard(amount);
  const amountInNMAS =
    unit === Unit.MAS
      ? fromMAS(unformattedAmount).toString()
      : unformattedAmount;
  if (!isBalanceHigherThanAmount(balance, +amountInNMAS)) {
    return 'errors.send.amount-to-high';
  }
  return null;
}

/**
 * Checks if the address is valid.
 * @param recipient The address to check.
 * @returns The error message or null if there is no error.
 */
export function getAddressError(recipient: string): string | null {
  if (recipient === '') {
    return 'errors.send.no-address';
  }
  if (!checkAddressFormat(recipient)) {
    return 'errors.send.invalid-address';
  }
  return null;
}
