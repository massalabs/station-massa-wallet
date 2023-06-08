import Intl from '../i18n/i18n';

/**
 * Validate the recipient and/or amount in a form
 * @param recipient - The recipient address to validate.
 * @param amount - The amount to validate.
 * @param balance - The balance to validate.
 * @returns An object containing the error object if there is an error, else null
 */

import {
  checkAddressFormat,
  reverseFormatStandard,
} from '../utils/MassaFormating';

export function validateInputs(
  amount: string,
  recipient: string,
  addressType = 'recipient',
  balance?: string,
): object | null {
  let errorsAmount = null;
  let errorsRecipient = null;
  errorsAmount = validateAmount(amount, balance);
  errorsRecipient = validateAddress(recipient, addressType);

  return errorsAmount || errorsRecipient;
}

/**
 * Validates an amount to send
 * @param amount - The amount to validate.
 * @param balance - The balance to validate.
 * @returns An object containing the error message if the amount is invalid, `null` otherwise.
 * @see validateInputs
 */
export function validateAmount(
  amount: string,
  balance?: string,
  amountType = 'amount',
): object | null {
  const amountNum = reverseFormatStandard(amount) * 10 ** 9;
  if (amountNum <= 0) {
    return {
      amount: Intl.t('errors.send.amount-to-low', { type: amountType }),
    };
  }
  if (Number.isNaN(amountNum)) {
    return {
      amount: Intl.t('errors.send.invalid-amount', { type: amountType }),
    };
  }

  console.log(balance);
  if (!balance) return null;
  console.log(balance);
  if (+balance === undefined)
    return { unexpectedError: Intl.t('errors.unexpected-error.description') };
  if (balance?.length > 0 && amountNum > +balance) {
    return { amount: Intl.t('errors.send.amount-to-high') };
  }
  return null;
}

/**
 * Validates a recipient address.
 * @param recipient - The recipient address to validate.
 * @returns An object containing the error message if the address is invalid, `null` otherwise.
 * @remarks This function does not check if the address is valid on the blockchain.
 * @see checkAddressFormat
 * @see validateInputs
 */

function validateAddress(
  recipient: string,
  addresstype = 'recipient',
): object | null {
  if (!recipient) {
    return {
      address: Intl.t('errors.send.no-address', { type: addresstype }),
    };
  }
  if (!checkAddressFormat(recipient)) {
    return {
      address: Intl.t('errors.send.invalid-address', { type: addresstype }),
    };
  }
  return null;
}
