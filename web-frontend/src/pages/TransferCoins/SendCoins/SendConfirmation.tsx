import { Balance, Button } from '@massalabs/react-ui-kit';
import { FiChevronLeft, FiHelpCircle } from 'react-icons/fi';
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

export function SendConfirmation({ ...props }) {
  const navigate = useNavigate();
  const { amount, nickname, recipient, valid, fees, setValid } = props;
  const reversedAmount = reverseFormatStandard(amount);
  const amountInNanoMAS = fromMAS(reversedAmount).toString();
  const total = +amountInNanoMAS + +fees;
  const formattedTotal = formatStandard(total, Unit.NanoMAS);
  const [showTooltip, setShowTooltip] = useState(false); // Add state for tooltip visibility

  const { mutate, isSuccess } = usePost<SendTransactionObject>(
    `accounts/${nickname}/transfer`,
  );
  if (isSuccess) {
    navigate(routeFor('index'));
  }
  const handleTransfer = () => {
    const transferData: SendTransactionObject = {
      amount: amount,
      recipientAddress: recipient,
      fee: fees,
    };

    mutate(transferData as SendTransactionObject);
  };

  return (
    <div>
      <div
        onClick={() => setValid(!valid)}
        className="flex flex-row just items-center hover:cursor-pointer pb-3.5"
      >
        <div className="mr-2">
          <FiChevronLeft />
        </div>
        <p>
          <u>{Intl.t('sendcoins.back-to-sending')}</u>
        </p>
      </div>
      <div className="pb-3.5">
        <p>{Intl.t('sendcoins.send-message')}</p>
      </div>
      <div className="flex flex-col p-10 bg-secondary rounded-lg mb-3.5">
        <div className="flex flex-row items-center pb-3.5 ">
          <div className="pr-2">
            <p>
              {Intl.t('sendcoins.amount')} ({formatStandard(reversedAmount)})
              {Intl.t('sendcoins.mas-gas')} {fees}
              {Intl.t('sendcoins.nano-mas')}
            </p>
          </div>
          <div
            className="flex flex-row relative items-center gap-1 "
            onMouseEnter={() => setShowTooltip(true)}
            onMouseLeave={() => setShowTooltip(false)}
          >
            <FiHelpCircle />
            {showTooltip && (
              <div className="flex flex-col w-[368px] absolute z-10 t-4 bg-tertiary p-3 rounded ">
                {Intl.t('sendcoins.gas-info').replace('XX', fees)}
              </div>
            )}
          </div>
        </div>
        <div className="pb-3.5">
          <Balance customClass="pl-0 bg-transparent" amount={formattedTotal} />
        </div>
        <div>
          <p>
            {Intl.t('sendcoins.recipient')}
            <u>{recipient.slice(0, 4) + '...' + recipient.slice(-4)}</u>
          </p>
        </div>
      </div>
      <div className="mt-3.5">
        <Button
          onClick={() => {
            handleTransfer();
          }}
        >
          {Intl.t('sendcoins.confirm-sign')}
        </Button>
      </div>
    </div>
  );
}
