import { Tabs, Button } from '@massalabs/react-ui-kit';
import WalletLayout from '../../../layouts/WalletLayout/WalletLayout';
import { useLocation } from 'react-router-dom';
import Send from './Send';

function SendCoins() {
  const { state } = useLocation();
  const nickname = state.nickname;

  console.log(nickname);

  const tabsConfig = [
    {
      label: 'Send',
      content: <Send />,
      onClickTab: () => console.log('Hello'),
    },
    {
      label: 'Receive',
      content: <Button> Tab 2 component </Button>,
    },
  ];

  const args = {
    tabsConfig,
    defaultIndex: 0,
  };

  return (
    <WalletLayout>
      <Tabs {...args} />
    </WalletLayout>
  );
}

export default SendCoins;
