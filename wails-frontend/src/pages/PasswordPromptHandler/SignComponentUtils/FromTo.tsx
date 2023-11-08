import { FiArrowRight } from 'react-icons/fi';

import Intl from '@/i18n/i18n';
import { maskAddress, maskNickname } from '@/utils';

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
      <div className="flex flex-col">
        <div className="flex gap-2">
          <p className="mas-menu-active">
            {Intl.t('password-prompt.sign.from')}
          </p>
          <p className="mas-menu-default">{maskNickname(nickname)}</p>
        </div>
        <p className="mas-caption">{maskAddress(walletAddress)}</p>
      </div>
      <div className="h-8 w-8 rounded-full flex items-center justify-center bg-neutral">
        <FiArrowRight size={24} className="text-primary" />
      </div>
      <div className="flex flex-col">
        <div className="flex gap-2">
          <p className="mas-menu-active">{Intl.t(label)}</p>
          {recipientNickname ? (
            <p className="mas-menu-default">
              {maskNickname(recipientNickname)}
            </p>
          ) : null}
        </div>
        <p className="mas-caption">{maskAddress(recipientAddress)}</p>
      </div>
    </div>
  );
}
