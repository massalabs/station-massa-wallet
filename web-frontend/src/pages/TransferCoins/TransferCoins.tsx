import { useEffect } from 'react';
import { Tabs } from '@massalabs/react-ui-kit';
import WalletLayout, {
  MenuItem,
} from '../../layouts/WalletLayout/WalletLayout';
import Intl from '../../i18n/i18n';
import SendCoins from './SendCoins/SendCoins';
import ReceiveCoins from './ReceiveCoins/ReceiveCoins';
import { useNavigate, useParams, useSearchParams } from 'react-router-dom';
import { useResource } from '../../custom/api';
import { AccountObject } from '../../models/AccountModel';
import { routeFor } from '../../utils';
import { TAB_SEND, TAB_RECEIVE } from '../../const/tabs/tabs';

function TransferCoins() {
  const navigate = useNavigate();

  const [searchParams] = useSearchParams();

  const tabName = searchParams.get('tab') || TAB_SEND;
  const tabIndex = tabName === TAB_RECEIVE ? 1 : 0;

  const { nickname } = useParams();
  const { data: account, error } = useResource<AccountObject>(
    `accounts/${nickname}`,
  );

  window.history.replaceState(null, '', `/${nickname}/transfer-coins`);

  useEffect(() => {
    if (error) {
      navigate('/error');
    } else if (!account) {
      navigate(routeFor(`${nickname}/home`));
    }
  }, [account, error, navigate]);

  const tabsConfig = [
    {
      label: Intl.t('transfer-coins.send-tab'),
      content: <SendCoins account={account} />,
    },
    {
      label: Intl.t('transfer-coins.receive-tab'),
      content: <ReceiveCoins account={account} />,
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
