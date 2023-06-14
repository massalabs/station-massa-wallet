import { useState, FormEvent } from 'react';
import { formatStandard } from '../../../utils/massaFormat';
import { AccountObject } from '../../../models/AccountModel';
import Intl from '../../../i18n/i18n';
import { parseForm } from '../../../utils/parseForm';
import { FiCopy, FiCheckCircle } from 'react-icons/fi';
import {
  Button,
  Identicon,
  Input,
  Currency,
  MassaLogo,
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
  Selector,
} from '@massalabs/react-ui-kit';
import { SendInputsErrors } from '../../../validation/sendInputs';

interface GenerateLinkProps {
  account: AccountObject;
  presetURL: string;
  url: string;
  setURL: (url: string) => void;
  setModal: (modal: boolean) => void;
}

function GenerateLink(props: GenerateLinkProps) {
  const { account, presetURL, setURL, setModal } = props;

  const colorSuccess = '#1AE19D';
  const recipient = account.nickname;
  const recipientBalance = parseInt(account.candidateBalance) / 10 ** 9;
  const formattedBalance = formatStandard(recipientBalance);

  const [amount, setAmount] = useState<number | string | undefined>('');
  const [provider, setProvider] = useState('');
  const [link, setLink] = useState('');
  const [error, setError] = useState<SendInputsErrors | null>(null);
  const [success, setSuccess] = useState<boolean>(false);

  function validate(formObject: any) {
    const { amount } = formObject;

    setError(null);

    if (!amount) {
      setError({ amount: Intl.t('errors.send-coins.invalid-amount') });
      return false;
    }

    return true;
  }

  function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const formObject = parseForm(e);

    if (!validate(formObject)) return;

    const amountArg = amount ? `&amount=${amount}` : '';
    const providerArg = provider ? `&provider=${provider}` : '';

    const newURL = presetURL + amountArg + providerArg;

    setURL(newURL);
    setLink(newURL);
  }

  function handleCopyToClipboard() {
    setError(null);

    if (!link) {
      setError({ link: Intl.t('errors.receive-coins.no-link') });
      return;
    }

    navigator.clipboard.writeText(link);
    setSuccess(true);
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
            <div className="flex flex-col gap-3 mb-6">
              <p className="mas-body2">{Intl.t('receive-coins.provider')}</p>
              <Input
                placeholder={Intl.t('receive-coins.provider-description')}
                value={provider}
                name="provider"
                onChange={(e) => setProvider(e.target.value)}
                error={error?.address}
              />
            </div>
            <div className="flex flex-col gap-3 mb-3">
              <p className="mas-body2">
                {Intl.t('receive-coins.link-to-share')}{' '}
                {success && (
                  <label className="text-s-success">
                    {Intl.t('receive-coins.copied')}
                  </label>
                )}
              </p>
              <Input
                placeholder={Intl.t('receive-coins.link-to-share')}
                value={link}
                name="link"
                icon={
                  success ? <FiCheckCircle color={colorSuccess} /> : <FiCopy />
                }
                onClickIcon={handleCopyToClipboard}
                onChange={(e) => setLink(e.target.value)}
                error={error?.link}
              />
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
