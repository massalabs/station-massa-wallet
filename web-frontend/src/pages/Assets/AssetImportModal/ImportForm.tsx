import { FormEvent, useEffect, useState } from 'react';

import { Button, Input, toast } from '@massalabs/react-ui-kit';
import { AxiosError } from 'axios';
import { FiPlus } from 'react-icons/fi';
import { useParams } from 'react-router-dom';

import { InputsErrors, assetImportErrors } from '@/const/assets/assets';
import { usePost } from '@/custom/api';
import Intl from '@/i18n/i18n';
import { ImportAssetsObject } from '@/models/ImportAssetModel';
import { IForm, parseForm } from '@/utils';

export function ImportForm({ ...props }) {
  const [inputError, setInputError] = useState<InputsErrors | null>(null);
  const [tokenAddress, setTokenAddress] = useState<string>('');

  const { setModal, refetch } = props;
  const { nickname } = useParams();
  const {
    mutate,
    isSuccess: postSuccess,
    isError: postError,
    error,
  } = usePost<ImportAssetsObject>(
    `accounts/${nickname}/assets?assetAddress=${tokenAddress}`,
  );

  useEffect(() => {
    if (postSuccess) {
      toast.success(Intl.t('assets.success'));
      refetch();
      setModal(false);
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

  function validate(formObject: IForm) {
    const { tokenAddress } = formObject;

    setInputError(null);

    if (!isValidAssetAddress(tokenAddress)) {
      setInputError({ address: Intl.t('assets.import.wrong-format') });
      return false;
    }
    return true;
  }

  function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const formObject = parseForm(e);

    if (!validate(formObject)) return;

    if (formObject.tokenAddress) {
      setTokenAddress(formObject.tokenAddress);
      handleImport();
    }
  }

  function handleImport() {
    mutate({} as ImportAssetsObject);
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
        console.log('Unknown Error:', postErrorStatus);
    }
  }

  return (
    <div className="mas-body2 pb-10">
      <form onSubmit={handleSubmit}>
        <Input
          placeholder={'Token Address'}
          default-value=""
          name="tokenAddress"
          error={inputError?.address}
        />
        <Button customClass="mt-2" preIcon={<FiPlus size={24} />} type="submit">
          {Intl.t('assets.import.add')}
        </Button>
      </form>
    </div>
  );
}
