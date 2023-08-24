import { WindowSetSize } from '@wailsjs/runtime/runtime';
import { FiArrowRight } from 'react-icons/fi';

import { PromptRequestData } from '../Sign';
import Intl from '@/i18n/i18n';
import { masToken, maskAddress } from '@/utils';

export function CallSc(props: PromptRequestData) {
  const {
    Coins,
    Address,
    WalletAddress,
    Function: CalledFunction,
    OperationType,
  } = props;

  WindowSetSize(545, 642);

  return (
    <div className="flex flex-col items-center gap-4">
      <div className="flex w-full items-center justify-between">
        <div className="flex flex-col">
          <p>{Intl.t('password-prompt.sign.from')}</p>
          <p className="mas-caption">{maskAddress(WalletAddress)}</p>
        </div>
        <div className="h-8 w-8 rounded-full flex items-center justify-center bg-neutral">
          <FiArrowRight size={24} className="text-primary" />
        </div>
        <div className="flex flex-col">
          <p>{Intl.t('password-prompt.sign.to')}</p>
          <p className="mas-caption">{maskAddress(Address)}</p>
        </div>
      </div>

      <hr className="h-0.25 bg-neutral opacity-40 w-full" />

      <div className="flex justify-between w-full">
        <div className="flex flex-col h-fit">
          <p>{Intl.t('password-prompt.sign.operation-type')}</p>
          <p className="mas-caption">
            {Intl.t('password-prompt.sign.function-called')}
          </p>
        </div>
        <div className="flex flex-col items-end h-fit">
          <p>{OperationType}</p>
          <p className="mas-caption">{CalledFunction}</p>
        </div>
      </div>

      <hr className="h-0.25 bg-neutral opacity-40 w-full" />

      <div className="flex flex-col gap-2 w-full">
        <div className="flex w-full items-center justify-between">
          <p>{Intl.t('password-prompt.sign.coins')}</p>
          <p>
            {Coins} {masToken}
          </p>
        </div>
      </div>

      <hr className="h-0.25 bg-neutral opacity-40 w-full" />
    </div>
  );
}
