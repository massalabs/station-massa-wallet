import { useEffect, useState } from 'react';

import {
  Button,
  Input,
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
} from '@massalabs/react-ui-kit';
import { AxiosError } from 'axios';
import { FiPlus } from 'react-icons/fi';
import { useParams } from 'react-router-dom';

import { InputsErrors, assetImportErrors } from '@/const/assets/assets';
import { usePost, useResource } from '@/custom/api';
import Intl from '@/i18n/i18n';
import { Asset } from '@/models/AssetModel';
import { ImportResult } from '@/pages/Assets';
import { isValidAssetAddress } from '@/validation/asset';

interface ImportModalProps {
  closeModal: () => void;
}

export function AssetsImportModal(props: ImportModalProps) {
  const { closeModal } = props;

  const [importResult, setImportResult] = useState<boolean | null>(null);
  const [inputError, setInputError] = useState<InputsErrors | null>(null);
  const [assetAddress, setAssetAddress] = useState<string>('');

  const { nickname } = useParams();

  const { refetch } = useResource<Asset[]>(`accounts/${nickname}/assets`);

  const {
    mutate,
    data,
    isSuccess: postSuccess,
    isError: postError,
    error,
  } = usePost<Asset>(
    `accounts/${nickname}/assets?assetAddress=${assetAddress}`,
  );
  const axiosError = error as AxiosError;
  const postErrorStatus = axiosError?.response?.status;

  useEffect(() => {
    if (postSuccess) {
      refetch();
      setImportResult(true);
    } else if (postError) {
      setImportResult(false);
      displayErrors(postErrorStatus);
    }
  }, [postSuccess, postError]);

  function displayErrors(postStatus: number | undefined) {
    switch (postStatus) {
      case assetImportErrors.badRequest:
        setInputError({
          address: Intl.t('assets.import.failure-screen.bad-request'),
        });
        break;
      case assetImportErrors.invalidAddress:
        setInputError({
          address: Intl.t('assets.import.failure-screen.invalid-address'),
        });
        break;
      case assetImportErrors.notFound:
        setInputError({
          address: Intl.t('assets.import.failure-screen.not-found'),
        });
        break;
      case assetImportErrors.serverError:
        setInputError({
          address: Intl.t('assets.import.failure-screen.internal-server-error'),
        });
        break;
      default:
        setInputError({
          address: Intl.t('assets.import.failure-screen.unkown-error'),
        });
    }
  }

  function validate() {
    setInputError(null);

    if (!assetAddress) {
      setInputError({
        address: Intl.t('assets.import.failure-screen.no-input'),
      });
      return false;
    }

    if (!isValidAssetAddress(assetAddress)) {
      setInputError({
        address: Intl.t('assets.import.failure-screen.wrong-format'),
      });
      return false;
    }
    return true;
  }

  function handleSubmit() {
    if (!validate()) return;
    handleImport();
  }

  function handleImport() {
    mutate({} as Asset);
  }

  return (
    <>
      {importResult !== null ? (
        <PopupModal
          customClass="w-[513px] h-[440px]"
          fullMode={true}
          onClose={() => closeModal()}
          customClassNested="w-full h-full"
        >
          <ImportResult
            closeModal={() => closeModal()}
            data={data as Asset}
            importResult={importResult}
            inputError={inputError}
          />
        </PopupModal>
      ) : (
        <PopupModal
          customClass="w-[580px] h-[300px]"
          fullMode={true}
          onClose={() => closeModal()}
        >
          <PopupModalHeader>
            <div className="mb-6">
              <p className="mas-title mb-6">
                {Intl.t('assets.import.import-title')}
              </p>
              <p className="mas-body2">
                {Intl.t('assets.import.import-subtitle')}
              </p>
            </div>
          </PopupModalHeader>
          <PopupModalContent>
            <div className="mas-body2 pb-10">
              <Input
                value={assetAddress}
                onChange={(e) => setAssetAddress(e.target.value)}
                name="tokenAddress"
                error={inputError?.address}
                placeholder={Intl.t('assets.import.placeholder')}
              />

              <Button
                customClass="mt-6"
                preIcon={<FiPlus size={24} />}
                onClick={() => handleSubmit()}
              >
                {Intl.t('assets.import.add')}
              </Button>
            </div>
          </PopupModalContent>
        </PopupModal>
      )}
    </>
  );
}
