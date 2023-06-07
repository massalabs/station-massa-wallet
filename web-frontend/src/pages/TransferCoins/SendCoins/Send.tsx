import Intl from '../../../i18n/i18n';
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
import { useQuery } from '../../../custom/api/useQuery';
import { AccountObject } from '../../../models/AccountModel';

export interface SendProps {
  account: AccountObject;
}

function Send(props: SendProps) {
  const { account } = props;
  const nickname = account.nickname ?? '';
  let query = useQuery();

  let presetTo = query.get('to');
  let presetAmount = query.get('amount') ?? '';
  const unformattedBalance = account?.candidateBalance ?? '0';
  const balance = parseInt(unformattedBalance);

  const [amount, setAmount] = useState<string>(presetAmount);
  const [fees, setFees] = useState<number>(1000);
  const [recipient, setRecipient] = useState<string>(presetTo ?? '');
  const [valid, setValid] = useState<boolean>(false);
  const [error, setError] = useState<object | null>(null);
  const [modal, setModal] = useState<boolean>(false);
  const [modalAccounts, setModalAccounts] = useState<boolean>(false);

  function validate(e: FormEvent<HTMLFormElement>) {
    const formObject = parseForm(e);
    const { amount, recipient, fees } = formObject;
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
    if (Number.isNaN(fees)) {
      setError({ fees: Intl.t('errors.send.invalid-amount') });
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
  function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (!validate(e)) return;
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
    // MAX 9 DECIMALS
    if (candidateAmount.split('.')[1]?.length > 9) {
      setAmount(amount);
      return;
    }
    setAmount(e.target.value);
  }

  function handleFees(num: number) {
    setFees(num);
  }

  function handleModalAccounts() {
    setModalAccounts(!modalAccounts);
  }

  function handleConfirm() {
    setModal(!modal);
  }

  const confirmArgs = {
    amount,
    nickname,
    recipient,
    valid: valid,
    fees: fees,
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

export default Send;
