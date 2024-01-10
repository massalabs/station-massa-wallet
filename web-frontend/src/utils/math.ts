import { toMASS } from './massaFormat';

export function handlePercent(
  amount = 0n,
  percent: bigint,
  fees: bigint,
  balance: bigint,
): number {
  let newAmount = (amount * percent) / 100n;

  if (newAmount > balance - fees) {
    if (balance - fees < 0) {
      newAmount = 0n;
    } else {
      newAmount = balance - fees;
    }
  }

  return toMASS(newAmount);
}
