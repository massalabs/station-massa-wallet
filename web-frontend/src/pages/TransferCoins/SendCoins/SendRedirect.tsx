import { useNavigate, useSearchParams } from 'react-router-dom';

import { fetchAccounts, routeFor } from '@/utils';

export default function Redirect() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { okAccounts: accounts, isLoading } = fetchAccounts();

  if (isLoading === false) {
    const nickname = accounts?.shift()?.nickname;

    if (nickname) {
      navigate(`${routeFor(`${nickname}/transfer-coins`)}?${searchParams}`);
    } else {
      navigate(routeFor('index'));
    }
  }

  return <div>Redirecting...</div>;
}
