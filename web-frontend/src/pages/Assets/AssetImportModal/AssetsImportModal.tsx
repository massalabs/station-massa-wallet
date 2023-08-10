import {
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
} from '@massalabs/react-ui-kit';

import { ImportForm } from './ImportForm';
import Intl from '@/i18n/i18n';

export function AssetsImportModal({ ...props }) {
  const { closeModal } = props;

  return (
    <PopupModal
      customClass="max-w-[580px] w-[580px]"
      fullMode={true}
      onClose={() => closeModal()}
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
