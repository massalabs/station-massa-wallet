import { useResource } from '@/custom/api';
import { AccountObject } from '@/models/AccountModel';

function separateAccounts(accounts: AccountObject[] | undefined) {
  if (!accounts) return { okAccounts: [], otherAccounts: [] };

  const { okAccounts, corruptedAccounts } = accounts.reduce<{
    okAccounts: AccountObject[];
    corruptedAccounts: AccountObject[];
  }>(
    (acc, account) => {
      if (account.status === 'ok') {
        acc.okAccounts.push(account);
      } else {
        acc.corruptedAccounts.push(account);
      }
      return acc;
    },
    { okAccounts: [], corruptedAccounts: [] },
  );

  return { okAccounts, corruptedAccounts };
}

export function useFetchAccounts() {
  const getAccounts = useResource<AccountObject[]>('accounts');
  const { data: accounts } = getAccounts;
  const { okAccounts, corruptedAccounts } = separateAccounts(accounts);
  return { okAccounts, corruptedAccounts, ...getAccounts };
}
