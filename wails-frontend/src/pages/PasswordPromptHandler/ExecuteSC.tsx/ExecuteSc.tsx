import { WindowSetSize } from '@wailsjs/runtime/runtime';

import { PromptRequestData } from '../Sign';
import Intl from '@/i18n/i18n';
import { formatStandard, masToken, maskAddress, Unit } from '@/utils';

export function ExecuteSC(props: PromptRequestData) {
  const {
    MaxCoins,
    WalletAddress,
    OperationType,
    Fees,
    MaxGas,
    Expiry,
    Nickname,
  } = props;

  WindowSetSize(460, 560);

  return (
    <div className="flex flex-col items-center gap-4 mas-menu-default">
      <div className="flex w-full justify-between">
        <p>{Intl.t('password-prompt.sign.from')}</p>
        <div className="flex flex-col">
          <p className="mas-menu-default">{Nickname}</p>
          <p className="mas-caption">{maskAddress(WalletAddress)}</p>
        </div>
      </div>
      <hr className="h-0.25 bg-neutral opacity-40 w-full" />

      <div className="flex w-full items-center justify-between">
        <p>{Intl.t('password-prompt.sign.operation-type')}</p>
        <p>{OperationType}</p>
      </div>
      <hr className="h-0.25 bg-neutral opacity-40 w-full" />
      <div className="flex w-full items-center justify-between">
        <p>{Intl.t('password-prompt.sign.max-coins')} </p>
        <p>
          {formatStandard(MaxCoins, Unit.NanoMAS)} {masToken}
        </p>
      </div>
      <div className="flex w-full items-center justify-between">
        <p>{Intl.t('password-prompt.sign.max-gas')} </p>
        <p>{formatStandard(MaxGas)}</p>
      </div>
      <div className="flex w-full items-center justify-between">
        <p>{Intl.t('password-prompt.sign.expiry')} </p>
        <p>{Expiry}</p>
      </div>
      <div className="flex w-full items-center justify-between">
        <p>{Intl.t('password-prompt.sign.fees')}</p>
        <p>
          {formatStandard(Fees, Unit.NanoMAS)} {masToken}
        </p>
      </div>
    </div>
  );
}
