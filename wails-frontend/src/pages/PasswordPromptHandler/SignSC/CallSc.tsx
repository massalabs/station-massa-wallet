import { FiArrowRight } from 'react-icons/fi';

import { Description } from '../Description';
import Intl from '@/i18n/i18n';
import { SignBodyProps } from '@/pages/PasswordPromptHandler/Sign';
import { maskAddress, maskNickname } from '@/utils';

export function CallSc(props: SignBodyProps) {
  const {
    Address,
    WalletAddress,
    Function: CalledFunction,
    OperationType,
    Description: description,
    Nickname,
    children,
  } = props;

  return (
    <div className="flex flex-col items-center gap-4 mas-menu-default w-[326px]">
      <div className="flex w-full items-center justify-between">
        <div className="flex flex-col">
          <div className="flex gap-2">
            <p className="mas-menu-active">
              {Intl.t('password-prompt.sign.from')}
            </p>
            <p className="mas-menu-default">{maskNickname(Nickname)}</p>
          </div>
          <p className="mas-caption">{maskAddress(WalletAddress)}</p>
        </div>
        <div className="h-8 w-8 rounded-full flex items-center justify-center bg-neutral">
          <FiArrowRight size={24} className="text-primary" />
        </div>
        <div className="flex flex-col">
          <p className="mas-menu-active">
            {Intl.t('password-prompt.sign.contract')}
          </p>
          <p className="mas-caption">{maskAddress(Address)}</p>
        </div>
      </div>

      <hr className="h-0.25 bg-neutral opacity-40 w-full" />

      <div className="flex justify-between w-full">
        <div className="flex flex-col h-fit">
          <p>{Intl.t('password-prompt.sign.operation-type')}</p>
          <p className="mas-caption">
            {Intl.t('password-prompt.sign.function-called')}
          </p>
        </div>
        <div className="flex flex-col items-end h-fit">
          <p>
            {Intl.t(`password-prompt.sign.operation-types.${OperationType}`)}
          </p>
          <p className="mas-caption">{CalledFunction}</p>
        </div>
      </div>

      <hr className="h-0.25 bg-neutral opacity-40 w-full" />

      <Description description={description} />

      {children}
    </div>
  );
}
