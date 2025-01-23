export function base64ToArray(base64: string): Uint8Array {
  // Decode Base64 string to a binary string
  const binaryString = atob(base64);

  // Convert binary string to a Uint8Array
  const len = binaryString.length;
  const uintArray = new Uint8Array(len);
  for (let i = 0; i < len; i++) {
    uintArray[i] = binaryString.charCodeAt(i);
  }

  return uintArray;
}
