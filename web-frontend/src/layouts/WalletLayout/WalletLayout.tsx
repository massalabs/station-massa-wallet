import { ReactNode } from 'react';
import { Dropdown, MassaWallet, SideMenu } from '@massalabs/react-ui-kit';
import {
  FiHome,
  FiList,
  FiArrowUpRight,
  FiUser,
  FiUsers,
  FiDisc,
  FiSettings,
  FiSun,
} from 'react-icons/fi';
import { useLocation, useNavigate } from 'react-router-dom';
import useResource from '../../custom/api/useResource';
import { AccountObject } from '../../models/AccountModel';
import { routeFor } from '../../utils';
import { useState } from 'react';

enum MenuItem {
  Home = 'home',
  SendReceive = 'send-receive',
  Transactions = 'transactions',
  Contacts = 'contacts',
  Assets = 'assets',
  Settings = 'settings',
  LightTheme = 'light-theme',
}

interface WalletProps {
  children: ReactNode;
}

function WalletLayout(props: WalletProps) {
  const navigate = useNavigate();
  const { state } = useLocation();
  var nickname: string = state.nickname;

  const { data: accounts = [] } = useResource<AccountObject[]>('accounts');

  // If no account, redirect to welcome page
  if (!accounts.length) {
    navigate(routeFor('index'));
  }

  const [selectedItem, setSelectedItem] = useState<MenuItem>(
    state.menuItem ?? MenuItem.Home,
  );

  function isActive(item: MenuItem) {
    return item === selectedItem;
  }

  let menuConf = {
    title: 'MassaWallet',
    logo: <MassaWallet />,
    fullMode: true,
  };

  let menuItems = [
    {
      label: 'Wallet',
      icon: <FiHome />,
      active: isActive(MenuItem.Home),
      footer: false,
      onClickItem: () => navigate(routeFor(MenuItem.Home), { state }),
    },
    {
      label: 'Transactions',
      icon: <FiList />,
      active: isActive(MenuItem.Transactions),
      footer: false,
      onClickItem: () => navigate(routeFor(MenuItem.Transactions), { state }),
    },
    {
      label: 'Send/Receive',
      icon: <FiArrowUpRight />,
      active: isActive(MenuItem.SendReceive),
      footer: false,
      onClickItem: () => navigate(routeFor(MenuItem.SendReceive), { state }),
    },
    {
      label: 'Contacts',
      icon: <FiUsers />,
      active: isActive(MenuItem.Contacts),
      footer: false,
      onClickItem: () => navigate(routeFor(MenuItem.Contacts), { state }),
    },
    {
      label: 'Assets',
      icon: <FiDisc />,
      active: isActive(MenuItem.Assets),
      footer: false,
      onClickItem: () => navigate(routeFor(MenuItem.Assets), { state }),
    },
    {
      label: 'Settings',
      icon: <FiSettings />,
      active: isActive(MenuItem.Settings),
      footer: true,
      onClickItem: () => navigate(routeFor(MenuItem.Settings), { state }),
    },
    {
      label: 'Light theme',
      icon: <FiSun />,
      active: isActive(MenuItem.LightTheme),
      footer: true,
      onClickItem: () => setSelectedItem(MenuItem.LightTheme),
    },
  ];

  const accountsItems = accounts.map((account) => ({
    icon: <FiUser size={32} />,
    item: account.nickname,
  }));

  const selectedAccountKey: number = parseInt(
    Object.keys(accounts).find(
      (_, idx) => accounts[idx].nickname === nickname,
    )!,
  );

  return (
    <div className="bg-primary min-h-screen">
      <div
        className="flex flex-col justify-center h-screen
        max-w-xs min-w-fit text-f-primary m-auto"
      >
        <SideMenu conf={menuConf} items={menuItems} />
        <div className="absolute top-0 right-0 m-6">
          <Dropdown options={accountsItems} select={selectedAccountKey} />
        </div>
        {props?.children}
      </div>
    </div>
  );
}

export default WalletLayout;
