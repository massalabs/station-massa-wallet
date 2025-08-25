import {
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
  Selector,
  Identicon,
  MassaLogo,
  formatAmount,
} from '@massalabs/react-ui-kit';
import { maskNickname } from '@massalabs/react-ui-kit/src/lib/massa-react/utils';

import { AccountObject } from '@/models/AccountModel';

interface MyAccountsListProps {
  setRecipient: React.Dispatch<string>;
  accounts: AccountObject[];
  onClose: () => void;
}

function MyAccountsList(props: MyAccountsListProps) {
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
              amount={formatAmount(filteredAccount.candidateBalance).full}
              onClick={() => handleSetRecipient(filteredAccount)}
            />
          ))}
        </div>
      </PopupModalContent>
    </PopupModal>
  );
}

export default MyAccountsList;
