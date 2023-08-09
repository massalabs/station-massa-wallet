import {
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
} from '@massalabs/react-ui-kit';

import { ImportForm } from './ImportForm';
import Intl from '@/i18n/i18n';

export function AssetsImportModal({ ...props }) {
  const { setModal } = props;

  return (
    <PopupModal
      customClass="max-w-[500px] w-[500px]"
      fullMode={true}
      onClose={() => setModal(false)}
    >
      <PopupModalHeader>
        <div className="flex flex-col">
          <div className="mas-title mb-6">
            {Intl.t('assets.import.import-title')}
          </div>
          <p className="mb-6">{Intl.t('assets.import.import-subtitle')}</p>
        </div>
      </PopupModalHeader>
      <PopupModalContent>
        <ImportForm {...props} />
      </PopupModalContent>
    </PopupModal>
  );
}
