import { Tooltip, formatAmount, massaToken } from '@massalabs/react-ui-kit';
import { FiAlertTriangle, FiInfo } from 'react-icons/fi';

import Intl from '@/i18n/i18n';
import { SignBodyProps } from '@/pages/PasswordPromptHandler/Sign';
import { Description } from '@/pages/PasswordPromptHandler/SignComponentUtils/Description';
import { From } from '@/pages/PasswordPromptHandler/SignComponentUtils/From';

export function ExecuteSC(props: SignBodyProps) {
  const {
    MaxCoins,
    WalletAddress,
    OperationType,
    Nickname,
    Description: description,
    children,
    DeployedByteCodeSize,
  } = props;

  const isDeploySC = DeployedByteCodeSize > 0;

  return (
    <div className="flex flex-col items-center gap-4">
      <From nickname={Nickname} walletAddress={WalletAddress} />

      <hr className="h-0.25 bg-neutral opacity-40 w-full" />

      <div className="flex w-full items-center justify-between">
        <p>{Intl.t('password-prompt.sign.operation-type')}</p>
        <div className="flex items-center">
          <p>
            {isDeploySC
              ? Intl.t('password-prompt.sign.deploy-sc-pseudo-operation-type')
              : Intl.t(`password-prompt.sign.operation-types.${OperationType}`)}
          </p>
          <Tooltip
            customClass="mas-caption !w-7/12 left-0  !ml-28"
            body={
              isDeploySC ? (
                Intl.t('password-prompt.sign.deploy-sc-tooltip')
              ) : (
                <>
                  {Intl.t('password-prompt.sign.execute-sc-tooltip.1')}
                  <br />
                  {Intl.t('password-prompt.sign.execute-sc-tooltip.2')}
                  <br />
                  <br />
                  {Intl.t('password-prompt.sign.execute-sc-tooltip.3')}
                </>
              )
            }
          >
            <FiInfo size={16} />
          </Tooltip>
        </div>
      </div>

      <div className="flex w-full items-center">
        <FiAlertTriangle
          size={42}
          className="text-s-warning pr-2"
          strokeWidth="2"
        />
        <p className="mas-caption">
          {isDeploySC ? (
            <>
              {Intl.t('password-prompt.sign.deploy-sc-warning.1')}
              <br />
              {Intl.t('password-prompt.sign.deploy-sc-warning.2')}
            </>
          ) : (
            <>
              {Intl.t('password-prompt.sign.execute-sc-warning.1')}
              <br />
              {Intl.t('password-prompt.sign.execute-sc-warning.2')}
            </>
          )}
        </p>
      </div>

      <hr className="h-0.25 bg-neutral opacity-40 w-full" />

      <Description text={description} />

      {!isDeploySC && (
        <div className="flex w-full items-center justify-between pb-2">
          <div className="flex items-center">
            <Tooltip
              className="w-fit pl-0 pr-2 hover:cursor-pointer"
              customClass="mas-caption !w-7/12"
              body={
                <>
                  {Intl.t(
                    'password-prompt.sign.execute-sc-max-coins-tooltip.1',
                  )}
                  <br />
                  {Intl.t(
                    'password-prompt.sign.execute-sc-max-coins-tooltip.2',
                  )}
                </>
              }
            >
              <FiInfo size={18} />
            </Tooltip>
            <p>{Intl.t('password-prompt.sign.max-coins')}</p>
          </div>
          <p>
            {formatAmount(MaxCoins).full} {massaToken}
          </p>
        </div>
      )}

      {children}
    </div>
  );
}
