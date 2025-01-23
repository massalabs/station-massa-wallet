import {
  maskAddress,
  maskNickname,
} from '@massalabs/react-ui-kit/src/lib/massa-react/utils';

import { CopyClip } from './clipBoardCopy';
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
          <p className="mas-caption">{maskAddress(walletAddress)}</p>
          <div className="ml-2">
            <CopyClip data={walletAddress} />
          </div>
        </div>
      </div>
    </div>
  );
}
