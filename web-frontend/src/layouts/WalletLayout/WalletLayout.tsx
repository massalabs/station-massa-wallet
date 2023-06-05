import { ReactNode, useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useResource } from '../../custom/api';
import { AccountObject } from '../../models/AccountModel';
import { routeFor } from '../../utils';
import Intl from '../../i18n/i18n';

import {
  Dropdown,
  MassaWallet,
  SideMenu,
  Identicon,
} from '@massalabs/react-ui-kit';
import {
  FiHome,
  FiList,
  FiArrowUpRight,
  FiUsers,
  FiDisc,
  FiSettings,
  FiSun,
  FiPlus,
} from 'react-icons/fi';

export enum MenuItem {
  Home = 'home',
  SendCoins = 'send-coins',
  Transactions = 'transactions',
  Contacts = 'contacts',
  Assets = 'assets',
  Settings = 'settings',
  LightTheme = 'light-theme',
}

interface WalletProps {
  menuItem: MenuItem;
  children: ReactNode;
}

function WalletLayout(props: WalletProps) {
  const { menuItem } = props;
  const navigate = useNavigate();
  const { nickname } = useParams();

  const { data: accounts = [], refetch } = useResource<AccountObject[]>(
    'accounts',
    {
      refetchOnWindowFocus: false,
      enabled: false, // disable this query from automatically running
      manual: true,
    },
  );

  refetch();

  const [selectedAccountKey, setSelectedAccountKey] = useState(0);

  useEffect(() => {
    const index = Object.keys(accounts).find((_, idx) => {
      return accounts[idx].nickname === nickname;
    });

    if (index === undefined) {
      setSelectedAccountKey(0);
    } else {
      setSelectedAccountKey(parseInt(index));
    }
  }, [accounts, refetch]);

  // If no account, redirect to welcome page
  if (!accounts.length) {
    navigate(routeFor('index'));
    return null;
  }

  function isActive(item: MenuItem) {
    return item === menuItem;
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
      onClickItem: () => navigate(routeFor(`${nickname}/${MenuItem.Home}`)),
    },
    {
      label: 'Transactions',
      icon: <FiList />,
      active: isActive(MenuItem.Transactions),
      footer: false,
      onClickItem: () =>
        navigate(routeFor(`${nickname}/${MenuItem.Transactions}`)),
    },
    {
      label: 'Send/Receive',
      icon: <FiArrowUpRight />,
      active: isActive(MenuItem.SendCoins),
      footer: false,
      onClickItem: () =>
        navigate(routeFor(`${nickname}/${MenuItem.SendCoins}`)),
    },
    {
      label: 'Contacts',
      icon: <FiUsers />,
      active: isActive(MenuItem.Contacts),
      footer: false,
      onClickItem: () => navigate(routeFor(`${nickname}/${MenuItem.Contacts}`)),
    },
    {
      label: 'Assets',
      icon: <FiDisc />,
      active: isActive(MenuItem.Assets),
      footer: false,
      onClickItem: () => navigate(routeFor(`${nickname}/${MenuItem.Assets}`)),
    },
    {
      label: 'Settings',
      icon: <FiSettings />,
      active: isActive(MenuItem.Settings),
      footer: true,
      onClickItem: () => navigate(routeFor(`${nickname}/${MenuItem.Settings}`)),
    },
    {
      label: 'Light theme',
      icon: <FiSun />,
      active: isActive(MenuItem.LightTheme),
      footer: true,
    },
  ];

  const accountsItems = accounts.map((account) => ({
    icon: <Identicon username={account.nickname} size={32} />,
    item: account.nickname,
    onClick: () => navigate(routeFor(`${account.nickname}/${MenuItem.Home}`)),
  }));

  accountsItems.push({
    icon: <FiPlus size={32} />,
    item: Intl.t('account.add'),
    onClick: () => navigate(routeFor('account-create')),
  });

  return (
    <div className="bg-primary min-h-screen">
      <div
        className="flex flex-col justify-center h-screen
        w-1/2 min-w-fit text-f-primary m-auto"
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
