import Intl from '../i18n/i18n';
import {
  checkAddressFormat,
  reverseFormatStandard,
} from '../utils/MassaFormating';

export interface SendInputsErrors {
  unexpectedError?: string;
  amount?: string;
  address?: string;
}

/**
 * Validate the recipient and/or amount in a form
 * @param recipient - The recipient address to validate.
 * @param amount - The amount to validate.
 * @param balance - The balance to validate.
 * @returns An object containing the error object if there is an error, else null
 */
export function validateInputs(
  amount: string,
  address: string,
  addressType = 'recipient',
  balance?: string,
): SendInputsErrors | null {
  let errorsAmount = null;
  let errorsRecipient = null;
  errorsAmount = validateAmount(amount, balance);
  if (addressType === 'provider' && !address) return errorsAmount; // provider can be empty
  errorsRecipient = validateAddress(address, addressType);

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
  amountType = 'Amount',
): SendInputsErrors | null {
  let amountNum = reverseFormatStandard(amount);
  let verb = 'are';
  if (amountType == 'Amount') {
    amountNum *= 10 ** 9;
    verb = 'in';
  }
  if (amountNum <= 0) {
    return {
      amount: Intl.t('errors.send.amount-to-low', { type: amountType, verb }),
    };
  }
  if (Number.isNaN(amountNum)) {
    return {
      amount: Intl.t('errors.send.invalid-amount', { type: amountType, verb }),
    };
  }

  if (!balance) return null;
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
): SendInputsErrors | null {
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
