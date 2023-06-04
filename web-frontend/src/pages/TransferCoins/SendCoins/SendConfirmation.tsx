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
import { useNavigate } from 'react-router-dom';
import { fromMAS } from '@massalabs/massa-web3';

export function SendConfirmation({ ...props }) {
  const navigate = useNavigate();
  const { amount, nickname, recipient, valid, fees, setValid } = props;
  const reversedAmount = reverseFormatStandard(amount);
  const amountInNanoMAS = fromMAS(reversedAmount).toString();
  const total = +amountInNanoMAS + +fees;
  const formattedTotal = formatStandard(total, Unit.NanoMAS);

  const { mutate, isSuccess } = usePost<SendTransactionObject>(
    `accounts/${nickname}/transfer`,
  );
  if (isSuccess) {
    navigate(routeFor('account-create-step-three'), { state: { nickname } });
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
          <u>Back to Sending page</u>
        </p>
      </div>
      <div className="pb-3.5">
        <p>You're going to send: </p>
      </div>
      <div className="flex flex-col p-10 bg-secondary rounded-lg mb-3.5">
        <div className="flex flex-row items-center pb-3.5 ">
          <div className="pr-2">
            <p>
              Amount ({formatStandard(reversedAmount)}) MAS + gas fee {fees}{' '}
              nMAS
            </p>
          </div>
          <div>
            <FiHelpCircle />
          </div>
        </div>
        <div className="pb-3.5">
          <Balance customClass="pl-0 bg-transparent" amount={formattedTotal} />
        </div>
        <div>
          <p>
            To : <u>{recipient.slice(0, 4) + '...' + recipient.slice(-4)}</u>
          </p>
        </div>
      </div>
      <div className="mt-3.5">
        <Button
          onClick={() => {
            handleTransfer();
          }}
        >
          Confirm and sign with my password
        </Button>
      </div>
    </div>
  );
}
