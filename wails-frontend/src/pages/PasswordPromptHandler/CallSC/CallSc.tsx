import { FiAlertTriangle } from 'react-icons/fi';

import { FTTransferInfo } from './FTTransferInfo';
import Intl from '@/i18n/i18n';
import { SignBodyProps } from '@/pages/PasswordPromptHandler/Sign';
import { Description } from '@/pages/PasswordPromptHandler/SignComponentUtils/Description';
import { FromTo } from '@/pages/PasswordPromptHandler/SignComponentUtils/FromTo';

export function CallSc(props: SignBodyProps) {
  const {
    Address,
    WalletAddress,
    Function: CalledFunction,
    OperationType,
    Description: description,
    Nickname,
    Assets,
    Parameters,
    children,
    ExpiredSignRule,
    EnabledSignRule,
  } = props;

  const asset = Assets.find((a) => a.address === Address);

  return (
    <div className="flex flex-col items-center gap-4 mas-menu-default w-[326px]">
      {ExpiredSignRule && EnabledSignRule && (
        <div className="p-4 flex items-center mb-4">
          <FiAlertTriangle size={42} className="text-s-warning" />
          <div className="ml-2 text-s-warning text-sm">
            <p>
              {Intl.t('signRule.expiredSignRule', {
                signRuleType: EnabledSignRule,
              })}
            </p>
          </div>
        </div>
      )}

      <FromTo
        fromNickname={Nickname}
        fromAddress={WalletAddress}
        toAddress={Address}
        label="password-prompt.sign.contract"
      />

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

      <FTTransferInfo
        targetFunction={CalledFunction}
        asset={asset}
        parameters={Parameters}
      />

      <hr className="h-0.25 bg-neutral opacity-40 w-full" />

      <Description text={description} />

      {children}
    </div>
  );
}
