import { useState, ChangeEvent, FormEvent } from 'react';

import {
  validateInputs,
  SendInputsErrors,
  validateAmount,
} from '../../../validation/sendInputs';
import { formatStandard, Unit } from '../../../utils/MassaFormating';
import { parseForm } from '../../../utils/parseForm';

import { useQuery } from '../../../custom/api/useQuery';
import { AccountObject } from '../../../models/AccountModel';

import { SendForm } from './SendForm';
import { SendConfirmation } from './SendConfirmation';

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

  function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const formObject = parseForm(e);
    const { amount, recipient } = formObject;
    const errors = validateInputs(
      amount,
      recipient,
      'Recipient',
      unformattedBalance,
    );
    setError(errors);
    if (errors !== null) return;
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

  function handleValidation(errors: object | null) {
    if (errors !== null) {
      return;
    } else {
      setModal(false);
    }
  }

  // we set the fees using a useState so they can render in the modal
  // We pass the validateAmount Fn as args directly ta handle the confirmation
  // Because the useState is too slow to be use as a param in th Fn
  function handleConfirm() {
    setErrorAdvanced(
      validateAmount(fees ?? '0.00', unformattedBalance, 'Fees'),
    );
    handleValidation(
      validateAmount(fees ?? '0.00', unformattedBalance, 'Fees'),
    );
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
