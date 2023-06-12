import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

import { routeFor } from '../../../utils';
import Intl from '../../../i18n/i18n';
import { fromMAS } from '@massalabs/massa-web3';

import { Balance, Button } from '@massalabs/react-ui-kit';
import { FiChevronLeft } from 'react-icons/fi';
import {
  Unit,
  formatStandard,
  maskAddress,
  reverseFormatStandard,
} from '../../../utils/MassaFormating';
import { usePost } from '../../../custom/api';
import { SendTransactionObject } from '../../../models/AccountModel';

import ToolTip from './ToolTip';

export function SendConfirmation({ ...props }) {
  const { amount, nickname, recipient, valid, fees, setValid } = props;

  const navigate = useNavigate();
  const formattedRecipientAddress = maskAddress(recipient);
  const reversedAmount = reverseFormatStandard(amount);
  const amountInNanoMAS = fromMAS(reversedAmount).toString();
  const total = +amountInNanoMAS + +fees;
  const formattedTotal = formatStandard(total, Unit.NanoMAS);
  const [showTooltip, setShowTooltip] = useState(false);

  const { mutate, isSuccess, error } = usePost<SendTransactionObject>(
    `accounts/${nickname}/transfer`,
  );

  useEffect(() => {
    if (error) {
      navigate(routeFor('error'));
    } else if (isSuccess) {
      navigate(routeFor(`${nickname}/home`));
    }
  }, [isSuccess]);

  const handleTransfer = () => {
    const transferData: SendTransactionObject = {
      amount: amountInNanoMAS,
      recipientAddress: recipient,
      fee: fees.toString(),
    };

    mutate(transferData as SendTransactionObject);
  };

  const [customFees, setCustomFees] = useState<string>('Standard');

  useEffect(() => {
    switch (fees) {
      case '1000':
        setCustomFees('Standard');
        break;
      case '1':
        setCustomFees('Low');
        break;
      case '5000':
        setCustomFees('High');
        break;
      default:
        setCustomFees('Custom');
    }
  }, [fees]);

  const content = `
  ${Intl.t('sendcoins.gas-info', {
    default: customFees,
    gasFees: fees,
  })}  \u26A0  ${Intl.t('sendcoins.gas-alert')}`;

  const ToolTipArgs = {
    showTooltip,
    content,
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
          <div className="pr-2 text-s-info">
            {Intl.t('sendcoins.send-confirmation', {
              amount: formatStandard(reversedAmount),
              fees: fees,
            })}
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
        <p className="text-s-info">
          {Intl.t('sendcoins.recipient')}
          <u>{formattedRecipientAddress}</u>
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
