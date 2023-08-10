import { PopupModal, PopupModalHeader } from '@massalabs/react-ui-kit';

import { ConfirmDelete } from './ConfirmDelete';
import Intl from '@/i18n/i18n';


export function DeleteAssetModal({ ...props }) {
  const { setModalOpen } = props;

  return (
    <PopupModal
      customClass="!w-[440px] h-1/2"
      fullMode={true}
      onClose={() => setModalOpen(false)}
    >
      <PopupModalHeader>
        <div className="flex flex-col mb-4">
          <div className="mas-title mb-4">{Intl.t('assets.delete.title')}</div>

          <div className="mas-body2">{Intl.t('assets.delete.subtitle')}</div>
        </div>
      </PopupModalHeader>
      <ConfirmDelete {...props} />
    </PopupModal>
  );
}
