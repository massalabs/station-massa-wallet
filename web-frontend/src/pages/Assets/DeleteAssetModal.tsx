import { useEffect, useState } from 'react';

import {
  Button,
  Input,
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
  toast,
} from '@massalabs/react-ui-kit';
import { FiTrash2 } from 'react-icons/fi';
import { useParams } from 'react-router-dom';

import { DeleteAssetsErrors } from '@/const/assets/assets';
import { useDelete } from '@/custom/api';
import Intl from '@/i18n/i18n';
import { IToken } from '@/models/AccountModel';

export function DeleteAssetModal({ ...props }) {
  const { tokenAddress, setModalOpen } = props;
  const { nickname } = useParams();

  const [deletePhrase, setDeletePhrase] = useState<string>('');
  const [error, setError] = useState<DeleteAssetsErrors | null>(null);

  const {
    mutate: mutateDelete,
    isSuccess: isSuccessDelete,
    isError: isErrorDelete,
    error: errorDelete,
  } = useDelete<IToken[]>(
    `accounts/${nickname}/assets?assetAddress=${tokenAddress}`,
  );

  useEffect(() => {
    if (isSuccessDelete) {
      toast.success(Intl.t('assets.delete.success'));
      setModalOpen(false);
    } else if (isErrorDelete) {
      console.log(errorDelete);
    }
  }, [isSuccessDelete, isErrorDelete]);

  function validate(phrase: string) {
    if (!phrase) {
      setError({ phrase: Intl.t('assets.delete.no-value') });
      return false;
    }

    if (phrase !== 'DELETE') {
      setError({ phrase: Intl.t('assets.delete.wrong-value') });
      return false;
    }

    return true;
  }

  function confirmDelete() {
    if (!validate(deletePhrase)) return;

    mutateDelete({} as IToken[]);
  }

  return (
    <PopupModal customClass="!w-[440px] h-1/2" fullMode={true}>
      <PopupModalHeader>
        <div className="flex flex-col mb-4">
          <div className="mas-title mb-4">{Intl.t('assets.delete.title')}</div>

          <div className="mas-body2">{Intl.t('assets.delete.subtitle')}</div>
        </div>
      </PopupModalHeader>
      <PopupModalContent>
        <div className="mb-4">
          <Input
            value={deletePhrase}
            name="provider"
            onChange={(e) => setDeletePhrase(e.target.value)}
            error={error?.phrase}
            placeholder={Intl.t('assets.delete.placeholder')}
          />
        </div>
        <div className="flex gap-4 pb-12">
          <Button onClick={() => setModalOpen(false)}>
            {Intl.t('assets.delete.cancel')}
          </Button>
          <Button posIcon={<FiTrash2 />} onClick={() => confirmDelete()}>
            {Intl.t('assets.delete.delete')}
          </Button>
        </div>
      </PopupModalContent>
    </PopupModal>
  );
}
