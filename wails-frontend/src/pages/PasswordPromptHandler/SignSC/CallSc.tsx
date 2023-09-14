import { AccordionCategory, AccordionContent } from '@massalabs/react-ui-kit';
import { WindowSetSize } from '@wailsjs/runtime/runtime';
import {
  FiArrowRight,
  FiChevronDown,
  FiChevronUp,
  FiInfo,
} from 'react-icons/fi';

import { PromptRequestData } from '../Sign';
import Intl from '@/i18n/i18n';
import { formatStandard, masToken, maskAddress, Unit } from '@/utils';

export function CallSc(props: PromptRequestData) {
  const {
    Coins,
    Fees,
    Address,
    WalletAddress,
    Function: CalledFunction,
    OperationType,
    Description,
  } = props;

  const toAddInHeigthDescription = Description ? 200 : 0;

  const winWidth = 460;
  const winHeight = 460 + toAddInHeigthDescription;

  WindowSetSize(winWidth, winHeight);

  return (
    <div className="flex flex-col items-center gap-4 mas-menu-default w-[326px]">
      <div className="flex w-full items-center justify-between">
        <div className="flex flex-col">
          <p className="mas-menu-active">
            {Intl.t('password-prompt.sign.from')}
          </p>
          <p className="mas-caption">{maskAddress(WalletAddress)}</p>
        </div>
        <div className="h-8 w-8 rounded-full flex items-center justify-center bg-neutral">
          <FiArrowRight size={24} className="text-primary" />
        </div>
        <div className="flex flex-col">
          <p className="mas-menu-active">
            {Intl.t('password-prompt.sign.contract')}
          </p>
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

      {Description && (
        <div className="flex flex-col w-full h-fit">
          <AccordionCategory
            isChild={false}
            iconOpen={<FiChevronDown />}
            iconClose={<FiChevronUp />}
            customClass={'!p-0'}
            categoryTitle={
              <div className="flex items-center w-full gap-4">
                <FiInfo size={18} />
                <p>{Intl.t('password-prompt.sign.description')}</p>
              </div>
            }
          >
            <AccordionContent customClass={'px-0 pt-4 pb-0'}>
              <div className="max-w-full overflow-hidden">
                <p>{Description}</p>
              </div>
            </AccordionContent>
          </AccordionCategory>
        </div>
      )}

      {Description && <hr className="h-0.25 bg-neutral opacity-40 w-full" />}

      <div className="flex flex-col gap-2 w-full">
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

      <hr className="h-0.25 bg-neutral opacity-40 w-full" />
    </div>
  );
}
