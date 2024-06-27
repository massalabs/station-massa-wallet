import { useState, FormEvent } from 'react';

import {
  Button,
  Money,
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
  Clipboard,
  formatAmount,
  Selector,
  Identicon,
} from '@massalabs/react-ui-kit';

import Intl from '@/i18n/i18n';
import { AccountObject } from '@/models/AccountModel';
import { Asset } from '@/models/AssetModel';
import { AssetSelector } from '@/pages/TransferCoins/SendCoins/AssetSelector';
import { parseForm } from '@/utils/';
import { tokenIcon } from '@/utils/tokenIcon';
import { SendInputsErrors } from '@/validation/sendInputs';

interface MoneyForm {
  amount: string;
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

  const [amount, setAmount] = useState<string>('');
  const [link, setLink] = useState('');
  const [error, setError] = useState<SendInputsErrors | null>(null);
  const [selectedAsset, setSelectedAsset] = useState<Asset | undefined>();

  const formattedBalance = selectedAsset
    ? formatAmount(selectedAsset.balance || '', selectedAsset.decimals).full
    : '';

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
    if (!selectedAsset) return;

    const amountArg = amount
      ? `&amount=${amount}&symbol=${selectedAsset.symbol}`
      : '';

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
              <p className="mas-body2">{Intl.t('receive-coins.coins')}</p>
              <div className="flex flex-row gap-2">
                <div className="grow">
                  <Money
                    data-testid="amount-to-send"
                    placeholder={Intl.t('receive-coins.amount-to-ask')}
                    name="amount"
                    value={amount}
                    suffix=""
                    onValueChange={(event) => setAmount(event.value)}
                    error={error?.amount}
                    decimalScale={selectedAsset?.decimals}
                    customClass="h-14 pb-3"
                  />
                </div>
                <AssetSelector
                  selectedAsset={selectedAsset}
                  setSelectedAsset={setSelectedAsset}
                />
              </div>
            </div>
            <div className="flex flex-col gap-3 mb-6">
              <p className="mas-body2">{Intl.t('receive-coins.recipient')}</p>
              <Selector
                preIcon={<Identicon username={account.nickname} />}
                content={account.nickname}
                amount={formattedBalance}
                posIcon={tokenIcon(selectedAsset?.symbol || '', 24)}
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
              <Button
                data-testid="generate-link-button"
                type="submit"
                disabled={!selectedAsset}
              >
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
