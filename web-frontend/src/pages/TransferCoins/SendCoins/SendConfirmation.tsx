import { useState } from 'react';

import { fromMAS } from '@massalabs/massa-web3';
import { Balance, Button } from '@massalabs/react-ui-kit';
import { FiChevronLeft } from 'react-icons/fi';

import { PRESET_HIGH, PRESET_LOW, PRESET_STANDARD } from './Advanced';
import ToolTip from './ToolTip';
import Intl from '@/i18n/i18n';
import { maskAddress } from '@/utils';
import { formatAmount } from '@/utils/parseAmount';

export interface SendConfirmationData {
  amount: string;
  fee: string;
  recipientAddress: string;
}

interface SendConfirmationProps {
  data: SendConfirmationData;
  handleConfirm: (confirm: boolean) => void;
  isLoading: boolean;
}

export function SendConfirmation(props: SendConfirmationProps) {
  const { data, handleConfirm, isLoading } = props;
  const { amount, fee, recipientAddress } = data;

  const FEES_STANDARD = Intl.t('send-coins.fee-standard');
  const FEES_LOW = Intl.t('send-coins.fee-low');
  const FEES_HIGH = Intl.t('send-coins.fee-high');
  const FEES_CUSTOM = Intl.t('send-coins.fee-custom');

  const formattedRecipientAddress = maskAddress(recipientAddress);
  const total = fromMAS(amount) + fromMAS(fee);
  const formattedAmount = formatAmount(
    fromMAS(amount).toString(),
  ).amountFormattedFull;

  const formattedTotal = formatAmount(total.toString()).amountFormattedFull;
  const [showTooltip, setShowTooltip] = useState(false);
  let selectedFees;

  switch (fee) {
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
    fee,
  });
  const gasAlert = `  \u26A0  ${Intl.t('send-coins.fee-alert')}`;
  let content = selectedFees == FEES_LOW ? feeInfo + gasAlert : gasAlert;

  return (
    <>
      <div
        onClick={() => handleConfirm(false)}
        className="flex flex-row just items-center hover:cursor-pointer mb-5 gap-2"
      >
        <FiChevronLeft />
        <p>{Intl.t('send-coins.back-to-sending')}</p>
      </div>
      <p className="mb-6">{Intl.t('send-coins.send-message')}</p>
      <div
        data-testid="send-confirmation"
        className="flex flex-col p-10 bg-secondary rounded-lg mb-6"
      >
        <div className="flex flex-row items-center pb-3 ">
          <div data-testid="send-confirmation-info" className="pr-2 text-info">
            {Intl.t('send-coins.send-confirmation', {
              amount: formattedAmount,
              fee,
            })}
          </div>
          <div
            className="flex flex-row relative items-center gap-1"
            onMouseEnter={() => setShowTooltip(true)}
            onMouseLeave={() => setShowTooltip(false)}
          >
            <ToolTip showTooltip={showTooltip} content={content} />
          </div>
        </div>
        <Balance
          customClass="p-0 bg-transparent mb-3"
          amount={formattedTotal}
        />
        <div className="text-info flex items-center gap-2">
          {Intl.t('send-coins.recipient')}
          <p data-testid="send-confirmation-recipient">
            {formattedRecipientAddress}
          </p>
        </div>
      </div>
      <Button disabled={isLoading} onClick={() => handleConfirm(true)}>
        {Intl.t('send-coins.confirm-sign')}
      </Button>
    </>
  );
}
