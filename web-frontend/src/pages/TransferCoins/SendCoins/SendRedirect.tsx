import { useNavigate } from 'react-router-dom';
import { useResource } from '../../../custom/api';
import { AccountObject } from '../../../models/AccountModel';
import { routeFor } from '../../../utils';
import { useState } from 'react';
import Intl from '../../../i18n/i18n';
import { useQuery } from '../../../custom/api/useQuery';
// this page is used when a user use a link in the receive page
// the link is set to redirect to the send coin page either for a specific sender
// or for another account if it is not specified
function Redirect() {
  const navigate = useNavigate();
  const [nickname, setNickname] = useState<string>('');
  let error: string | null = null;
  const { data: accounts = [] } = useResource<AccountObject[]>('accounts');
  let query = useQuery();
  let provider = query.get('provider');
  if (provider) {
    // if the provider is in the query params, find the account with that address
    const accountProvider = accounts.find((acc) => acc.address === provider);
    if (accountProvider) {
      const nicknameProvider = accountProvider?.nickname;
      setNickname(nicknameProvider);
    } else {
      error = Intl.t('errors.send-redirect.account-not-found');
    }
  } else if (accounts.length > 0) {
    // if there is no provider in the query params, use the first account
    const nicknameProvider = accounts[0].nickname;
    setNickname(nicknameProvider);
  } else {
    error = Intl.t('errors.send-redirect.no-account');
  }
  if (nickname) {
    const newQueryParams = new URLSearchParams(location.search);
    newQueryParams.delete('provider');
    navigate(
      `${routeFor(`${nickname}/send-coins`)}?${newQueryParams.toString()}`,
    );
  }
  if (error) {
    navigate(routeFor('index'));
  }
  return (
    <div>
      Redirecting...
      {nickname && (
        <p>{`${Intl.t('send-redirect.redirecting', { nickname })} }`}</p>
      )}
      {error && <p>{error}</p>}
    </div>
  );
}
export default Redirect;
