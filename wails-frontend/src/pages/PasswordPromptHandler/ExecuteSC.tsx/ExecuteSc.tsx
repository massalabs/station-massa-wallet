import { WindowSetSize } from '@wailsjs/runtime/runtime';

import { PromptRequestData } from '../Sign';
import Intl from '@/i18n/i18n';
import { formatStandard, masToken, maskAddress, toMASString } from '@/utils';

export function ExecuteSC(props: PromptRequestData) {
  const { MaxCoins, WalletAddress, OperationType, Fees, MaxGas, Expiry } =
    props;

  WindowSetSize(460, 460);

  return (
    <div className="flex flex-col items-center gap-4 mas-menu-default">
      <div className="flex w-full items-center justify-between">
        <p>{Intl.t('password-prompt.sign.from')}</p>
        <p className="mas-caption">{maskAddress(WalletAddress)}</p>
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
          {formatStandard(toMASString(MaxCoins))} {masToken}
        </p>
      </div>
      <div className="flex w-full items-center justify-between">
        <p>{Intl.t('password-prompt.sign.max-gas')} </p>
        <p>
          {formatStandard(toMASString(MaxGas))} {masToken}
        </p>
      </div>
      <div className="flex w-full items-center justify-between">
        <p>{Intl.t('password-prompt.sign.expiry')} </p>
        <p>{Expiry}</p>
      </div>
      <div className="flex w-full items-center justify-between">
        <p>{Intl.t('password-prompt.sign.fees')}</p>
        <p>
          {formatStandard(toMASString(Fees))} {masToken}
        </p>
      </div>
    </div>
  );
}
