import {
  maskAddress,
  maskNickname,
} from '@massalabs/react-ui-kit/src/lib/massa-react/utils';

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
      <div className="flex flex-col">
        <p className="mas-menu-default truncate">{maskNickname(nickname)}</p>
        <p className="mas-caption">{maskAddress(walletAddress)}</p>
      </div>
    </div>
  );
}
