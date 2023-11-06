import { Description } from '../Description';
import { From } from '../From';
import { SignBodyProps } from '../Sign';
import Intl from '@/i18n/i18n';
import { formatStandard, masToken, Unit } from '@/utils';

export function ExecuteSC(props: SignBodyProps) {
  const {
    MaxCoins,
    WalletAddress,
    OperationType,
    Nickname,
    Description: description,
    children,
  } = props;

  return (
    <div className="flex flex-col items-center gap-4 mas-menu-default">
      <From nickname={Nickname} walletAddress={WalletAddress} />

      <hr className="h-0.25 bg-neutral opacity-40 w-full" />

      <div className="flex w-full items-center justify-between">
        <p>{Intl.t('password-prompt.sign.operation-type')}</p>
        <p>{OperationType}</p>
      </div>

      <div className="flex w-full items-center justify-between">
        <p>{Intl.t('password-prompt.sign.max-coins')} </p>
        <p>
          {formatStandard(MaxCoins, Unit.NanoMAS)} {masToken}
        </p>
      </div>

      <hr className="h-0.25 bg-neutral opacity-40 w-full" />

      <Description description={description} />

      {children}
    </div>
  );
}
