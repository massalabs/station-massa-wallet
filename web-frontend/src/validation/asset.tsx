export function isValidAssetAddress(input: string): boolean {
  const regexPattern = /^AS[0-9a-zA-Z]+$/;
  return regexPattern.test(input);
}
