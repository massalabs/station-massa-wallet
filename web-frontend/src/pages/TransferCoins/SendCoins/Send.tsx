import { useResource } from '../../../custom/api/useResource';
import { AccountObject } from '../../../models/AccountModel';
import Intl from '../../../i18n/i18n';
import { useLocation } from 'react-router-dom';
import { parseForm } from '../../../utils/parseForm';
import { useState, ChangeEvent, FormEvent } from 'react';
import {
  formatStandard,
  reverseFormatStandard,
  checkRecipientFormat,
  Unit,
} from '../../../utils/MassaFormating';
import { SendForm } from './SendForm';
import { SendConfirmation } from './SendConfirmation';

interface IErrorObject {
  amount?: string;
  fees?: string;
  recipient?: string;
  error?: string;
}

function Send() {
  const { state } = useLocation();
  const nickname = state.nickname;
  const { data: accounts = [] } = useResource<AccountObject[]>('accounts');
  let account = accounts.find((account) => account.nickname === nickname);
  const unformattedBalance = account?.candidateBalance ?? '0';
  const balance = parseInt(unformattedBalance);

  const [amount, setAmount] = useState<string>('');
  const [fees, setFees] = useState<number>(1000);
  const [recipient, setRecipient] = useState('');
  const [valid, setValid] = useState<boolean>(false);
  const [error, setError] = useState<IErrorObject | null>(null);
  const [modal, setModal] = useState<boolean>(false);
  const [modalAccounts, setModalAccounts] = useState<boolean>(false);

  // TODO : implement address check to see if it is real

  function validate(e: any) {
    const formObject = parseForm(e);
    const { amount, recipient } = formObject;
    console.log(fees);
    const amountNum = reverseFormatStandard(amount);

    if (amountNum > balance) {
      setError({ amount: Intl.t('errors.send.amount-to-high') });
      return false;
    }
    if (amountNum <= 0) {
      setError({ amount: Intl.t('errors.send.amount-to-low') });
      return false;
    }
    if (Number.isNaN(amountNum)) {
      setError({ amount: Intl.t('errors.send.invalid-amount') });
      return false;
    }
    if (!recipient) {
      setError({ recipient: Intl.t('errors.send.no-recipient') });
      return false;
    }
    if (!checkRecipientFormat(recipient)) {
      setError({ recipient: Intl.t('errors.send.invalid-recipient') });
      return false;
    }
    setError(null);
    return true;
  }

  // This function is going to be used for the backend
  // Some parts are commented to avoid a build errors
  // because the amount and recipient aren't used yet
  function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (!validate(e)) return;
    // const formObject = parseForm(e);
    // const { amount, recipient } = formObject;
    else setValid(!valid);
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
    // if there is more than 9 decimals, set the amount to the previous value
    if (candidateAmount.split('.')[1]?.length > 9) {
      setAmount(amount);
      return;
    }
    setAmount(e.target.value);
  }

  function handleFeesConfirm() {
    setModal(!modal);
  }

  function handleFees(num: number) {
    console.log(num);
    setFees(num);
  }

  function handleModal() {
    setModal(!modal);
  }
  function handleModalAccounts() {
    setModalAccounts(!modalAccounts);
  }

  const confirmArgs = {
    amount,
    nickname,
    recipient,
    valid: valid,
    fees: fees,
    setValid,
  };

  const sendArgs = {
    amount,
    account,
    formattedBalance,
    recipient,
    error,
    fees,
    modal,
    modalAccounts,
    setModal,
    setModalAccounts,
    handleModal,
    handleModalAccounts,
    setFees,
    handleFees,
    handleFeesConfirm,
    setRecipient,
    handleSubmit,
    handleChange,
    SendPercentage,
  };

  return (
    <div className="mt-3.5">
      {valid ? (
        <SendConfirmation {...confirmArgs} />
      ) : (
        <SendForm {...sendArgs} />
      )}
    </div>
  );
}

export default Send;
