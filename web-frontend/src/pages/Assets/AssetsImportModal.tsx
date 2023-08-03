import {
  Button,
  Input,
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
} from '@massalabs/react-ui-kit';

import Intl from '@/i18n/i18n';

import { FiPlus } from 'react-icons/fi';

export function AssetsImportModal({ ...props }) {
  const { setModal } = props;
  return (
    <PopupModal
      customClass="max-w-[500px] w-[500px]"
      fullMode={true}
      onClose={() => setModal(false)}
    >
      <PopupModalHeader>
        <div className="mas-title mb-6">{Intl.t('assets.import-title')}</div>
      </PopupModalHeader>
      <PopupModalContent>
        <div className="mas-body2 pb-10">
          <p className="mb-6">{Intl.t('assets.import-subtitle')}</p>
          <Input customClass="mb-2" />
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
