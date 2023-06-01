export function formatStandard(num: number, maximumFractionDigits = 2) {
  return num.toLocaleString('en-US', {
    minimumFractionDigits: 2,
    maximumFractionDigits,
  });
}

export function reverseFormatStandard(str: string) {
  const formattedString = str.replace(/[^0-9.-]/g, ''); // Remove non-numeric characters
  return parseFloat(formattedString);
}
