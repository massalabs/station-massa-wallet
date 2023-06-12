import { Tabs } from '@massalabs/react-ui-kit';
import WalletLayout, {
  MenuItem,
} from '../../layouts/WalletLayout/WalletLayout';
import Intl from '../../i18n/i18n';
import SendCoins from './SendCoins/SendCoins';
import ReceiveCoins from './ReceiveCoins/ReceiveCoins';
import { useNavigate, useParams } from 'react-router-dom';
import { useResource } from '../../custom/api';
import { AccountObject } from '../../models/AccountModel';
import { routeFor } from '../../utils';
import { useQuery } from '../../custom/api/useQuery';
import { TAB_SEND, TAB_RECEIVE } from '../../const/tabs/tabs';

function TransferCoins() {
  const navigate = useNavigate();
  const query = useQuery();

  const tabName = query.get('tab') || TAB_SEND;
  const tabIndex = tabName === TAB_RECEIVE ? 1 : 0;

  const { nickname } = useParams();
  const { data: account } = useResource<AccountObject>(`accounts/${nickname}`);

  if (account === undefined) {
    navigate(routeFor(`${nickname}/home`));
    return null;
  }

  const tabsConfig = [
    {
      label: Intl.t('transfer-coins.send-tab'),
      content: <SendCoins account={account} />,
    },
    {
      label: Intl.t('transfer-coins.receive-tab'),
      content: <ReceiveCoins />,
    },
  ];

  return (
    <WalletLayout menuItem={MenuItem.TransferCoins}>
      <div className="w-1/2 h-1/2">
        <Tabs tabsConfig={tabsConfig} defaultIndex={tabIndex} />
      </div>
    </WalletLayout>
  );
}

export default TransferCoins;
