import { useState, FormEvent, useEffect } from 'react';

import { Balance, Button, Money, Input } from '@massalabs/react-ui-kit';
import { FiArrowUpRight, FiPlus } from 'react-icons/fi';

import Advanced from './Advanced';
import ContactList from './ContactList';
import { SendConfirmationData } from './SendConfirmation';
import Intl from '@/i18n/i18n';
import { AccountObject } from '@/models/AccountModel';
import {
  IForm,
  parseForm,
  toNanoMASS,
  toMASS,
  formatStandard,
  reverseFormatStandard,
  fetchAccounts,
  checkAddressFormat,
} from '@/utils';
import { handlePercent } from '@/utils/math';

interface InputsErrors {
  amount?: string;
  recipient?: string;
}

interface SendFormProps {
  handleSubmit: (confirmed: SendConfirmationData) => void;
  account: AccountObject;
  data: SendConfirmationData;
  redirect: {
    amount: string;
    to: string;
  };
}

export function SendForm(props: SendFormProps) {
  const {
    handleSubmit: sendCoinsHandleSubmit,
    account: currentAccount,
    data,
    redirect,
  } = props;
  const { amount: redirectAmount, to } = redirect;

  const balance = BigInt(currentAccount.candidateBalance) || 0n;
  const formattedBalance = formatStandard(toMASS(balance));

  const [error, setError] = useState<InputsErrors | null>(null);
  const [advancedModal, setAdvancedModal] = useState<boolean>(false);
  const [ContactListModal, setContactListModal] = useState<boolean>(false);
  const [amount, setAmount] = useState<string>(
    data.amount ? toMASS(toNanoMASS(data.amount)).toString() : '',
  );
  const [fees, setFees] = useState<bigint>(1000n);
  const [recipient, setRecipient] = useState<string>(data?.recipient ?? '');
  const { okAccounts: accounts } = fetchAccounts();
  const filteredAccounts = accounts?.filter(
    (account: AccountObject) => account?.nickname !== currentAccount?.nickname,
  );
  useEffect(() => {
    setAmount(redirectAmount);
    setRecipient(to);
  }, []);

  function validate(formObject: IForm) {
    const { amount, recipient } = formObject;

    setError(null);

    if (!amount) {
      setError({ amount: Intl.t('errors.send-coins.invalid-amount') });
      return false;
    }

    if (reverseFormatStandard(amount) <= 0) {
      setError({ amount: Intl.t('errors.send-coins.amount-to-low') });
      return false;
    }

    if (toMASS(toNanoMASS(amount)) > toMASS(balance)) {
      setError({ amount: Intl.t('errors.send-coins.amount-to-high') });
      return false;
    }

    if (toNanoMASS(amount) + fees > balance) {
      setError({
        amount: Intl.t('errors.send-coins.amount-plus-fees-to-high'),
      });
      return false;
    }

    if (!recipient) {
      setError({ recipient: Intl.t('errors.send-coins.no-address') });
      return false;
    }

    // It needs starts with AU and after AU have at least 4 chars
    if (!checkAddressFormat(recipient)) {
      setError({ recipient: Intl.t('errors.send-coins.invalid-address') });
      return false;
    }

    return true;
  }

  function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const formObject = parseForm(e);

    if (!validate(formObject)) return;

    sendCoinsHandleSubmit({
      ...(formObject as SendConfirmationData),
      fees: fees.toString(),
    });
  }

  return (
    <div>
      <form onSubmit={handleSubmit}>
        <p className="mas-subtitle mb-5">
          {Intl.t('send-coins.account-balance')}
        </p>
        <Balance customClass="mb-5" amount={formattedBalance} />
        <div className="flex flex-row justify-between w-full pb-3.5 ">
          <p className="mas-body2"> {Intl.t('send-coins.send-action')} </p>
          <p className="mas-body2">
            {Intl.t('send-coins.available-balance')} <u>{formattedBalance}</u>
          </p>
        </div>
        <div className="pb-3.5">
          <Money
            placeholder={Intl.t('send-coins.amount-to-send')}
            name="amount"
            value={amount}
            onValueChange={(event) => setAmount(event.value)}
            error={error?.amount}
          />
        </div>
        <div className="flex flex-row-reverse">
          <ul className="flex flex-row mas-body2">
            <li
              data-testid="send-percent-25"
              onClick={() =>
                setAmount(handlePercent(balance, 25n, fees, balance).toString())
              }
              className="mr-3.5 hover:cursor-pointer"
            >
              25%
            </li>
            <li
              data-testid="send-percent-50"
              onClick={() =>
                setAmount(handlePercent(balance, 50n, fees, balance).toString())
              }
              className="mr-3.5 hover:cursor-pointer"
            >
              50%
            </li>
            <li
              data-testid="send-percent-75"
              onClick={() =>
                setAmount(handlePercent(balance, 75n, fees, balance).toString())
              }
              className="mr-3.5 hover:cursor-pointer"
            >
              75%
            </li>
            <li
              data-testid="send-percent-100"
              onClick={() =>
                setAmount(
                  handlePercent(balance, 100n, fees, balance).toString(),
                )
              }
              className="mr-3.5 hover:cursor-pointer"
            >
              Max
            </li>
          </ul>
        </div>
        <p className="pb-3.5 mas-body2">{Intl.t('send-coins.recipient')}</p>
        <div className="pb-3.5">
          <Input
            placeholder={Intl.t('receive-coins.recipient')}
            value={recipient}
            name="recipient"
            onChange={(e) => setRecipient(e.target.value)}
            error={error?.recipient}
          />
        </div>
        {filteredAccounts.length > 0 && (
          <div className="flex flex-row-reverse pb-3.5">
            <p className="hover:cursor-pointer">
              <u
                data-testid="transfer-between-accounts"
                className="mas-body2"
                onClick={() => setContactListModal(true)}
              >
                {Intl.t('send-coins.transfer-between-acc')}
              </u>
            </p>
          </div>
        )}
        <div className="flex flex-col w-full gap-3.5">
          <Button
            onClick={() => setAdvancedModal(!advancedModal)}
            variant="secondary"
            posIcon={<FiPlus />}
          >
            {Intl.t('send-coins.advanced')}
          </Button>

          <div>
            <Button type="submit" posIcon={<FiArrowUpRight />}>
              {Intl.t('send-coins.send')}
            </Button>
          </div>
        </div>
      </form>
      {advancedModal && (
        <Advanced
          setFees={setFees}
          fees={fees}
          onClose={() => setAdvancedModal(false)}
        />
      )}
      {ContactListModal && (
        <ContactList
          setRecipient={setRecipient}
          accounts={filteredAccounts}
          onClose={() => setContactListModal(false)}
        />
      )}
    </div>
  );
}
