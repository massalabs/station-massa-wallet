import { useState, FormEvent, useEffect } from 'react';

import { fromMAS } from '@massalabs/massa-web3';
import {
  Button,
  Money,
  Input,
  formatAmount,
  Dropdown,
  IOption,
  getIcon,
  Spinner,
  parseAmount,
} from '@massalabs/react-ui-kit';
import { FiArrowUpRight, FiPlus } from 'react-icons/fi';
import { useParams } from 'react-router-dom';

import Advanced, { PRESET_LOW } from './Advanced';
import ContactList from './ContactList';
import { SendConfirmationData } from './SendConfirmation';
import { MAS } from '@/const/assets/assets';
import { useResource } from '@/custom/api';
import { useFTTransfer } from '@/custom/smart-contract/useFTTransfer';
import Intl from '@/i18n/i18n';
import { AccountObject } from '@/models/AccountModel';
import { Asset } from '@/models/AssetModel';
import {
  IForm,
  parseForm,
  useFetchAccounts,
  checkAddressFormat,
} from '@/utils';
import { handlePercent } from '@/utils/math';
import { symbolDict } from '@/utils/tokenIcon';

interface InputsErrors {
  amount?: string;
  recipient?: string;
}

interface SendFormProps {
  handleSubmit: (confirmed: SendConfirmationData) => void;
  account: AccountObject;
  data?: SendConfirmationData;
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
  const { amount: redirectAmount, to: redirectedTo } = redirect;
  const { nickname } = useParams();
  const { isMainnet } = useFTTransfer(nickname || '');

  const [error, setError] = useState<InputsErrors | null>(null);
  const [advancedModal, setAdvancedModal] = useState<boolean>(false);
  const [ContactListModal, setContactListModal] = useState<boolean>(false);
  const [amount, setAmount] = useState<string>(
    data && data.amount ? data.amount : '',
  );
  const [fees, setFees] = useState<string>(PRESET_LOW);
  const [recipient, setRecipient] = useState<string>(
    data?.recipientAddress || '',
  );
  const { okAccounts: accounts } = useFetchAccounts();
  const filteredAccounts = accounts?.filter(
    (account: AccountObject) => account?.nickname !== currentAccount?.nickname,
  );
  const [selectedAsset, setSelectedAsset] = useState<Asset | undefined>(
    data?.asset,
  );
  const { data: assets, isLoading: isAssetsLoading } = useResource<Asset[]>(
    `accounts/${nickname}/assets`,
    false,
  );

  const balance = BigInt(selectedAsset?.balance || '0');

  useEffect(() => {
    if (data) {
      setAmount(redirectAmount || data.amount);
      setRecipient(redirectedTo || data.recipientAddress);
      setFees(data.fees || PRESET_LOW);
    }
  }, [data, redirectAmount, redirectedTo]);

  function validate(formObject: IForm) {
    const { recipientAddress } = formObject;
    setError(null);
    if (!selectedAsset) return;

    if (!amount) {
      setError({ amount: Intl.t('errors.send-coins.invalid-amount') });
      return false;
    }

    const rawAmount = parseAmount(amount, selectedAsset.decimals);
    if (rawAmount <= 0n) {
      setError({ amount: Intl.t('errors.send-coins.amount-to-low') });
      return false;
    }

    if (rawAmount > balance) {
      setError({ amount: Intl.t('errors.send-coins.amount-to-high') });
      return false;
    }

    if (selectedAsset.symbol === MAS) {
      if (fromMAS(amount) + fromMAS(fees) > balance) {
        setError({
          amount: Intl.t('errors.send-coins.amount-plus-fees-to-high'),
        });
        return false;
      }
    }

    if (!recipientAddress) {
      setError({ recipient: Intl.t('errors.send-coins.no-address') });
      return false;
    }

    if (!checkAddressFormat(recipientAddress)) {
      setError({ recipient: Intl.t('errors.send-coins.invalid-address') });
      return false;
    }

    return true;
  }

  function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const formObject = parseForm(e);

    if (!validate(formObject) || !selectedAsset) return;

