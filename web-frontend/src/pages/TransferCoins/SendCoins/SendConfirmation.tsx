import { Balance } from '@massalabs/react-ui-kit';
import { FiChevronLeft, FiHelpCircle } from 'react-icons/fi';

export function SendConfirmation() {
  return (
    <div>
      <div className="flex flex-row just items-center hover:cursor-pointer mb-3.5">
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
      <div className="flex flex-col p-10 bg-secondary rounded-lg">
        <div className="flex flex-row items-center">
          <p>Amount (10 MAS) + gas fee (0,2 MAS)</p>
          <FiHelpCircle />
        </div>
        <Balance customClass="pl-0" amount={1000} />
      </div>
    </div>
  );
}
