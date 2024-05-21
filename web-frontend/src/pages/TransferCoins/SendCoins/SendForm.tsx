import {
  useState,
  FormEvent,
  useEffect,
  useCallback,
  ChangeEvent,
} from 'react';

import { fromMAS } from '@massalabs/massa-web3';
import {
  Button,
  Money,
  Input,
  formatAmount,
  Spinner,
  parseAmount,
} from '@massalabs/react-ui-kit';
import { FiArrowUpRight, FiPlus } from 'react-icons/fi';

import { MAS } from '@/const/assets/assets';
import { useMNS } from '@/custom/useMNS';
import Intl from '@/i18n/i18n';
import { AccountObject } from '@/models/AccountModel';
import { Asset } from '@/models/AssetModel';
import Advanced, { PRESET_LOW } from '@/pages/TransferCoins/SendCoins/Advanced';
import { AssetSelector } from '@/pages/TransferCoins/SendCoins/AssetSelector';
import ContactList from '@/pages/TransferCoins/SendCoins/ContactList';
import { SendConfirmationData } from '@/pages/TransferCoins/SendCoins/SendConfirmation';
import { Redirect } from '@/pages/TransferCoins/TransferCoins';
import {
  IForm,
  parseForm,
  useFetchAccounts,
  checkAddressFormat,
} from '@/utils';
import { mnsExtension } from '@/utils/const';
import { handlePercent } from '@/utils/math';

interface InputsErrors {
  amount?: string;
  recipient?: string;
}

interface SendFormProps {
  handleSubmit: (confirmed: SendConfirmationData) => void;
  account: AccountObject;
  data?: SendConfirmationData;
  redirect: Redirect;
}

export function SendForm(props: SendFormProps) {
  const {
    handleSubmit: sendCoinsHandleSubmit,
    account: currentAccount,
    data: sendOpData,
    redirect,
  } = props;
  const {
    amount: redirectAmount,
    to: redirectedTo,
    symbol: redirectSymbol,
  } = redirect;

  const { okAccounts: accounts } = useFetchAccounts();
  const { resolveDns, targetMnsAddress, resetTargetMnsAddress } = useMNS();

  const [amount, setAmount] = useState<string>('');
  const [fees, setFees] = useState<string>(PRESET_LOW);
  const [recipient, setRecipient] = useState<string>('');
  const [selectedAsset, setSelectedAsset] = useState<Asset | undefined>();

  const [error, setError] = useState<InputsErrors | null>(null);
  const [advancedModal, setAdvancedModal] = useState<boolean>(false);
  const [ContactListModal, setContactListModal] = useState<boolean>(false);
  const [mnsAddressCorrelation, setMnsAddressCorrelation] = useState<boolean>();
  const filteredAccounts = accounts?.filter(
    (account: AccountObject) => account?.nickname !== currentAccount?.nickname,
  );

  const balance = BigInt(selectedAsset?.balance || '0');

  useEffect(() => {
    if (!sendOpData) {
      setAmount(redirectAmount);
      setFees(PRESET_LOW);
      setRecipient(redirectedTo);
    } else {
      if (sendOpData.recipientDomainName) {
        setRecipient(sendOpData.recipientDomainName + mnsExtension);
        resolveDns(sendOpData.recipientDomainName);
      } else {
        setRecipient(sendOpData.recipientAddress);
      }
      setAmount(sendOpData.amount);
      setFees(sendOpData.fees);
      setSelectedAsset(sendOpData.asset);
    }
  }, [sendOpData, redirectAmount, redirectedTo, resolveDns]);

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
    let formObject = parseForm(e);

    if (targetMnsAddress && /\.massa$/.test(recipient)) {
      // recipient is a domain name
      formObject = {
        ...formObject,
        recipientAddress: targetMnsAddress, // set the resolved address
        recipientDomainName: recipient.replace(mnsExtension, ''), // set the domain name
      };
    } else {
      setMnsAddressCorrelation(false);
      resetTargetMnsAddress();
    }

    if (!validate(formObject) || !selectedAsset) return;

    sendCoinsHandleSubmit({
      ...(formObject as SendConfirmationData),
      amount,
      asset: selectedAsset,
      fees,
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

  const handleRecipientChange = useCallback(
    async (e: ChangeEvent<HTMLInputElement>) => {
      const mnsExtensionRegex = /\.massa$/;
      const value = e.target.value;
      setRecipient(value);
      setError((error) => ({ ...error, recipient: '' }));
      if (value !== '' && mnsExtensionRegex.test(value)) {
        const inputMns = value.replace(mnsExtension, '');
        const targetAddress = await resolveDns(inputMns);
        if (targetAddress && checkAddressFormat(targetAddress)) {
          setMnsAddressCorrelation(true);
        } else {
          setMnsAddressCorrelation(false);
          resetTargetMnsAddress();
        }
      }
    },
    [resolveDns, resetTargetMnsAddress],
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
            <AssetSelector
              selectedAsset={selectedAsset || sendOpData?.asset}
              setSelectedAsset={setSelectedAsset}
              selectSymbol={redirectSymbol}
            />
            <div data-testid="available-amount">
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
            onChange={(e) => handleRecipientChange(e)}
            error={error?.recipient}
            success={
              mnsAddressCorrelation
                ? Intl.t('send-coins.mns.mns-correlation', {
                    mns: recipient,
                    address: targetMnsAddress,
                  })
                : ''
            }
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
