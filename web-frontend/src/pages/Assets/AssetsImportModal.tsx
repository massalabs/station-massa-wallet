import {
  Button,
  Input,
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
} from '@massalabs/react-ui-kit';
import { FiPlus } from 'react-icons/fi';

import Intl from '@/i18n/i18n';

export function AssetsImportModal({ ...props }) {
  const { closeModal } = props;
  return (
    <PopupModal
      customClass="w-[580px] h-[300px]"
      fullMode={true}
      onClose={() => closeModal()}
    >
      <PopupModalHeader>
        <div className="mas-title mb-6">{Intl.t('assets.import-title')}</div>
      </PopupModalHeader>
      <PopupModalContent>
        <div className="mas-body2 pb-10">
          <p className="mb-6">{Intl.t('assets.import-subtitle')}</p>
          <Input customClass="mb-6" value="" placeholder="" />
          <Button
            type="submit"
            onClick={() => {
              console.log('submit');
            }}
            preIcon={<FiPlus size={24} />}
          >
            {Intl.t('assets.add')}
          </Button>
        </div>
      </PopupModalContent>
    </PopupModal>
  );
}
