const errorsEN: Record<string, string> = {
  'Nickname-0001': 'Invalid nickname',
  'PrivateKey-0001': 'Invalid private key',
  'AccountFile-0001': 'Filesystem error',
  'DuplicateKey-0001': 'Private key already exists',
  'Unknown-0001': 'Unknown error, try again',
  'DuplicateNickname-001': 'This username already exists',
  'Timeout-0001': 'Timeout error',
  'WrongPassword-0001': 'Wrong password',
};

export function getErrorMessage(code: string): string {
  const errorMessage = errorsEN[code];
  if (errorMessage) {
    return errorMessage;
  }

  return code;
}
