import { toMASS } from './massaFormat';

export function handlePercent(
  amount = 0,
  percent: number,
  fees: string,
  balance: number,
) {
  let newAmount = amount * percent;
  const feesAsNumber = Number(fees);

  if (newAmount > balance - feesAsNumber)
    newAmount = Math.max(balance - feesAsNumber, 0);

  return toMASS(newAmount);
}
