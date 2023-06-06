import { Balance, Button } from '@massalabs/react-ui-kit';
import { FiChevronLeft } from 'react-icons/fi';
import {
  Unit,
  formatStandard,
  reverseFormatStandard,
} from '../../../utils/MassaFormating';
import { usePost } from '../../../custom/api';
import { SendTransactionObject } from '../../../models/AccountModel';
import { routeFor } from '../../../utils';
import Intl from '../../../i18n/i18n';
import { useNavigate } from 'react-router-dom';
import { fromMAS } from '@massalabs/massa-web3';
import { useState } from 'react';
import ToolTip from './ToolTip';

export function SendConfirmation({ ...props }) {
  const navigate = useNavigate();
  const { amount, nickname, recipient, valid, fees, setValid } = props;
  const reversedAmount = reverseFormatStandard(amount);
  const amountInNanoMAS = fromMAS(reversedAmount).toString();
  const total = +amountInNanoMAS + +fees;
  const formattedTotal = formatStandard(total, Unit.NanoMAS);
  const [showTooltip, setShowTooltip] = useState(false);

  const { mutate, isSuccess } = usePost<SendTransactionObject>(
    `accounts/${nickname}/transfer`,
  );
  if (isSuccess) {
    navigate(routeFor(`${nickname}/home`));
  }
  const handleTransfer = () => {
    const transferData: SendTransactionObject = {
      amount: amountInNanoMAS,
      recipientAddress: recipient,
      fee: fees.toString(),
    };

    mutate(transferData as SendTransactionObject);
  };
  const ToolTipArgs = {
    showTooltip,
    content: Intl.t('sendcoins.gas-info').replace('XX', fees),
  };

  return (
    <>
      <div
        onClick={() => setValid(!valid)}
        className="flex flex-row just items-center hover:cursor-pointer mb-5 gap-2"
      >
        <FiChevronLeft />

        <u>{Intl.t('sendcoins.back-to-sending')}</u>
      </div>
      <p className="mb-6">{Intl.t('sendcoins.send-message')}</p>
      <div className="flex flex-col p-10 bg-secondary rounded-lg mb-6">
        <div className="flex flex-row items-center pb-3 ">
          <div className="pr-2">
            {Intl.t('sendcoins.amount')} ({formatStandard(reversedAmount)})
            {Intl.t('sendcoins.mas-gas')} ({fees}){Intl.t('sendcoins.nano-mas')}
          </div>
          <div
            className="flex flex-row relative items-center gap-1 "
            onMouseEnter={() => setShowTooltip(true)}
            onMouseLeave={() => setShowTooltip(false)}
          >
            <ToolTip {...ToolTipArgs} />
          </div>
        </div>
        <Balance
          customClass="p-0 bg-transparent mb-3"
          amount={formattedTotal}
        />
        <p>
          {Intl.t('sendcoins.recipient')}
          <u>{recipient.slice(0, 4) + '...' + recipient.slice(-4)}</u>
        </p>
      </div>
      <Button
        onClick={() => {
          handleTransfer();
        }}
      >
        {Intl.t('sendcoins.confirm-sign')}
      </Button>
    </>
  );
}
