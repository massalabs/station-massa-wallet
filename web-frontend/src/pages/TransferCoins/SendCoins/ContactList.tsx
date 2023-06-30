import {
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
  Selector,
  Identicon,
  MassaLogo,
} from '@massalabs/react-ui-kit';
import { formatStandard } from '../../../utils/massaFormat';
import { AccountObject } from '../../../models/AccountModel';
import { ContactListProps } from './SendForm';

function AccountSelect(props: ContactListProps) {
  const { onClose, setRecipient, accounts } = props;

  function handleSetRecipient(filteredAccount: AccountObject) {
    let address = filteredAccount?.address || '';

    setRecipient?.(address);
    onClose?.();
  }

  return (
    <PopupModal
      fullMode={true}
      onClose={() => onClose?.()}
      customClass="!w-1/2 h-1/2 "
    >
      <PopupModalHeader>
        <label className="mas-title mb-6">My accounts</label>
      </PopupModalHeader>
      <PopupModalContent>
        <div className="overflow-scroll h-80">
          {accounts.map((filteredAccount: AccountObject) => (
            <Selector
              customClass="pb-4"
              key={filteredAccount.nickname}
              preIcon={<Identicon username={filteredAccount.nickname} />}
              posIcon={<MassaLogo size={24} />}
              content={filteredAccount.nickname}
              variant="secondary"
              amount={formatStandard(
                +filteredAccount.candidateBalance / 10 ** 9,
              )}
              onClick={() => handleSetRecipient(filteredAccount)}
            />
          ))}
        </div>
      </PopupModalContent>
    </PopupModal>
  );
}

export default AccountSelect;
