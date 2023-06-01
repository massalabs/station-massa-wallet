export function formatStandard(num: number, maximumFractionDigits = 2) {
  const locale = localStorage.getItem('locale') || 'en-US';
  return num.toLocaleString(locale, {
    minimumFractionDigits: 2,
    maximumFractionDigits,
  });
}

export function reverseFormatStandard(str: string) {
  const formattedString = str.replace(/[^0-9.-]/g, ''); // Remove non-numeric characters
  return parseFloat(formattedString);
}