    sendCoinsHandleSubmit({
      ...(formObject as SendConfirmationData),
      amount,
      asset: selectedAsset,
      fees,
    });
  }

  useEffect(() => {
    if (!selectedAsset && assets && assets?.length > 0) {
      setSelectedAsset(assets?.[0]);
    }
  }, [assets, setSelectedAsset, selectedAsset]);

  const selectedAssetKey: number = selectedAsset
    ? assets?.indexOf(selectedAsset) || 0
    : 0;

  let options: IOption[] = [];

  if (assets) {
    options = assets.map((asset) => {
      const formattedBalance = formatAmount(
        asset.balance,
        asset.decimals,
      ).amountFormattedFull;
      return {
        itemPreview: asset.name,
        item: (
          <div>
            <p>{asset.name}</p>
            <p className="mas-caption">
              {Intl.t('send-coins.balance')} {formattedBalance}
            </p>
          </div>
        ),
        icon: getIcon(
          symbolDict[asset.symbol as keyof typeof symbolDict],
          true,
          isMainnet,
          28,
        ),
        onClick: () => setSelectedAsset(asset),
      };
    });
  }

  const formattedBalance = selectedAsset?.balance ? (
    <u>
      {
        formatAmount(selectedAsset.balance, selectedAsset.decimals)
          .amountFormattedFull
      }
    </u>
  ) : (
    <Spinner size={12} customClass="inline-block" />
  );

  return (
    <div>
      <form onSubmit={handleSubmit}>
        <div className="flex flex-row justify-between w-full pb-3.5 ">
          <p className="mas-body2"> {Intl.t('send-coins.send-action')} </p>
        </div>
        <div className="flex flex-row gap-2">
          <div className="grow">
            <div className="pb-3">
              <Money
                placeholder={Intl.t('send-coins.amount-to-send')}
                name="amount"
                value={amount}
                suffix=""
                onValueChange={(event) => setAmount(event.value)}
                error={error?.amount}
                decimalScale={selectedAsset?.decimals}
                customClass="h-14"
              />
            </div>
            <div className="flex flex-row-reverse">
              <ul className="flex flex-row mas-body2">
                <li
                  data-testid="send-percent-25"
                  onClick={() =>
                    setAmount(
                      handlePercent(
                        balance,
                        25n,
                        fromMAS(fees),
                        balance,
                        selectedAsset?.decimals || 9,
                        selectedAsset?.symbol || '',
                      ),
                    )
                  }
                  className="mr-3.5 hover:cursor-pointer"
                >
                  25%
                </li>
                <li
                  data-testid="send-percent-50"
                  onClick={() =>
                    setAmount(
                      handlePercent(
                        balance,
                        50n,
                        fromMAS(fees),
                        balance,
                        selectedAsset?.decimals || 9,
                        selectedAsset?.symbol || '',
                      ),
                    )
                  }
                  className="mr-3.5 hover:cursor-pointer"
                >
                  50%
                </li>
                <li
                  data-testid="send-percent-75"
                  onClick={() =>
                    setAmount(
                      handlePercent(
                        balance,
                        75n,
                        fromMAS(fees),
                        balance,
                        selectedAsset?.decimals || 9,
                        selectedAsset?.symbol || '',
                      ),
                    )
                  }
                  className="mr-3.5 hover:cursor-pointer"
                >
                  75%
                </li>
                <li
                  data-testid="send-percent-100"
                  onClick={() =>
                    setAmount(
                      handlePercent(
                        balance,
                        100n,
                        fromMAS(fees),
                        balance,
                        selectedAsset?.decimals || 9,
                        selectedAsset?.symbol || '',
                      ),
                    )
                  }
                  className="mr-3.5 hover:cursor-pointer"
                >
                  Max
                </li>
              </ul>
            </div>
          </div>
          <div>
            <Dropdown
              select={selectedAssetKey}
              readOnly={isAssetsLoading}
              size="md"
              options={options}
              className="pb-3.5"
              fullWidth={false}
            />
            <div>
              <p className="mas-caption inline-block mr-1">
                {Intl.t('send-coins.available-balance')}
              </p>
              {formattedBalance}
            </div>
          </div>
        </div>
        <p className="pb-3.5 mas-body2">{Intl.t('send-coins.recipient')}</p>
        <div className="pb-3.5">
          <Input
            placeholder={Intl.t('receive-coins.recipient')}
            value={recipient}
            name="recipientAddress"
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
