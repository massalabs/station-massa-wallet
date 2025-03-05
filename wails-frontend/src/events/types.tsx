import { walletapp } from '@wailsjs/go/models';

/* eslint-disable @typescript-eslint/no-explicit-any */
export type promptResult<T = null> = {
  Success: boolean;
  CodeMessage: string;
  Data: T;
};

export type promptRequest = {
  Action: walletapp.PromptRequestAction;
  Msg: string;
  Data: any;
  CodeMessage: string;
  DisablePassword: boolean | null;
};

export const backupMethods = {
  ymlFileBackup: 'yaml',
  privateKeyBackup: 'privateKey',
};
