import { useState } from 'react';

import {
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
} from '@massalabs/react-ui-kit';


import { ImportForm } from '.';
import { ImportSuccess } from './ImportSuccess';
import Intl from '@/i18n/i18n';

export function AssetsImportModal({ ...props }) {
  const { closeModal } = props;

  const [importSuccess, setImportSuccess] = useState<boolean>(false);

  return (
    <>
      <PopupModal
        customClass={
          importSuccess ? `w-[513px] h-[440px]` : `w-[580px] h-[300px]`
        }
        fullMode={true}
        onClose={() => closeModal()}
        customClassNested={importSuccess ? `w-full h-full` : ``}
      >
        {importSuccess ? (
          <ImportSuccess closeModal={() => closeModal()} />
        ) : (
          <>
            <PopupModalHeader>
              <div className="flex flex-col">
                <div className="mas-title mb-6">
                  {Intl.t('assets.import.import-title')}
                </div>
                <p className="mb-6">
                  {Intl.t('assets.import.import-subtitle')}
                </p>
              </div>
            </PopupModalHeader>
            <PopupModalContent>
              <ImportForm
                {...props}
                setImportSuccess={() => setImportSuccess(true)}
              />
            </PopupModalContent>
          </>
        )}
      </PopupModal>
    </>
  );
}
