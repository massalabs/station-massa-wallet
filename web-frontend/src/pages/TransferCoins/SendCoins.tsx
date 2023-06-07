import { Tabs } from '@massalabs/react-ui-kit';
import WalletLayout, {
  MenuItem,
} from '../../layouts/WalletLayout/WalletLayout';
import Intl from '../../i18n/i18n';
import Send from './SendCoins/Send';
import ReceiveCoins from './ReceiveCoins/ReceiveCoins';

function SendCoins() {
  // this is place holder code for texting purpose
  // The nickname will be passed as an argument to the CopyAddress.tsx
  const address = 'AU12irbDfYNwyZRbnpBrfCBPCxrktp8f8riK2sQddWbzQ3g43G7bb';
  const formattedAddress = address.slice(0, 4) + '...' + address.slice(-4);

  const componentArgs = {
    formattedAddress,
    address,
  };

  const tabsConfig = [
    {
      label: Intl.t('sendcoins.send-tab'),
      content: <Send />,
    },
    {
      label: Intl.t('sendcoins.receive-tab'),
      content: (
        <div className="mt-20">
          <ReceiveCoins {...componentArgs} />
        </div>
      ),
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
