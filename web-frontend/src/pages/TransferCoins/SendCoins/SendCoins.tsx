import { Tabs, Button } from '@massalabs/react-ui-kit';
import WalletLayout, {
  MenuItem,
} from '../../../layouts/WalletLayout/WalletLayout';
import Intl from '../../../i18n/i18n';
import Send from './Send';

function SendCoins() {
  const tabsConfig = [
    {
      label: Intl.t('sendcoins.send-tab'),
      content: <Send />,
    },
    {
      label: Intl.t('sendcoins.receive-tab'),
      content: <Button> Receive component </Button>,
    },
  ];

  const args = {
    tabsConfig,
    defaultIndex: 0,
  };

  return (
    <WalletLayout menuItem={MenuItem.SendCoins}>
      <Tabs {...args} />
    </WalletLayout>
  );
}

export default SendCoins;
