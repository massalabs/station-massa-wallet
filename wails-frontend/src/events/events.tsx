/* eslint-disable @typescript-eslint/no-explicit-any */
export type promptResult = {
  Success: boolean;
  CodeMessage: string;
};

export enum promptAction {
  deleteReq = 0,
  newPasswordReq = 1,
  signReq = 2,
  importReq = 3,
  backupReq = 4,
  transferReq = 5,
  tradeRollsReq = 6,
}

export type promptRequest = {
  Action: promptAction;
  Msg: string;
  Data: any;
};

export const events = {
  promptResult: 'promptResult',
  promptRequest: 'promptRequest',
};
