import { AxiosError } from 'axios';

import Intl from '@/i18n/i18n';

interface IBackendErrorObject {
  code?: string;
  message?: string;
}

type BackendErrors = {
  [key: string]: string;
};

const backendErrorsCode: BackendErrors = {
  'Nickname-0001': Intl.t('errors.invalid-nickname'),
  'WrongPassword-0001': Intl.t('errors.invalid-password'),
  'PrivateKey-0001': Intl.t('errors.invalid-private-key'),
  'NoFile-0001': Intl.t('errors.no-file'),
  'AccountFile-0001': Intl.t('errors.invalid-file'),
  'Unknown-0001': Intl.t('errors.unknown'),
  'DuplicateKey-0001': Intl.t('errors.duplicate-key'),
  'DuplicateNickname-001': Intl.t('errors.duplicate-nickname'),
  'Timeout-0001': Intl.t('errors.timeout'),
  'ActionCanceled-0001': Intl.t('errors.action-canceled'),
  PasswordRequired: Intl.t('errors.password-required'),
  'Wallet-0002': Intl.t('errors.unknown'),
};

export function parseErrors(err: AxiosError) {
  const data: IBackendErrorObject = err.response?.data as IBackendErrorObject;

  return (
    backendErrorsCode[data.code || 'Unknown-0001'] || Intl.t('errors.unknown')
  );
}
