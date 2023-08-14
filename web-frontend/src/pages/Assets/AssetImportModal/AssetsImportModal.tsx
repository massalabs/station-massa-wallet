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

import { ImportResult } from './ImportModalScreens/ImportResult';
import { InputsErrors, assetImportErrors } from '@/const/assets/assets';
import { usePost, useResource } from '@/custom/api';
import Intl from '@/i18n/i18n';
import { IToken } from '@/models/AssetModel';


export function AssetsImportModal({ ...props }) {
  const { closeModal } = props;

  const [importResult, setImportResult] = useState<boolean | null>(null);
  const [inputError, setInputError] = useState<InputsErrors | null>(null);
  const [tokenAddress, setTokenAddress] = useState<string>('');

  const { nickname } = useParams();

  const { refetch } = useResource<IToken[]>(`accounts/${nickname}/assets`);

  const {
    mutate,
    data,
    isSuccess: postSuccess,
    isError: postError,
    error,
  } = usePost<IToken>(
    `accounts/${nickname}/assets?assetAddress=${tokenAddress}`,
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
  function isValidAssetAddress(input: string): boolean {
    const regexPattern = /^AS[0-9a-zA-Z]+$/;
    return regexPattern.test(input);
  }

  function validate() {
    setInputError(null);

    if (!tokenAddress) {
      setInputError({ address: Intl.t('assets.import.no-input') });
      return false;
    }

    if (!isValidAssetAddress(tokenAddress)) {
      setInputError({ address: Intl.t('assets.import.wrong-format') });
      return false;
    }
    return true;
  }

  function handleSubmit() {
    if (!validate()) return;
    setTokenAddress(tokenAddress);
    handleImport();
  }

  function handleImport() {
    mutate({} as IToken);
  }

  return (
    <>
      <PopupModal
        customClass={
          importResult !== null ? `w-[513px] h-[440px]` : `w-[580px] h-[300px]`
        }
        fullMode={true}
        onClose={() => closeModal()}
        customClassNested={importResult !== null ? `w-full h-full` : ``}
      >
        {importResult !== null ? (
          <ImportResult
            closeModal={() => closeModal()}
            token={data}
            importResult={importResult}
            inputError={inputError}
          />
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
              {/* <ImportForm
                {...props}
                setTokenAddress={setTokenAddress}
                tokenAddress={tokenAddress}
                mutate={mutate}
                setInputError={setInputError}
                inputError={inputError}
              /> */}
              <div className="mas-body2 pb-10">
                <Input
                  value={tokenAddress}
                  onChange={(e) => setTokenAddress(e.target.value)}
                  name="tokenAddress"
                  error={inputError?.address}
                  placeholder={Intl.t('assets.import.placeholder')}
                />
                <Button
                  customClass="mt-2 mt-6"
                  preIcon={<FiPlus size={24} />}
                  onClick={() => handleSubmit()}
                >
                  {Intl.t('assets.import.add')}
                </Button>
              </div>
            </PopupModalContent>
          </>
        )}
      </PopupModal>
    </>
  );
}
