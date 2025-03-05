import { maskAddress, maskNickname, Clipboard } from '@massalabs/react-ui-kit';

import Intl from '@/i18n/i18n';

interface FromProps {
  nickname: string;
  walletAddress: string;
}

export function From(props: FromProps) {
  const { nickname, walletAddress } = props;

  return (
    <div className="flex w-full justify-between">
      <p>{Intl.t('password-prompt.sign.from')}</p>
      <div className="flex flex-col gap-y-2">
        <p className="mas-menu-default truncate">{maskNickname(nickname)}</p>
        <Clipboard
          customClass="h-8 w-10"
          rawContent={walletAddress}
          displayedContent={maskAddress(walletAddress)}
        />
      </div>
    </div>
  );
}
