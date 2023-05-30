import { AccountObject } from '../../../models/AccountModel';

interface TransferCoinProps {
  account: AccountObject;
}

export function TransferCoin({ account }: TransferCoinProps) {
  return <p className="mas-banner text-neutral">{account.nickname}</p>;
}
