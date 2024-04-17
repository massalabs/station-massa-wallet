import { useState, FormEvent } from 'react';

import {
  Button,
  Identicon,
  Money,
  MassaLogo,
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
  Selector,
  Clipboard,
} from '@massalabs/react-ui-kit';

import Intl from '@/i18n/i18n';
import { AccountObject } from '@/models/AccountModel';
import { parseForm } from '@/utils/';
import { formatAmount } from '@/utils/parseAmount';
import { SendInputsErrors } from '@/validation/sendInputs';

export type AmountValue = number | string | undefined;

interface MoneyForm {
  amount: AmountValue;
}

interface GenerateLinkProps {
  account: AccountObject;
  presetURL: string;
  url: string;
  setURL: (url: string) => void;
  setModal: (modal: boolean) => void;
}

function GenerateLink(props: GenerateLinkProps) {
  const { account, presetURL, setURL, setModal } = props;

  const recipient = account.nickname;
  const formattedBalance = formatAmount(
    account.candidateBalance,
  ).amountFormattedFull;

  const [amount, setAmount] = useState<number | string | undefined>('');
  const [link, setLink] = useState('');
  const [error, setError] = useState<SendInputsErrors | null>(null);

  function validate(formObject: MoneyForm) {
    const { amount } = formObject;

    setError(null);

    if (!amount) {
      setError({ amount: Intl.t('errors.send-coins.invalid-amount') });
      return false;
    }

    return true;
  }

  async function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const formObject = parseForm(e) as MoneyForm;

    if (!validate(formObject)) return;

    const amountArg = amount ? `&amount=${amount}` : '';

    const newURL = presetURL + amountArg;

    setURL(newURL);
    setLink(newURL);
  }

  return (
    <PopupModal
      customClass="!w-1/2 min-w-[775px]"
      fullMode={true}
      onClose={() => setModal(false)}
    >
      <PopupModalHeader>
        <h1 className="mas-banner mb-6">
          {Intl.t('receive-coins.generate-link')}
        </h1>
      </PopupModalHeader>
      <PopupModalContent>
        <form onSubmit={handleSubmit}>
          <div className="pb-10">
            <div className="flex flex-col gap-3 mb-6">
              <p className="mas-body2">
                {Intl.t('receive-coins.receive-amount')}
              </p>
              <Money
                data-testid="amount-to-send"
                placeholder={Intl.t('receive-coins.amount-to-ask')}
                name="amount"
                value={amount}
                onValueChange={(event) => setAmount(event.value)}
                error={error?.amount}
              />
            </div>
            <div className="flex flex-col gap-3 mb-6">
              <p className="mas-body2">{Intl.t('receive-coins.recipient')}</p>
              <Selector
                preIcon={<Identicon username={account.nickname} />}
                content={recipient}
                amount={formattedBalance}
                posIcon={<MassaLogo size={24} />}
                variant="secondary"
              />
            </div>
            <div className="flex flex-col gap-3 mb-3">
              <p className="mas-body2">
                {Intl.t('receive-coins.link-to-share')}
              </p>
              <div className="h-16">
                <Clipboard
                  data-testid="clipboard-link"
                  rawContent={link}
                  error={Intl.t('errors.no-content-to-copy')}
                />
              </div>
            </div>
            <div className="pb-3">
              <Button data-testid="generate-link-button" type="submit">
                {Intl.t('receive-coins.generate-link')}
              </Button>
            </div>
          </div>
        </form>
      </PopupModalContent>
    </PopupModal>
  );
}

export default GenerateLink;
