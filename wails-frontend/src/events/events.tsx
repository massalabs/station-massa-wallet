/* eslint-disable @typescript-eslint/no-explicit-any */
export type promptResult<T = null> = {
  Success: boolean;
  CodeMessage: string;
  Data: T;
};

export enum promptAction {
  deleteReq = 0,
  newPasswordReq = 1,
  signReq = 2,
  importReq = 3,
  backupReq = 4,
  tradeRollsReq = 6,
  unprotectReq = 7,
}

export type promptRequest = {
  Action: promptAction;
  Msg: string;
  Data: any;
  CodeMessage: string;
  CorrelationID: string;
};

export const events = {
  promptResult: 'promptResult',
  promptData: 'promptData',
  promptRequest: 'promptRequest',
};

export const backupMethods = {
  ymlFileBackup: 'yaml',
  privateKeyBackup: 'privateKey',
};
