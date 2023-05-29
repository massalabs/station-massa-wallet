export enum ErrorCode {
  InvalidNickname = 'Nickname-0001',
  InvalidPrivateKey = 'PrivateKey-0001',
  FilesystemError = 'AccountFile-0001',
  DuplicateKey = 'DuplicateKey-0001',
  UnknownError = 'Unknown-0001',
  DuplicateNickname = 'DuplicateNickname-001',
  TimeoutError = 'Timeout-0001',
  WrongPassword = 'WrongPassword-0001',
  InvalidPromptInput = 'InvalidPromptInput-0001',
}

const errorsEN: Record<ErrorCode, string> = {
  [ErrorCode.InvalidNickname]: 'Invalid nickname',
  [ErrorCode.InvalidPrivateKey]: 'Invalid private key',
  [ErrorCode.FilesystemError]: 'Filesystem error',
  [ErrorCode.DuplicateKey]: 'Private key already exists',
  [ErrorCode.UnknownError]: 'Unknown error, try again',
  [ErrorCode.DuplicateNickname]: 'This username already exists',
  [ErrorCode.TimeoutError]: 'Timeout error',
  [ErrorCode.WrongPassword]: 'Wrong password',
  [ErrorCode.InvalidPromptInput]: 'Invalid user input',
};

export function getErrorMessage(code: ErrorCode | string): string {
  if (typeof code === 'string' && !(code in ErrorCode)) {
    console.log('Unknown error code', code);
    return code;
  }

  const errorMessage = errorsEN[code as ErrorCode];
  if (errorMessage) {
    return errorMessage;
  }

  return code;
}

export interface IErrorObject {
  password?: string;
  error?: string;
}
