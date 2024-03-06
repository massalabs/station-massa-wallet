export enum ErrorCode {
  WrongPassword = 'WrongPassword-0001',
  Timeout = 'Timeout-0001',
}

export interface IErrorObject {
  nickname?: string;
  password?: string;
  privateKey?: string;
  error?: string;
}
