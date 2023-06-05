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
  const { data: accounts = [] } = useResource<AccountObject[]>('accounts');
  const otherAccounts = accounts.filter(
    (acnt: AccountObject) => acnt.nickname !== account.nickname,
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
        <div>
          {otherAccounts.map((account: AccountObject) => (
            <div className="pb-4">
              <Selector
                key={account.nickname}
                preIcon={<Identicon username={account.nickname} size={32} />}
                posIcon={<MassaToken size={24} />}
                content={account.nickname}
                variant="secondary"
                amount={formatStandard(+account.candidateBalance)}
                onClick={() => setRecipientAndClose(account)}
              />
            </div>
          ))}
        </div>
      </PopupModalContent>
    </PopupModal>
  );
}

export default AccountSelect;
