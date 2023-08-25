import { WindowSetSize } from '@wailsjs/runtime/runtime';

import { PromptRequestData } from '../Sign';
import Intl from '@/i18n/i18n';
import { masToken, maskAddress } from '@/utils';

export function ExecuteSC(props: PromptRequestData) {
  const { MaxCoins, WalletAddress, OperationType } = props;

  WindowSetSize(470, 470);

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
          {MaxCoins} {masToken}
        </p>
      </div>
    </div>
  );
}
