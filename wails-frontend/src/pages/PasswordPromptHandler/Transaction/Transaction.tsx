import { FiArrowRight } from 'react-icons/fi';

import { SignBodyProps } from '../Sign';
import { Description } from '../SignComponentUtils/Description';
import Intl from '@/i18n/i18n';
import {
  formatStandard,
  masToken,
  maskAddress,
  Unit,
  maskNickname,
} from '@/utils';

export function Transaction(props: SignBodyProps) {
  const {
    WalletAddress,
    RecipientAddress,
    RecipientNickname,
    OperationType,
    Amount,
    Description: description,
    Nickname,
    children,
  } = props;

  return (
    <div className="flex flex-col items-center gap-4 mas-menu-default">
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
          <div className="flex gap-2">
            <p className="mas-menu-active">
              {Intl.t('password-prompt.sign.to')}
            </p>
            {RecipientAddress ? (
              <p className="mas-menu-default">
                {maskNickname(RecipientNickname)}
              </p>
            ) : null}
          </div>
          <p className="mas-caption">{maskAddress(RecipientAddress)}</p>
        </div>
      </div>

      <hr className="h-0.25 bg-neutral opacity-40 w-full" />

      <div className="flex justify-between w-full">
        <div className="flex flex-col h-fit">
          <p>{Intl.t('password-prompt.sign.operation-type')}</p>
        </div>
        <div className="flex flex-col items-end h-fit">
          <p>
            {Intl.t(`password-prompt.sign.operation-types.${OperationType}`)}
          </p>
        </div>
      </div>

      <div className="flex flex-col gap-2 w-full">
        <div className="flex w-full items-center justify-between">
          <p>{Intl.t('password-prompt.sign.sending-amount')}</p>
          <p>
            {formatStandard(Amount, Unit.NanoMAS)} {masToken}
          </p>
        </div>
      </div>

      <hr className="h-0.25 bg-neutral opacity-40 w-full" />

      <Description description={description} />

      {children}
    </div>
  );
}
