import { useEffect, useState } from 'react';

import {
  Button,
  Input,
  PopupModalContent,
  toast,
  PopupModal,
} from '@massalabs/react-ui-kit';
import { AxiosError } from 'axios';
import { FiTrash2 } from 'react-icons/fi';
import { useParams } from 'react-router-dom';

import {
  IDeleteAssetsErrors,
  assetDeleteErrors,
  deleteConfirm,
} from '@/const/assets/assets';
import { useDelete, useResource } from '@/custom/api';
import Intl from '@/i18n/i18n';
import { IToken } from '@/models/AssetModel';

interface DeleteAssetModal {
  closeModal: () => void;
  tokenAddress: string;
}

export function DeleteAssetModal({ ...props }: DeleteAssetModal) {
  const { tokenAddress, closeModal } = props;
  const { nickname } = useParams();

  const { refetch: refetchAssets } = useResource<IToken[]>(
    `accounts/${nickname}/assets`,
  );

  const [deletePhrase, setDeletePhrase] = useState<string>('');
  const [error, setError] = useState<IDeleteAssetsErrors | null>(null);

  const {
    mutate: mutateDelete,
    isSuccess: isSuccessDelete,
    isError: isErrorDelete,
    error: errorDelete,
  } = useDelete(`accounts/${nickname}/assets?assetAddress=${tokenAddress}`);

  const axiosError = errorDelete as AxiosError;
  const deleteErrorStatus = axiosError?.response?.status;

  useEffect(() => {
    if (isSuccessDelete) {
      handleDeleteSuccess();
    } else if (isErrorDelete) {
      displayErrors(deleteErrorStatus);
    }
  }, [isSuccessDelete, isErrorDelete]);

  function displayErrors(postStatus: number | undefined) {
    switch (postStatus) {
      case assetDeleteErrors.badRequest:
        toast.error(Intl.t('assets.delete.bad-request'));
        break;
      case assetDeleteErrors.invalidAddress:
        toast.error(Intl.t('assets.delete.invalid-address'));
        break;
      case assetDeleteErrors.serverError:
        toast.error(Intl.t('assets.internal-server-error'));
        break;
      default:
        toast.error(Intl.t('assets.unkown-error'));
    }
  }

  function handleDeleteSuccess() {
    toast.success(Intl.t('assets.delete.success'));
    closeModal();
    refetchAssets();
  }

  function validate(phrase: string) {
    if (!phrase) {
      setError({ phrase: Intl.t('assets.delete.no-value') });
      return false;
    }

    if (phrase !== deleteConfirm) {
      setError({ phrase: Intl.t('assets.delete.wrong-value') });
      return false;
    }
    return true;
  }

  function confirmDelete() {
    if (!validate(deletePhrase)) return;

    mutateDelete(null);
  }
  return (
    <PopupModal
      customClass="w-[520px] h-[440px] "
      customClassNested="px-10 py-20"
      fullMode={true}
    >
      <PopupModalContent customClassContent="py-10 ">
        <div className="mb-4">
          <h1 className="mas-title mb-4">{Intl.t('assets.delete.title')}</h1>

          <h3 className="mas-body2">{Intl.t('assets.delete.subtitle')}</h3>
        </div>
        <div className="mb-4">
          <Input
            value={deletePhrase}
            name="provider"
            onChange={(e) => setDeletePhrase(e.target.value)}
            error={error?.phrase}
            placeholder={Intl.t('assets.delete.placeholder')}
          />
        </div>
        <div className="flex gap-4">
          <Button onClick={() => closeModal()}>
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
