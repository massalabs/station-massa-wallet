import { fromMAS } from '@massalabs/massa-web3';
import { checkAddressFormat } from '../utils/MassaFormating';

export interface SendInputsErrors {
  unexpectedError?: string;
  amount?: string;
  address?: string;
  fees?: string;
}


export function getAddressError(
  recipient: string,
): string | null{
  if (!recipient) {
    return 'errors.send.no-address'

  }
  if (!checkAddressFormat(recipient)) {
    return 'errors.send.invalid-address'
  }
  return null;
}

function isPositiveAmount(amount: number): boolean {
  const amountIsPositive = amount > 0;
  return amountIsPositive;
}

function isValidAmount(number: string): boolean {
  const numberIsNotNaN = !isNaN(+number);
  return numberIsNotNaN;
}

function isAmountAnInteger(amount: number): boolean {
  return amount % 1 === 0;
}

function isBalanceHigherThanAmount(
  amount: number,
  balance: number,
): boolean {
  return balance > amount;
}
export const getAmountFormatError = (
  amount: string,
  isAmountInMAS: boolean = false,
) => {
  if (!isValidAmount(amount)) {
    return 'errors.send.invalid-amount';
  }
  if (!isPositiveAmount(+amount)) {
    return 'errors.send.amount-to-low';
  }
  const amountInNMAS = isAmountInMAS ? fromMAS(amount).toString() : amount;
  if (!isAmountAnInteger(+amountInNMAS)) {
    return 'errors.send.invalid-amount-decimals';
  }
};
export function getAmountTooHighError(
  amount: number,
  balance: number,
  isAmountInMAS: boolean = false,
): string | null {
  const amountInNMAS = isAmountInMAS ? fromMAS(amount).toString() : amount;
  if (!isBalanceHigherThanAmount(+amountInNMAS, balance)) {
    return 'errors.send.amount-to-high';
  }
  return null;
}
