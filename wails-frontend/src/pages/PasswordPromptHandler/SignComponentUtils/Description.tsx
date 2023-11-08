import { AccordionCategory, AccordionContent } from '@massalabs/react-ui-kit';
import { FiChevronDown, FiChevronUp, FiInfo } from 'react-icons/fi';

import Intl from '@/i18n/i18n';

export interface DescriptionProps {
  text?: string;
  label?: string;
}

export function Description(props: DescriptionProps) {
  const { text, label = 'password-prompt.sign.description' } = props;

  if (!text) return null;

  return (
    <>
      <div className="flex flex-col w-full h-fit">
        <AccordionCategory
          isChild={false}
          iconOpen={<FiChevronDown />}
          iconClose={<FiChevronUp />}
          customClass="px-0 py-0"
          categoryTitle={
            <div className="flex items-center w-full gap-4">
              <FiInfo size={18} />
              <p>{Intl.t(label)}</p>
            </div>
          }
        >
          <AccordionContent customClass="px-0 pt-4 pb-0">
            <div className="max-w-full overflow-hidden">
              <p>{text}</p>
            </div>
          </AccordionContent>
        </AccordionCategory>
      </div>
      <hr className="h-0.25 bg-neutral opacity-40 w-full" />
    </>
  );
}
