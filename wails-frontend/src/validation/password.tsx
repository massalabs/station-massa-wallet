export function hasMoreThanFiveChars(password: string) {
  return password.length >= 5;
}

export function hasSamePassword(password: string, passwordConfirm: string) {
  return password === passwordConfirm;
}
