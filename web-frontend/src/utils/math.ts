import { formatAmount } from '@massalabs/react-ui-kit';

import { MAS } from '@/const/assets/assets';

export function handlePercent(
  amount = 0n,
  percent: bigint,
  fees: bigint,
  balance: bigint,
  decimals: number,
  symbol: string,
): string {
  let newAmount = (amount * percent) / 100n;

  if (symbol === MAS) {
    if (newAmount > balance - fees) {
      if (balance - fees < 0) {
        newAmount = 0n;
      } else {
        newAmount = balance - fees;
      }
    }
  }

  return formatAmount(newAmount.toString(), decimals).full;
}
