/* eslint-disable @typescript-eslint/no-explicit-any */
export type promptResult = {
  Success: boolean;
  Data: any;
  Error: string;
};

export enum promptAction {
  deleteReq = 0,
  newPasswordReq = 1,
  signReq = 2,
  importReq = 3,
  exportReq = 4,
  transferReq = 5,
}

export type promptRequest = {
  Action: promptAction;
  Msg: string;
  Data: any;
};

export const events = {
  passwordResult: 'passwordResult',
  promptRequest: 'promptRequest',
};
