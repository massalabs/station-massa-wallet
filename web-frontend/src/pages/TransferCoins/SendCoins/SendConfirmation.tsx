import { useState } from 'react';

import { Balance, Button } from '@massalabs/react-ui-kit';
import { FiChevronLeft } from 'react-icons/fi';

import ToolTip from './ToolTip';
import Intl from '@/i18n/i18n';
import { maskAddress, formatStandard, toNanoMASS, toMASS } from '@/utils';

export interface SendConfirmationData {
  amount: string;
  fees: string;
  recipient: string;
}

interface SendConfirmationProps {
  data: SendConfirmationData;
  handleConfirm: (confirm: boolean) => void;
  isLoading: boolean;
}

export function SendConfirmation(props: SendConfirmationProps) {
  const { data, handleConfirm, isLoading } = props;
  const { amount, fees, recipient } = data;

  const GAS_STANDARD = `${Intl.t('send-coins.fee-standard')}`;
  const GAS_LOW = `${Intl.t('send-coins.fee-low')}`;
  const GAS_HIGH = `${Intl.t('send-coins.fee-high')}`;
  const GAS_CUSTOM = `${Intl.t('send-coins.fee-custom')}`;

  const formattedRecipientAddress = maskAddress(recipient);
  const total = toNanoMASS(amount) + BigInt(fees);

  const formattedTotal = formatStandard(toMASS(total));
  const [showTooltip, setShowTooltip] = useState(false);
  let selectedFees;

  switch (fees) {
    case '1000':
      selectedFees = GAS_STANDARD;
      break;
    case '1':
      selectedFees = GAS_LOW;
      break;
    case '5000':
      selectedFees = GAS_HIGH;
      break;
    default:
      selectedFees = GAS_CUSTOM;
      break;
  }

  const feeInfo = `${Intl.t('send-coins.fee-info', {
    feeType: selectedFees,
    fees: fees,
  })}`;
  const gasAlert = `  \u26A0  ${Intl.t('send-coins.fee-alert')}`;
  let content = selectedFees == GAS_STANDARD ? feeInfo + gasAlert : gasAlert;

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
            {Intl.t('send-coins.send-confirmation', { amount, fees })}
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
