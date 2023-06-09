import { Tabs } from '@massalabs/react-ui-kit';
import WalletLayout, {
  MenuItem,
} from '../../layouts/WalletLayout/WalletLayout';
import Intl from '../../i18n/i18n';
import Send from './SendCoins/Send';
import ReceiveCoins from './ReceiveCoins/ReceiveCoins';
import { useNavigate, useParams } from 'react-router-dom';
import { useResource } from '../../custom/api';
import { AccountObject } from '../../models/AccountModel';
import { routeFor } from '../../utils';
import { useQuery } from '../../custom/api/useQuery';

function SendCoins() {
  const navigate = useNavigate();
  const query = useQuery();

  const tabIndex = parseInt(query.get('tabIndex') || '0');

  const { nickname } = useParams();
  const { data: account } = useResource<AccountObject>(`accounts/${nickname}`);

  if (account === undefined) {
    navigate(routeFor(`${nickname}/home`));
    return null;
  }

  const tabsConfig = [
    {
      label: Intl.t('sendcoins.send-tab'),
      content: <Send account={account} />,
    },
    {
      label: Intl.t('sendcoins.receive-tab'),
      content: <ReceiveCoins />,
    },
  ];

  return (
    <WalletLayout menuItem={MenuItem.SendCoins}>
      <div className="w-1/2 h-1/2">
        <Tabs tabsConfig={tabsConfig} defaultIndex={tabIndex} />
      </div>
    </WalletLayout>
  );
}

export default SendCoins;
