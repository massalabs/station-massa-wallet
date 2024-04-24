export function base64ToArray(base64: string): number[] {
  // Decode Base64 string to a binary string
  const binaryString = atob(base64);

  // Convert binary string to an array of numbers (byte values)
  const len = binaryString.length;
  const numbers = new Array<number>(len);
  for (let i = 0; i < len; i++) {
    numbers[i] = binaryString.charCodeAt(i);
  }

  return numbers;
}
