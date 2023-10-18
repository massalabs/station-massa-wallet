export enum ErrorCode {
  InvalidNickname = 'Nickname-0001',
  InvalidPrivateKey = 'PrivateKey-0001',
  FilesystemError = 'AccountFile-0001',
  ErrNoFile = 'NoFile-0001',
  DuplicateKey = 'DuplicateKey-0001',
  UnknownError = 'Unknown-0001',
  DuplicateNickname = 'DuplicateNickname-001',
  TimeoutError = 'Timeout-0001',
  WrongPassword = 'WrongPassword-0001',
  InvalidPromptInput = 'InvalidPromptInput-0001',
}

export interface IErrorObject {
  nickname?: string;
  password?: string;
  privateKey?: string;
  error?: string;
}
