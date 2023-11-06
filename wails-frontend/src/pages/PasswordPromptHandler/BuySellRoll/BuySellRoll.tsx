import { Description } from '../Description';
import { From } from '../From';
import { SignBodyProps } from '../Sign';
import Intl from '@/i18n/i18n';
import { formatStandard, masToken, Unit } from '@/utils';

export function BuySellRoll(props: SignBodyProps) {
  const {
    Nickname,
    WalletAddress,
    RollCount,
    OperationType,
    Coins,
    children,
    Description: description,
  } = props;

  return (
    <div className="flex flex-col items-center gap-4 mas-menu-default">
      <From nickname={Nickname} walletAddress={WalletAddress} />

      <div className="flex justify-between w-full">
        <p>{Intl.t('password-prompt.sign.operation-type')}</p>
        <p>{OperationType}</p>
      </div>

      <hr className="h-0.25 bg-neutral opacity-40 w-full" />

      <div className="flex justify-between w-full">
        <p>{Intl.t('password-prompt.sign.roll-amount')}</p>
        <p className="mas-menu-default">{RollCount}</p>
      </div>

      <div className="flex w-full items-center justify-between">
        <p>{Intl.t('password-prompt.sign.spend-amount')}</p>
        {/* or receive amount */}
        <p>
          {formatStandard(Coins, Unit.NanoMAS)} {masToken}
        </p>
      </div>

      <hr className="h-0.25 bg-neutral opacity-40 w-full" />

      <Description description={description} />

      {children}
    </div>
  );
}
