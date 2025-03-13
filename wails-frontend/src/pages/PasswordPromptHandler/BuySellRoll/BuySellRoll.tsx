import { formatAmount, massaToken } from '@massalabs/react-ui-kit';

import { OPER_BUY_ROLL } from '@/const/operations';
import Intl from '@/i18n/i18n';
import { SignBodyProps } from '@/pages/PasswordPromptHandler/Sign';
import { Description } from '@/pages/PasswordPromptHandler/SignComponentUtils/Description';
import { From } from '@/pages/PasswordPromptHandler/SignComponentUtils/From';

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

  const label =
    OperationType === OPER_BUY_ROLL
      ? 'password-prompt.sign.spend-amount'
      : 'password-prompt.sign.receive-amount';

  return (
    <div className="flex flex-col items-center gap-4 mas-menu-default">
      <From nickname={Nickname} walletAddress={WalletAddress} />

      <div className="flex justify-between w-full">
        <p>{Intl.t('password-prompt.sign.operation-type')}</p>
        <p>{Intl.t(`password-prompt.sign.operation-types.${OperationType}`)}</p>
      </div>

      <hr className="h-0.25 bg-neutral opacity-40 w-full" />

      <div className="flex justify-between w-full">
        <p>{Intl.t('password-prompt.sign.roll-amount')}</p>
        <p className="mas-menu-default">{RollCount}</p>
      </div>

      <div className="flex w-full items-center justify-between">
        <p>{Intl.t(label)}</p>
        <p>
          {formatAmount(Coins).full} {massaToken}
        </p>
      </div>

      <hr className="h-0.25 bg-neutral opacity-40 w-full" />

      <Description text={description} />

      {children}
    </div>
  );
}
