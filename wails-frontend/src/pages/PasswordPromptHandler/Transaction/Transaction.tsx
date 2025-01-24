import { formatAmount, massaToken } from '@massalabs/react-ui-kit';

import Intl from '@/i18n/i18n';
import { AmountBox } from '@/pages/PasswordPromptHandler/AmountBox';
import { SignBodyProps } from '@/pages/PasswordPromptHandler/Sign';
import { Description } from '@/pages/PasswordPromptHandler/SignComponentUtils/Description';
import { FromTo } from '@/pages/PasswordPromptHandler/SignComponentUtils/FromTo';

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

      <AmountBox>
        {formatAmount(Amount).full} {massaToken}
      </AmountBox>

      <hr className="h-0.25 bg-neutral opacity-40 w-full" />

      <Description text={description} />

      {children}
    </div>
  );
}
