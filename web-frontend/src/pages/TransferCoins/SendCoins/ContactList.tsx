import {
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
  Selector,
  Identicon,
  MassaLogo,
} from '@massalabs/react-ui-kit';

import { AccountObject } from '@/models/AccountModel';
import { formatStandard, maskNickname } from '@/utils/massaFormat';

interface ContactListProps {
  setRecipient: React.Dispatch<string>;
  accounts: AccountObject[];
  onClose: () => void;
}

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
      customClass="!w-1/2 h-1/2"
    >
      <PopupModalHeader>
        <label className="mas-title mb-6">My accounts</label>
      </PopupModalHeader>
      <PopupModalContent>
        <div
          data-testid="selector-account-list"
          className="overflow-scroll h-80"
        >
          {accounts.map((filteredAccount: AccountObject, index: number) => (
            <Selector
              customClass="pb-4"
              key={index}
              data-testid={`selector-account-${index}`}
              preIcon={<Identicon username={filteredAccount.nickname} />}
              posIcon={<MassaLogo size={24} />}
              content={maskNickname(filteredAccount.nickname)}
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
