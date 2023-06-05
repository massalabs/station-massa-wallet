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

import { toMAS } from '@massalabs/massa-web3';

function AccountSelect({ ...props }) {
  const { modalAccounts, setModalAccounts, setRecipient, account } = props;
  const { data: accounts = [] } = useResource<AccountObject[]>('accounts');
  const otherAccounts = accounts.filter(
    (acnt: AccountObject) => acnt.nickname !== account.nickname,
  );

  function getFormattedBalance(account: AccountObject): string {
    return formatStandard(toMAS(account.candidateBalance).toNumber());
  }
  function setRecipientAndClose(account: AccountObject) {
    console.log(
      'setting recipient as ',
      account.nickname,
      ' address ',
      account.address,
    );
    setRecipient(account.address);
    setModalAccounts(false);
  }

  return (
    <PopupModal
      fullMode={true}
      onClose={() => setModalAccounts(!modalAccounts)}
    >
      <PopupModalHeader>
        <div>
          <label className="mas-title">My accounts</label>
        </div>
      </PopupModalHeader>
      <PopupModalContent>
        <div>
          {otherAccounts.map((account: AccountObject) => (
            <div className="pb-4" key={account.nickname}>
              <Selector
                preIcon={<Identicon username={account.nickname} size={32} />}
                posIcon={<MassaToken size={24} />}
                content={account.nickname}
                variant="secondary"
                amount={getFormattedBalance(account)}
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
