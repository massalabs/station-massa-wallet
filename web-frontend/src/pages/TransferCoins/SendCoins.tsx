import { Tabs } from '@massalabs/react-ui-kit';
import WalletLayout, {
  MenuItem,
} from '../../layouts/WalletLayout/WalletLayout';
import Intl from '../../i18n/i18n';
import Send from './SendCoins/Send';
import ReceiveCoins from './ReceiveCoins/ReceiveCoins';
import { useParams } from 'react-router-dom';
import { useResource } from '../../custom/api';
import { AccountObject } from '../../models/AccountModel';

function SendCoins() {
  const { nickname } = useParams();
  const { data: accounts = [] } = useResource<AccountObject[]>('accounts');
  const account = accounts.find((account) => account.nickname === nickname);
  const params = { account };

  const tabsConfig = [
    {
      label: Intl.t('sendcoins.send-tab'),
      content: <Send {...params} />,
    },
    {
      label: Intl.t('sendcoins.receive-tab'),
      content: <ReceiveCoins {...params} />,
    },
  ];

  return (
    <WalletLayout menuItem={MenuItem.SendCoins}>
      <div className="w-1/2 h-1/2">
        <Tabs tabsConfig={tabsConfig} defaultIndex={0} />
      </div>
    </WalletLayout>
  );
}

export default SendCoins;
