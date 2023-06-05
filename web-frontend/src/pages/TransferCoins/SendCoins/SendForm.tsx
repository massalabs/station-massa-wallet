import { Balance, Button, Input } from '@massalabs/react-ui-kit';
import { FiArrowUpRight, FiPlus } from 'react-icons/fi';
import Intl from '../../../i18n/i18n';
import Modal from './Advanced';
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
    handleModal,
    setModalAccounts,
    handleModalAccounts,
    handleFees,
    setRecipient,
    handleSubmit,
    handleChange,
    SendPercentage,
  } = props;

  const modalArgs = {
    fees,
    handleModal,
    handleFees,
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
          <p> {Intl.t('sendcoins.account-balance')}</p>
          <Balance
            customClass="pl-0 bg-transparent"
            amount={formattedBalance}
          />
        </div>
        <div className="pb-3.5">
          <div className="flex flex-row justify-between w-full pb-3.5 ">
            <p> {Intl.t('sendcoins.send-action')} </p>
            <p>
              {Intl.t('sendcoins.available-balance')} <u>{formattedBalance}</u>
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
            <p>{Intl.t('sendcoins.recipient')}</p>
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
              <u>{Intl.t('sendcoins.transfer-between-acc')}</u>
            </p>
          </div>
        </div>
        {/* Button Section */}
        <div className="flex flex-col w-full">
          <div className="pb-3.5">
            <Button
              onClick={() => handleModal()}
              variant={'secondary'}
              posIcon={<FiPlus />}
            >
              {Intl.t('sendcoins.advanced')}
            </Button>
          </div>
          <div>
            <Button type="submit" posIcon={<FiArrowUpRight />}>
              {Intl.t('sendcoins.send')}
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
