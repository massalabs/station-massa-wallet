import { Tooltip } from '@massalabs/react-ui-kit';
import { FiAlertTriangle, FiInfo } from 'react-icons/fi';

import Intl from '@/i18n/i18n';
import { SignBodyProps } from '@/pages/PasswordPromptHandler/Sign';
import { Description } from '@/pages/PasswordPromptHandler/SignComponentUtils/Description';
import { From } from '@/pages/PasswordPromptHandler/SignComponentUtils/From';
import { Unit, formatStandard, masToken } from '@/utils';
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
    <div className="flex flex-col items-center gap-4">
      <From nickname={Nickname} walletAddress={WalletAddress} />

      <hr className="h-0.25 bg-neutral opacity-40 w-full" />

      <div className="flex w-full items-center justify-between">
        <p>{Intl.t('password-prompt.sign.operation-type')}</p>
        <div className="flex items-center">
          <p>
            {Intl.t(`password-prompt.sign.operation-types.${OperationType}`)}
          </p>
          <Tooltip
            icon={<FiInfo size={16} />}
            customClass="mas-caption !w-7/12 left-0  !ml-28"
            content={
              <>
                {Intl.t('password-prompt.sign.execute-sc-tooltip.1')}
                <br />
                {Intl.t('password-prompt.sign.execute-sc-tooltip.2')}
                <br />
                <br />
                {Intl.t('password-prompt.sign.execute-sc-tooltip.3')}
              </>
            }
          />
        </div>
      </div>

      <div className="flex w-full items-center">
        <FiAlertTriangle
          size={42}
          className="text-s-warning pr-2"
          strokeWidth="2"
        />
        <p className="mas-caption">
          {Intl.t('password-prompt.sign.execute-sc-warning.1')}
          <br />
          {Intl.t('password-prompt.sign.execute-sc-warning.2')}
        </p>
      </div>

      <hr className="h-0.25 bg-neutral opacity-40 w-full" />

      <Description description={description} />

      <div className="flex w-full items-center justify-between pb-2">
        <div className="flex items-center">
          <Tooltip
            icon={<FiInfo size={18} />}
            className="w-fit pl-0 pr-2 hover:cursor-pointer"
            customClass="mas-caption !w-7/12"
            content={
              <>
                {Intl.t('password-prompt.sign.execute-sc-max-coins-tooltip.1')}
                <br />
                {Intl.t('password-prompt.sign.execute-sc-max-coins-tooltip.2')}
              </>
            }
          />
          <p>{Intl.t('password-prompt.sign.max-coins')}</p>
        </div>
        <p>
          {formatStandard(MaxCoins, Unit.NanoMAS)} {masToken}
        </p>
      </div>

      {children}
    </div>
  );
}
