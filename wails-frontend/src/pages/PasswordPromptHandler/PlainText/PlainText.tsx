import { AccordionCategory, AccordionContent } from '@massalabs/react-ui-kit';
import { WindowSetSize } from '@wailsjs/runtime/runtime';
import { FiChevronDown, FiChevronUp, FiInfo } from 'react-icons/fi';

import { From } from '../From';
import { SignBodyProps } from '../Sign';
import Intl from '@/i18n/i18n';

export function PlainText(props: SignBodyProps) {
  const { PlainText, DisplayData, WalletAddress, Description, Nickname } =
    props;

  const toAddInHeightDescription = Description ? 50 : 0;
  const toAddInHeightDisplayData = DisplayData ? 50 : -50;

  const winWidth = 460;
  const winHeight = 460 + toAddInHeightDescription + toAddInHeightDisplayData;

  WindowSetSize(winWidth, winHeight);

  return (
    <div className="flex flex-col items-center gap-4 mas-menu-default w-[326px]">
      <From nickname={Nickname} walletAddress={WalletAddress} />

      {Description && (
        <>
          <hr className="h-0.25 bg-neutral opacity-40 w-full" />

          <div className="flex flex-col w-full h-fit">
            <AccordionCategory
              isChild={false}
              iconOpen={<FiChevronDown />}
              iconClose={<FiChevronUp />}
              customClass="px-0 py-0"
              categoryTitle={
                <div className="flex items-center w-full gap-4">
                  <FiInfo size={18} />
                  <p>{Intl.t('password-prompt.sign.description')}</p>
                </div>
              }
            >
              <AccordionContent customClass="px-0 pt-4 pb-0">
                <div className="max-w-full overflow-hidden">
                  <p>{Description}</p>
                </div>
              </AccordionContent>
            </AccordionCategory>
          </div>
        </>
      )}

      {DisplayData && (
        <>
          <hr className="h-0.25 bg-neutral opacity-40 w-full" />

          <div className="flex flex-col w-full h-fit">
            <AccordionCategory
              isChild={false}
              iconOpen={<FiChevronDown />}
              iconClose={<FiChevronUp />}
              customClass="px-0 py-0"
              categoryTitle={
                <div className="flex items-center w-full gap-4">
                  <FiInfo size={18} />
                  <p>{Intl.t('password-prompt.sign.message')}</p>
                </div>
              }
            >
              <AccordionContent customClass="px-0 pt-4 pb-0">
                <div className="max-w-full overflow-hidden">
                  <p>{PlainText}</p>
                </div>
              </AccordionContent>
            </AccordionCategory>
          </div>
        </>
      )}
    </div>
  );
}
