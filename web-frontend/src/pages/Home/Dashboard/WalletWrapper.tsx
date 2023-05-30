import { DashboardLayout } from '../../../layouts/DashboardLayout';
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
import useResource from '../../../custom/api/useResource';
import { AccountObject } from '../../../models/AccountModel';
import { routeFor } from '../../../utils';
import { useState } from 'react';
import { TransferCoin } from './transferCoin';

enum MenuItem {
  Home = 'home',
  SendReceive = 'send-receive',
  Transactions = 'transactions',
  Contacts = 'contacts',
  Assets = 'assets',
  Settings = 'settings',
  LightTheme = 'light-theme',
}

export default function DashBoard() {
  const navigate = useNavigate();
  const { state } = useLocation();
  var nickname: string = state.nickname;

  const { data: accounts = [] } = useResource<AccountObject[]>('accounts');

  // If no account, redirect to welcome page
  if (!accounts.length) {
    navigate(routeFor('index'));
  }

  const [selectedAccount, setSelectedAccount] = useState<AccountObject>(
    accounts.find((account) => account.nickname === nickname)!,
  );

  const [selectedItem, setSelectedItem] = useState<MenuItem>(
    state.menuItem ?? MenuItem.SendReceive,
  );

  const isActive = (item: MenuItem) => {
    return item === selectedItem;
  };

  let menuConf = {
    title: 'MassaWallet',
    logo: <MassaWallet />,
    fullMode: true,
  };

  let menuItems = [
    {
      label: 'Home',
      icon: <FiHome />,
      active: false,
      footer: false,
      onClickItem: () => setSelectedItem(MenuItem.Home),
    },
    {
      label: 'Transactions',
      icon: <FiList />,
      active: isActive(MenuItem.Transactions),
      footer: false,
      onClickItem: () => setSelectedItem(MenuItem.Transactions),
    },
    {
      label: 'Send/Receive',
      icon: <FiArrowUpRight />,
      active: isActive(MenuItem.SendReceive),
      footer: false,
      onClickItem: () => setSelectedItem(MenuItem.SendReceive),
    },
    {
      label: 'Contacts',
      icon: <FiUsers />,
      active: isActive(MenuItem.Contacts),
      footer: false,
      onClickItem: () => setSelectedItem(MenuItem.Contacts),
    },
    {
      label: 'Assets',
      icon: <FiDisc />,
      active: isActive(MenuItem.Assets),
      footer: false,
      onClickItem: () => setSelectedItem(MenuItem.Assets),
    },
    {
      label: 'Settings',
      icon: <FiSettings />,
      active: isActive(MenuItem.Settings),
      footer: true,
      onClickItem: () => setSelectedItem(MenuItem.Settings),
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
    onClick: () => setSelectedAccount(account),
  }));

  const selectedAccountKey: number = parseInt(
    Object.keys(accounts).find(
      (_, idx) => accounts[idx].nickname === nickname,
    )!,
  );

  let content;
  switch (selectedItem) {
    case MenuItem.SendReceive:
      content = <TransferCoin account={selectedAccount} />;
      break;
    default:
      content = null;
  }

  return (
    <DashboardLayout>
      <SideMenu conf={menuConf} items={menuItems} />
      <div className="absolute top-0 right-0 m-6">
        <Dropdown options={accountsItems} select={selectedAccountKey} />
      </div>
      {content}
    </DashboardLayout>
  );
}
