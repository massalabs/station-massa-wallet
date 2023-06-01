import { Balance, Button } from '@massalabs/react-ui-kit';
import { FiChevronLeft, FiHelpCircle } from 'react-icons/fi';
import {
  formatStandard,
  reverseFormatStandard,
} from '../../../utils/MassaFormating';

export function SendConfirmation({ ...props }) {
  const { amount, recipient, valid, setValid } = props;

  const reversedAmount = reverseFormatStandard(amount);

  console.log(formatStandard(reversedAmount));

  return (
    <div>
      <div
        onClick={() => setValid(!valid)}
        className="flex flex-row just items-center hover:cursor-pointer mb-3.5"
      >
        <div className="mr-2">
          <FiChevronLeft />
        </div>
        <p>
          <u>Back to Sending page</u>
        </p>
      </div>
      <div className="mb-3.5">
        <p>You're going to send: </p>
      </div>
      <div className="flex flex-col p-10 bg-secondary rounded-lg mb-3.5">
        <div className="flex flex-row items-center mb-3.5 ">
          <div className="mr-2">
            <p>Amount ({formatStandard(reversedAmount)}) + gas fee (0.2 MAS)</p>
          </div>
          <div>
            <FiHelpCircle />
          </div>
        </div>
        <div className="mb-3.5">
          <Balance
            customClass="pl-0 bg-transparent"
            amount={formatStandard(reversedAmount)}
          />
        </div>
        <div>
          <p>
            To : <u>{recipient}</u>
          </p>
        </div>
      </div>
      <div>
        <Button> Confirm and sign with my password </Button>
      </div>
    </div>
  );
}
