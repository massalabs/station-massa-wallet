import { WindowSetSize } from '@wailsjs/runtime/runtime';

import { PromptRequestData } from '../Sign';
import Intl from '@/i18n/i18n';
import { formatStandard, masToken, Unit } from '@/utils';

export function BuySellRoll(props: PromptRequestData) {
  const { RollCount, OperationType, Coins, Fees } = props;

  WindowSetSize(460, 460);

  return (
    <div className="flex flex-col items-center gap-4 mas-menu-default">
      <div className="flex justify-between w-full">
        <p>{Intl.t('password-prompt.sign.operation-type')}</p>
        <p>{OperationType}</p>
      </div>

      <hr className="h-0.25 bg-neutral opacity-40 w-full" />

      <div className="flex justify-between w-full">
        <p>{Intl.t('password-prompt.sign.roll-amount')}</p>
        <p className="mas-menu-default">{RollCount}</p>
      </div>

      <hr className="h-0.25 bg-neutral opacity-40 w-full" />

      <div className="flex w-full items-center justify-between">
        <p>{Intl.t('password-prompt.sign.coins')}</p>
        <p>
          {formatStandard(Coins, Unit.NanoMAS)} {masToken}
        </p>
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
