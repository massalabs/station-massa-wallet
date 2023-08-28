import { WindowSetSize } from '@wailsjs/runtime/runtime';

import { PromptRequestData } from '../Sign';
import Intl from '@/i18n/i18n';

export function PlainByteText(props: PromptRequestData) {
  const { PlainText, OperationType } = props;

  WindowSetSize(460, 460);

  return (
    <div className="flex flex-col items-center gap-4 mas-menu-default">
      <div className="flex justify-between w-full">
        <p className="mas-body">
          {Intl.t('password-prompt.sign.message-format')}
          {OperationType}
        </p>
      </div>

      {/* If plain text is not null or undefiened show it, otherwise (byte text case) don't shwo anything */}
      {PlainText && (
        <div className="flex w-full items-center justify-between">
          <p className="mas-caption">{PlainText}</p>
        </div>
      )}
    </div>
  );
}
