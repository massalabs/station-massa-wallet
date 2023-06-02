import {
  Balance,
  Button,
  Input,
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
} from '@massalabs/react-ui-kit';
import { useState } from 'react';
import { FiArrowUpRight, FiPlus } from 'react-icons/fi';

export function SendForm({ ...props }) {
  const {
    amount,
    formattedBalance,
    recipient,
    error,
    setRecipient,
    handleSubmit,
    handleChange,
    SendPercentage,
  } = props;

  const [modal, setModal] = useState<boolean>(false);

  function Modal() {
    return (
      <PopupModal fullMode={true} onClose={() => setModal(!modal)}>
        <PopupModalHeader>
          <label className="text-f-primary">The title</label>
        </PopupModalHeader>
        <PopupModalContent>
          <label className="text-f-primary">any content</label>
        </PopupModalContent>
      </PopupModal>
    );
  }

  function handleModal() {
    setModal(!modal);
  }

  return (
    <form onSubmit={handleSubmit}>
      {/* Balance Section */}
      <div>
        <p> Account Balance</p>
        <Balance customClass="pl-0 bg-transparent" amount={formattedBalance} />
      </div>
      <div className="mb-3.5">
        <div className="flex flex-row justify-between w-full mb-3.5 ">
          <p>Send Token</p>
          <p>
            Available : <u>{formattedBalance}</u>
          </p>
        </div>
        <div className="mb-3.5">
          <Input
            placeholder={'Amount to send'}
            value={amount}
            name="amount"
            onChange={(e) => handleChange(e)}
            error={error?.amount}
          />
        </div>
        <div className="flex flex-row-reverse">
          <ul className="flex flex-row">
            <SendPercentage percentage={25} />
            <SendPercentage percentage={50} />
            <SendPercentage percentage={75} />
            <SendPercentage percentage={100} />
          </ul>
        </div>
      </div>
      <div>
        <div className="mb-3.5">
          <p>To :</p>
        </div>
        <div className="mb-3.5">
          <Input
            placeholder={'Recipient'}
            value={recipient}
            onChange={(e) => setRecipient(e.target.value)}
            name="recipient"
            error={error?.recipient}
          />
        </div>
        <div className="flex flex-row-reverse mb-3.5">
          <p
            className="hover:cursor-pointer"
            onClick={() => console.log('transfer between accounts')}
          >
            <u>Transfer between my account</u>
          </p>
        </div>
      </div>
      {/* Button Section */}

      {modal ? <Modal /> : null}

      <div className="flex flex-col w-full">
        <div className="mb-3.5">
          <Button
            onClick={handleModal}
            variant={'secondary'}
            posIcon={<FiPlus />}
          >
            Advanced
          </Button>
        </div>
        <div>
          <Button type="submit" posIcon={<FiArrowUpRight />}>
            Send
          </Button>
        </div>
      </div>
    </form>
  );
}
