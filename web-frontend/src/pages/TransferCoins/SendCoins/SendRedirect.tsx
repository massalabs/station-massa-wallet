import { useNavigate, useSearchParams } from 'react-router-dom';
import { fetchAccounts, routeFor } from '@/utils';

export default function Redirect() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const provider = searchParams.get('provider');
  const { okAccounts: accounts, isLoading } = fetchAccounts();

  if (isLoading === false) {
    let nickname;

    if (provider) {
      const accountProvider = accounts?.find((acc) => acc.address === provider);

      if (accountProvider) {
        nickname = accountProvider?.nickname;
      }
    } else if (accounts?.length) {
      nickname = accounts?.shift()?.nickname;
    }

    if (nickname) {
      navigate(`${routeFor(`${nickname}/transfer-coins`)}?${searchParams}`);
    } else {
      navigate(routeFor('index'));
    }
  }

  return <div>Redirecting...</div>;
}
