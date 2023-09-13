import { useState, FormEvent } from 'react';

import {
  Button,
  Identicon,
  Currency,
  MassaLogo,
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
  Selector,
  Clipboard,
} from '@massalabs/react-ui-kit';
import axios from 'axios';

import Intl from '@/i18n/i18n';
import { AccountObject } from '@/models/AccountModel';
import { parseForm, formatStandard } from '@/utils/';
import { SendInputsErrors } from '@/validation/sendInputs';

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
  const recipientBalance = parseInt(account.candidateBalance) / 10 ** 9;
  const formattedBalance = formatStandard(recipientBalance);

  const [amount, setAmount] = useState<number | string | undefined>('');
  const [link, setLink] = useState('');
  const [error, setError] = useState<SendInputsErrors | null>(null);

  function validate(formObject: any) {
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
    const formObject = parseForm(e);

    if (!validate(formObject)) return;

    const amountArg = amount ? `&amount=${amount}` : '';

    const newURL = presetURL + amountArg;

    setURL(newURL);
    setLink(newURL);

    // TODO: This must be removed when the quest is finished
    // eslint-disable-next-line max-len
    const url = `https://dashboard.massa.net/quest_validation/register_quest/massastation/GENERATE_LINK/${account.address}`;

    await axios.post(url).catch((err) => {
      console.log('Error registering quest: ', err);
    });
  }

  return (
    <PopupModal
      customClass="!w-1/2 min-w-[775px]"
      fullMode={true}
      onClose={() => setModal(false)}
    >
      <PopupModalHeader>
        <h1 className="mas-banner mb-6">
          {Intl.t('receive-coins.receive-account')}
        </h1>
      </PopupModalHeader>
      <PopupModalContent>
        <form onSubmit={handleSubmit}>
          <div className="pb-10">
            <div className="flex flex-col gap-3 mb-6">
              <p className="mas-body2">
                {Intl.t('receive-coins.receive-amount')}
              </p>
              <Currency
                placeholder={Intl.t('receive-coins.amount-to-ask')}
                name="amount"
                value={amount}
                onValueChange={(value) => setAmount(value)}
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
                  rawContent={link}
                  toggleHover={false}
                  error={Intl.t('errors.no-content-to-copy')}
                />
              </div>
            </div>
            <div className="pb-3">
              <Button type="submit">
                {Intl.t('receive-coins.receive-account')}
              </Button>
            </div>
          </div>
        </form>
      </PopupModalContent>
    </PopupModal>
  );
}

export default GenerateLink;
