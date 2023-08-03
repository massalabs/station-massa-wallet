import { FormEvent, useEffect, useState } from 'react';

import {
  Button,
  Input,
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
  toast,
} from '@massalabs/react-ui-kit';
import { FiPlus } from 'react-icons/fi';
import { useParams } from 'react-router-dom';

import { usePost } from '@/custom/api';
import Intl from '@/i18n/i18n';
import { ImportAssetsObject } from '@/models/assetModel';
import { IForm, parseForm } from '@/utils';

export interface InputsErrors {
  address?: string;
}

export function AssetsImportModal({ ...props }) {
  const { setModal } = props;

  const [tokenAddress, setTokenAddress] = useState<string>('');
  const [error, setError] = useState<InputsErrors | null>(null);
  const { nickname } = useParams();

  const {
    mutate,
    error: postError,
    isSuccess: postSuccess,
  } = usePost<ImportAssetsObject>(`accounts/${nickname}/assets`);

  function isValidAssetAddress(input: string): boolean {
    const regexPattern = /^AS[0-9a-zA-Z]+$/;
    return regexPattern.test(input);
  }

  function validate(formObject: IForm) {
    const { tokenAddress } = formObject;

    setError(null);

    if (isValidAssetAddress(tokenAddress) === false) {
      setError({ address: 'Invalid Address' });
      return false;
    }

    return true;
  }

  function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const formObject = parseForm(e);

    if (!validate(formObject)) return;

    if (formObject.tokenAddress) {
      mutate(formObject.tokenAddress);
    }

    setModal(false);
  }

  useEffect(() => {
    if (postError) {
      toast.error('error');
    } else if (postSuccess) {
      toast.success('success');
    }
  }, [postSuccess]);

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
          <form onSubmit={handleSubmit}>
            <Input
              placeholder={'Token Address'}
              value={tokenAddress}
              name="tokenAddress"
              onChange={(e) => setTokenAddress(e.target.value)}
              error={error?.address}
            />
            <Button
              customClass="mt-2"
              preIcon={<FiPlus size={24} />}
              type="submit"
            >
              {Intl.t('assets.add')}
            </Button>
          </form>
        </div>
      </PopupModalContent>
    </PopupModal>
  );
}
