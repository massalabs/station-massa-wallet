import { useEffect, useState } from 'react';

import { Button, Input, toast } from '@massalabs/react-ui-kit';
import { AxiosError } from 'axios';
import { FiPlus } from 'react-icons/fi';
import { useParams } from 'react-router-dom';

import { InputsErrors, assetImportErrors } from '@/const/assets/assets';
import { usePost, useResource } from '@/custom/api';
import Intl from '@/i18n/i18n';
import { Asset } from '@/models/AssetModel';

export function ImportForm({ ...props }) {
  const { closeModal } = props;
  const { nickname } = useParams();

  const { refetch } = useResource<Asset[]>(`accounts/${nickname}/assets`);

  const [inputError, setInputError] = useState<InputsErrors | null>(null);
  const [tokenAddress, setTokenAddress] = useState<string>('');

  const {
    mutate,
    isSuccess: postSuccess,
    isError: postError,
    error,
  } = usePost<Asset>(
    `accounts/${nickname}/assets?assetAddress=${tokenAddress}`,
  );

  useEffect(() => {
    if (postSuccess) {
      toast.success(Intl.t('assets.import.success'));
      refetch();
      closeModal();
    } else if (postError) {
      displayErrors(postErrorStatus);
    }
  }, [postSuccess, postError]);

  const axiosError = error as AxiosError;
  const postErrorStatus = axiosError?.response?.status;

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

    handleImport();
  }

  function handleImport() {
    mutate({} as Asset);
  }

  function displayErrors(postStatus: number | undefined) {
    switch (postStatus) {
      case assetImportErrors.badRequest:
        toast.error(Intl.t('assets.import.bad-request'));
        break;
      case assetImportErrors.invalidAddress:
        setInputError({ address: Intl.t('assets.import.invalid-address') });
        break;
      case assetImportErrors.notFound:
        toast.error(Intl.t('assets.import.not-found'));
        break;
      case assetImportErrors.serverError:
        toast.error(Intl.t('assets.import.internal-server-error'));
        break;
      default:
        toast.error(Intl.t('assets.import.unkown-error'));
    }
  }

  return (
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
  );
}
