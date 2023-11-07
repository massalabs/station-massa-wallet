import { SignBodyProps } from '../Sign';
import { Description } from '../SignComponentUtils/Description';
import { FromTo } from '../SignComponentUtils/FromTo';
import Intl from '@/i18n/i18n';
import { formatStandard, masToken, Unit } from '@/utils';

export function Transaction(props: SignBodyProps) {
  const {
    WalletAddress,
    RecipientAddress,
    RecipientNickname = '',
    OperationType,
    Amount,
    Description: description,
    Nickname,
    children,
  } = props;

  return (
    <div className="flex flex-col items-center gap-4 mas-menu-default">
      <FromTo
        fromNickname={Nickname}
        fromAddress={WalletAddress}
        toNickname={RecipientNickname}
        toAddress={RecipientAddress}
        label="password-prompt.sign.to"
      />

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
