import { useState, ChangeEvent, FormEvent } from 'react';

import {
  SendInputsErrors,
  getAddressError,
  getAmountFormatError,
  getAmountTooHighError,
} from '../../../validation/sendInputs';
import { formatStandard, Unit } from '../../../utils/MassaFormating';
import { parseForm } from '../../../utils/parseForm';

import { useQuery } from '../../../custom/api/useQuery';
import { AccountObject } from '../../../models/AccountModel';

import { SendForm } from './SendForm';
import { SendConfirmation } from './SendConfirmation';
import Intl from '../../../i18n/i18n';

export interface SendCoinsProps {
  account: AccountObject;
}

function SendCoins(props: SendCoinsProps) {
  let query = useQuery();

  let presetTo = query.get('to');
  let presetAmount = query.get('amount') ?? '';

  const [amount, setAmount] = useState<string>(presetAmount);
  const [fees, setFees] = useState<string>('1000');
  const [recipient, setRecipient] = useState<string>(presetTo ?? '');
  const [valid, setValid] = useState<boolean>(false);
  const [error, setError] = useState<SendInputsErrors | null>(null);
  const [modal, setModal] = useState<boolean>(false);
  const [modalAccounts, setModalAccounts] = useState<boolean>(false);
  const [errorAdvanced, setErrorAdvanced] = useState<object | null>(null);

  const { account } = props;
  const nickname = account.nickname ?? '';

  const unformattedBalance = account?.candidateBalance ?? '0.00';
  const balance = parseInt(unformattedBalance);

  function validateSendForm(amount: string, recipient: string) {
    const errorAmountFormat = getAmountFormatError(amount, true);
    if (errorAmountFormat) {
      setError({
        amount: Intl.t(errorAmountFormat, { type: 'Amount', verb: 'is' }),
      });
      return false;
    }
    const errorAmountTooHigh = getAmountTooHighError(+amount, balance, true);
    if (errorAmountTooHigh) {
      setError({
        amount: Intl.t(errorAmountTooHigh, { type: 'Amount', verb: 'is' }),
      });
      return false;
    }
    const addressError = getAddressError(recipient);
    if (addressError) {
      setError({ address: Intl.t(addressError) });
      return false;
    }
    return true;
  }

  function validateAdvanced(fees: string) {
    const errorAmountFormat = getAmountFormatError(fees);
    if (errorAmountFormat) {
      setErrorAdvanced({
        fees: Intl.t(errorAmountFormat, { type: 'Fees', verb: 'are' }),
      });

      const errorAmountTooHigh = getAmountTooHighError(+fees, balance);
      if (errorAmountTooHigh) {
        setErrorAdvanced({
          fees: Intl.t(errorAmountTooHigh, { type: 'Fees', verb: 'are' }),
        });
      }
      return false;
    }
  }

  function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const formObject = parseForm(e);
    const { amount, recipient } = formObject;
    const errors = validateSendForm(amount, recipient);

    if (!errors) return;
    setValid(!valid);
  }

  const formattedBalance = formatStandard(balance, Unit.NanoMAS, 2);
  function overrideAmount(pct: number) {
    const newAmount = balance * pct;
    const newFormatedAmount = formatStandard(newAmount, Unit.NanoMAS, 9);
    setAmount(newFormatedAmount);
  }

  function SendPercentage({ percentage }: { percentage: number }) {
    return (
      <li
        onClick={() => {
          overrideAmount(percentage / 100);
        }}
        className="mr-3.5 hover:cursor-pointer"
      >
        {percentage === 100 ? 'Max' : `${percentage}%`}
      </li>
    );
  }

  function handleChange(e: ChangeEvent<HTMLInputElement>): void {
    e.preventDefault();
    const candidateAmount = e.target.value;
    // MAX 9 DECIMALS
    if (candidateAmount.split('.')[1]?.length > 9) {
      setAmount(amount);
      return;
    }
    setAmount(e.target.value);
  }

  // we set the fees using a useState so they can render in the modal
  // We pass the validateAmount Fn as args directly ta handle the confirmation
  // Because the useState is too slow to be use as a param in th Fn
  function handleConfirm() {
    const errors = validateAdvanced(fees);
    if (!errors) {
      return;
    } else {
      setModal(false);
    }
  }

  function handleFees(num: string) {
    // check if '.' is present
    if (num.includes('.')) return;
    setFees(num);
  }
  function handleModalAccounts() {
    setModalAccounts(!modalAccounts);
  }

  const confirmArgs = {
    amount,
    nickname,
    recipient,
    valid,
    fees,
    setValid,
    modal,
    setModal,
  };

  const sendArgs = {
    amount,
    account,
    formattedBalance,
    recipient,
    error,
    setErrorAdvanced,
    errorAdvanced,
    fees,
    modal,
    modalAccounts,
    setModal,
    setModalAccounts,
    handleModalAccounts,
    setFees,
    handleFees,
    handleConfirm,
    setRecipient,
    handleSubmit,
    handleChange,
    SendPercentage,
  };

  return (
    <>
      {valid ? (
        <div className="mt-5">
          <SendConfirmation {...confirmArgs} />
        </div>
      ) : (
        <div className="mt-5">
          <SendForm {...sendArgs} />
        </div>
      )}
    </>
  );
}

export default SendCoins;
