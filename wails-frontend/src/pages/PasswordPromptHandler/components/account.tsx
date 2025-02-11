import { Clipboard, maskAddress, maskNickname } from '@massalabs/react-ui-kit';

import Intl from '@/i18n/i18n';

interface FromProps {
  nickname: string;
  walletAddress: string;
}

export function Account(props: FromProps) {
  const { nickname, walletAddress } = props;

  return (
    <div className="flex w-full justify-between pt-4">
      <p>{Intl.t('password-prompt.account')}</p>
      <div className="flex flex-col">
        <p className="mas-menu-default truncate">{maskNickname(nickname)}</p>
        <div className="flex items-center">
          <Clipboard
            customClass="h-8 w-10"
            rawContent={walletAddress}
            displayedContent={maskAddress(walletAddress)}
          />
        </div>
      </div>
    </div>
  );
}
