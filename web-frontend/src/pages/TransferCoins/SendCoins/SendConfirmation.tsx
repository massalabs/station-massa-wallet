import { useEffect } from 'react';

import {
  Balance,
  Button,
  Tooltip,
  formatFTAmount,
  getAssetIcons,
  parseAmount,
  Clipboard,
} from '@massalabs/react-ui-kit';
import { FiChevronLeft } from 'react-icons/fi';

import { PRESET_HIGH, PRESET_LOW, PRESET_STANDARD } from './Advanced';
import { useFTTransfer } from '@/custom/smart-contract/useFTTransfer';
import { useMNS } from '@/custom/useMNS';
import Intl from '@/i18n/i18n';
import { Asset } from '@/models/AssetModel';
import { maskAddress } from '@/utils';
import { symbolDict } from '@/utils/tokenIcon';

export interface SendConfirmationData {
  amount: string;
  asset: Asset;
  fees: string;
  recipientAddress: string;
}

interface SendConfirmationProps {
  data: SendConfirmationData;
  handleConfirm: (confirm: boolean) => void;
  isLoading: boolean;
}

export function SendConfirmation(props: SendConfirmationProps) {
  const { data, handleConfirm, isLoading } = props;

  const { isMainnet } = useFTTransfer();

  const { amount, asset, fees, recipientAddress } = data;
  const { symbol, decimals } = asset;

  const FEES_STANDARD = Intl.t('send-coins.fee-standard');
  const FEES_LOW = Intl.t('send-coins.fee-low');
  const FEES_HIGH = Intl.t('send-coins.fee-high');
  const FEES_CUSTOM = Intl.t('send-coins.fee-custom');

  const formattedRecipientAddress = maskAddress(recipientAddress);
  // amount is the value given by the Money input component
  // we convert to the smallest unit with parseAmount
  // and then format it with formatFTAmount
  const formattedAmount = formatFTAmount(
    parseAmount(amount, data.asset.decimals),
    decimals,
  ).amountFormattedFull;

  const { reverseResolveDns, mns } = useMNS();

  useEffect(() => {
    reverseResolveDns();
  }, [reverseResolveDns]);

  let selectedFees;

  switch (fees) {
    case PRESET_STANDARD:
      selectedFees = FEES_STANDARD;
      break;
    case PRESET_LOW:
      selectedFees = FEES_LOW;
      break;
    case PRESET_HIGH:
      selectedFees = FEES_HIGH;
      break;
    default:
      selectedFees = FEES_CUSTOM;
      break;
  }

  const feeInfo = Intl.t('send-coins.fee-info', {
    feeType: selectedFees,
    fee: fees,
  });
  const accountCreationFeeAlert = `  \u26A0  ${Intl.t(
    'send-coins.account-fee-alert',
  )}`;
  let content =
    selectedFees == FEES_LOW
      ? feeInfo + accountCreationFeeAlert
      : accountCreationFeeAlert;

  return (
    <>
      <div
        onClick={() => handleConfirm(false)}
        className="flex flex-row just items-center hover:cursor-pointer mb-5 gap-2"
      >
        <FiChevronLeft />
        <p>{Intl.t('send-coins.back-to-sending')}</p>
      </div>
      <div
        data-testid="send-confirmation"
        className="flex flex-col items-center gap-4 p-10 bg-secondary rounded-lg mb-6"
      >
        <p className="mb-6 mas-body">{Intl.t('send-coins.send-message')}</p>

        <Balance
          customClass="p-0 bg-transparent"
          amount={formattedAmount}
          symbol={symbol}
          icon={getAssetIcons(
            symbolDict[symbol as keyof typeof symbolDict],
            true,
            isMainnet,
            32,
            'mr-3',
          )}
        />
        <div className="flex items-center gap-8 ">
          <div className="flex items-center">
            <Tooltip body={content} />
            <p>{Intl.t('send-coins.fee')}</p>
          </div>
          <Balance amount={fees} symbol="MAS" size="xs" />
        </div>
        <div
          data-testid="send-confirmation-recipient"
          className="text-info flex items-center gap-2"
        >
          <div>{Intl.t('send-coins.recipient')}</div>
          <Clipboard
            displayedContent={formattedRecipientAddress}
            rawContent={recipientAddress}
            error={Intl.t('errors.no-content-to-copy')}
            className="flex flex-row items-center mas-body2 justify-between
              w-fit h-fit px-3 rounded bg-primary cursor-pointer"
          />
        </div>
        {mns && (
          <div className="flex items-center gap-2 text-info">
            <div>{Intl.t('send-coins.mns.mns')}</div>
            <div>{mns}</div>
          </div>
        )}
      </div>
      <Button disabled={isLoading} onClick={() => handleConfirm(true)}>
        {Intl.t('send-coins.confirm-sign')}
      </Button>
    </>
  );
}
