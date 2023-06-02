import {
  Balance,
  Button,
  Input,
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
  RadioButton,
} from '@massalabs/react-ui-kit';
import { useState } from 'react';
import { FiArrowUpRight, FiPlus } from 'react-icons/fi';

export function SendForm({ ...props }) {
  const {
    amount,
    formattedBalance,
    recipient,
    fees,
    error,
    handleFees,
    setRecipient,
    handleSubmit,
    handleChange,
    handleFeesConfirm,
    SendPercentage,
  } = props;

  const radioArgs = {
    name: 'radio',
  };

  const [modal, setModal] = useState<boolean>(false);

  function Modal() {
    return (
      <PopupModal fullMode={true} onClose={() => setModal(!modal)}>
        <PopupModalHeader>
          <div>
            <label className="mas-title">Advanced</label>
            <p className="mas-body">
              You pay gas fees to reward block validators and maximize your
              chances to see your transaction validated. It’s a tip for people
              that support the blockchain network.
            </p>
          </div>
        </PopupModalHeader>
        <PopupModalContent>
          <div>
            <div className="flex flex-row items-center mas-buttons pb-3.5">
              <RadioButton defaultChecked={true} {...radioArgs} />
              <p>Preset : </p>
            </div>
            <div className="flex flex-row items-center w-full gap-4  pb-3.5">
              <div className="w-full">
                <Button
                  variant={fees == 1000 ? 'primary' : 'secondary'}
                  onClick={() => handleFees(1000)}
                >
                  Standard
                  <label className="text-info text-xs flex ml-1 items-center">
                    (1000 nMAS)
                  </label>
                </Button>
              </div>
              <div className="w-full">
                <Button
                  variant={fees == 5000 ? 'primary' : 'secondary'}
                  onClick={() => handleFees(5000)}
                >
                  High
                  <label className="text-info text-xs flex pl-1 items-center">
                    (5000 nMAS)
                  </label>
                </Button>
              </div>
              <div className="w-full">
                <Button
                  variant={fees === 1 ? 'primary' : 'secondary'}
                  onClick={() => handleFees(1)}
                >
                  Low
                  <label className="text-info text-xs flex pl-1 items-center">
                    (1 uMAS)
                  </label>
                </Button>
              </div>
            </div>
            <div>
              <div className="flex flex-row items-center mas-buttons pb-3.5">
                <RadioButton defaultChecked={false} {...radioArgs} />
                <p> Custom fees : </p>
              </div>
              <form className="pb-3.5" onSubmit={handleFeesConfirm}>
                <div className="pb-3.5">
                  <Input placeholder={'Gas fees amount'} name="fees" />
                </div>
                <div>
                  <Button type="submit"> Confirm Fees </Button>
                </div>
              </form>
            </div>
          </div>
        </PopupModalContent>
      </PopupModal>
    );
  }

  function handleModal() {
    setModal(!modal);
  }

  return (
    <div>
      <form onSubmit={handleSubmit}>
        {/* Balance Section */}
        <div>
          <p> Account Balance</p>
          <Balance
            customClass="pl-0 bg-transparent"
            amount={formattedBalance}
          />
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
      <div>{modal ? <Modal /> : null}</div>
    </div>
  );
}
