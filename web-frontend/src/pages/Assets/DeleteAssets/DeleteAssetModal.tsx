import { PopupModal, PopupModalHeader } from '@massalabs/react-ui-kit';

import Intl from '@/i18n/i18n';
import { ConfirmDelete } from '@/pages/Assets/DeleteAssets/ConfirmDelete';

interface DeleteAssetModal {
  closeModal: () => void;
  tokenAddress: string;
}

export function DeleteAssetModal({ ...props }: DeleteAssetModal) {
  const { closeModal } = props;
  return (
    <PopupModal
      customClass="!w-[440px] h-1/2"
      fullMode={true}
      onClose={() => closeModal()}
    >
      <PopupModalHeader>
        <div className="mb-4">
          <h1 className="mas-title mb-4">{Intl.t('assets.delete.title')}</h1>

          <h3 className="mas-body2">{Intl.t('assets.delete.subtitle')}</h3>
        </div>
      </PopupModalHeader>
      <ConfirmDelete {...props} />
    </PopupModal>
  );
}
