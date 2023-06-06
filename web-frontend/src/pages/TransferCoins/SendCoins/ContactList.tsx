import {
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
  Selector,
  Identicon,
  MassaLogo,
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
      customClass="!w-1/2 h-1/2 "
    >
      <PopupModalHeader>
        <label className="mas-title mb-6">My accounts</label>
      </PopupModalHeader>
      <PopupModalContent>
        <div className="overflow-scroll h-80">
          {filteredAccounts.map((account: AccountObject) => (
            <Selector
              customClass="pb-4"
              key={account.nickname}
              preIcon={<Identicon username={account.nickname} size={32} />}
              posIcon={<MassaLogo size={24} />}
              content={account.nickname}
              variant="secondary"
              amount={formatStandard(+account.candidateBalance / 10 ** 9)}
              onClick={() => setRecipientAndClose(account)}
            />
          ))}
        </div>
      </PopupModalContent>
    </PopupModal>
  );
}

export default AccountSelect;
