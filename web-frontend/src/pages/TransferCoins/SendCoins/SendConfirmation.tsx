import { useState } from 'react';
import Intl from '../../../i18n/i18n';
import { Balance, Button } from '@massalabs/react-ui-kit';
import { FiChevronLeft } from 'react-icons/fi';

import {
  maskAddress,
  formatStandard,
  toNanoMASS,
  toMASS,
} from '../../../utils/massaFormat';
import ToolTip from './ToolTip';

export function SendConfirmation({ ...props }) {
  const { data, handleConfirm, isLoading } = props;
  const { amount, fees, recipient } = data;

  const formattedRecipientAddress = maskAddress(recipient);
  const total = toNanoMASS(amount) + parseInt(fees);

  const formattedTotal = formatStandard(toMASS(total));
  const [showTooltip, setShowTooltip] = useState(false);

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
      <div className="flex flex-col p-10 bg-secondary rounded-lg mb-6">
        <div className="flex flex-row items-center pb-3 ">
          <div className="pr-2 text-s-info">
            {Intl.t('send-coins.send-confirmation', { amount, fees })}
          </div>
          <div
            className="flex flex-row relative items-center gap-1 "
            onMouseEnter={() => setShowTooltip(true)}
            onMouseLeave={() => setShowTooltip(false)}
          >
            <ToolTip
              showTooltip={showTooltip}
              content={Intl.t('send-coins.gas-info', { default: '1000' })}
            />
          </div>
        </div>
        <Balance
          customClass="p-0 bg-transparent mb-3"
          amount={formattedTotal}
        />
        <p className="text-s-info">
          {Intl.t('send-coins.recipient')}
          <p>{formattedRecipientAddress}</p>
        </p>
      </div>
      <Button disabled={isLoading} onClick={() => handleConfirm(true)}>
        {Intl.t('send-coins.confirm-sign')}
      </Button>
    </>
  );
}
