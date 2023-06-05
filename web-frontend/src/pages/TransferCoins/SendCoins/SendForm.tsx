import { Balance, Button, Input } from '@massalabs/react-ui-kit';
import { FiArrowUpRight, FiPlus } from 'react-icons/fi';
import Modal from './Modal';
import ContactList from './ContactList';

export function SendForm({ ...props }) {
  const {
    amount,
    account,
    formattedBalance,
    recipient,
    error,
    fees,
    modal,
    modalAccounts,
    setModal,
    setModalAccounts,
    handleModal,
    handleModalAccounts,
    handleFees,
    handleFeesConfirm,
    setRecipient,
    handleSubmit,
    handleChange,
    SendPercentage,
  } = props;

  const modalArgs = {
    modal,
    setModal,
    handleModal,
    fees,
    handleFees,
    handleFeesConfirm,
  };
  const modalArgsAccounts = {
    modalAccounts,
    setModalAccounts,
    handleModalAccounts,
    setRecipient,
    account,
  };

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
        <div className="pb-3.5">
          <div className="flex flex-row justify-between w-full pb-3.5 ">
            <p>Send Token</p>
            <p>
              Available : <u>{formattedBalance}</u>
            </p>
          </div>
          <div className="pb-3.5">
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
          <div className="pb-3.5">
            <p>To :</p>
          </div>
          <div className="pb-3.5">
            <Input
              placeholder={'Recipient'}
              value={recipient}
              onChange={(e) => setRecipient(e.target.value)}
              name="recipient"
              error={error?.recipient}
            />
          </div>
          <div className="flex flex-row-reverse pb-3.5">
            <p
              className="hover:cursor-pointer"
              onClick={() => setModalAccounts(!modalAccounts)}
            >
              <u>Transfer between my accounts</u>
            </p>
          </div>
        </div>
        {/* Button Section */}
        <div className="flex flex-col w-full">
          <div className="pb-3.5">
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
      <div>
        {modal ? (
          <Modal {...modalArgs} />
        ) : modalAccounts ? (
          <ContactList {...modalArgsAccounts} />
        ) : null}
      </div>
    </div>
  );
}
