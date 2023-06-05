import {
  MassaToken,
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
  Selector,
  Identicon,
} from '@massalabs/react-ui-kit';
import { formatStandard } from '../../../utils/MassaFormating';
import { useResource } from '../../../custom/api';
import { AccountObject } from '../../../models/AccountModel';

function AccountSelect({ ...props }) {
  const { modalAccounts, setModalAccounts, setRecipient, account } = props;
  const selectedAccount = account;
  const { data: accounts = [] } = useResource<AccountObject[]>('accounts');
  const filteredAccounts = accounts.filter(
    (account: AccountObject) => account.nickname !== selectedAccount.nickname,
  );

  function setRecipientAndClose(account: AccountObject) {
    setRecipient(account.address);
    setModalAccounts(false);
  }

  return (
    <PopupModal
      fullMode={true}
      onClose={() => setModalAccounts(!modalAccounts)}
    >
      <PopupModalHeader>
        <label className="mas-title">My accounts</label>
      </PopupModalHeader>
      <PopupModalContent>
        {filteredAccounts.map((account: AccountObject) => (
          <Selector
            customClass="pb-4"
            key={account.nickname}
            preIcon={<Identicon username={account.nickname} size={32} />}
            posIcon={<MassaToken size={24} />}
            content={account.nickname}
            variant="secondary"
            amount={formatStandard(+account.candidateBalance)}
            onClick={() => setRecipientAndClose(account)}
          />
        ))}
      </PopupModalContent>
    </PopupModal>
  );
}

export default AccountSelect;
