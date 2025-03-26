import { maskAddress, maskNickname, Clipboard } from '@massalabs/react-ui-kit';
import { FiArrowRight } from 'react-icons/fi';

import Intl from '@/i18n/i18n';

interface FromToProps {
  fromNickname: string;
  fromAddress: string;
  toNickname?: string;
  toAddress: string;
  label: string;
}

export function FromTo(props: FromToProps) {
  const {
    fromNickname: nickname,
    fromAddress: walletAddress,
    toNickname: recipientNickname,
    toAddress: recipientAddress,
    label,
  } = props;

  return (
    <div className="flex w-full items-center justify-between">
      <div className="flex flex-col gap-y-2">
        <div className="flex gap-2">
          <p className="mas-menu-active">
            {Intl.t('password-prompt.sign.from')}
          </p>
          <p className="mas-menu-default">{maskNickname(nickname)}</p>
        </div>
        <Clipboard
          customClass="h-8 w-10"
          rawContent={walletAddress}
          displayedContent={maskAddress(walletAddress)}
        />
      </div>
      <div className="h-8 w-8 rounded-full flex items-center justify-center bg-neutral">
        <FiArrowRight size={24} className="text-primary" />
      </div>
      <div className="flex flex-col gap-y-2">
        <div className="flex gap-2">
          <p className="mas-menu-active">{Intl.t(label)}</p>
          {recipientNickname ? (
            <p className="mas-menu-default">
              {maskNickname(recipientNickname)}
            </p>
          ) : null}
        </div>
        <Clipboard
          customClass="h-8 w-10"
          rawContent={recipientAddress}
          displayedContent={maskAddress(recipientAddress)}
        />
      </div>
    </div>
  );
}
