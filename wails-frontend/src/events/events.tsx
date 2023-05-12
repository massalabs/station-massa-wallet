export type promptResult = {
  Success: boolean;
  Data: [];
  Error: string;
};

export enum promptAction {
  deleteReq = 0,
  newPasswordReq = 1,
  signReq = 2,
  importReq = 3,
  exportReq = 4,
}

export type promptRequest = {
  Action: promptAction;
  Msg: string;
  Data: [];
};

export const events = {
  passwordResult: 'passwordResult',
  promptRequest: 'promptRequest',
};
