import { ReactNode, useEffect } from 'react';
import { useNavigate, useParams, useOutletContext } from 'react-router-dom';
import { useResource } from '../../custom/api';
import { AccountObject } from '../../models/AccountModel';
import { routeFor } from '../../utils';
import Intl from '../../i18n/i18n';
import { IOutletContextType } from './../../pages/Base/Base';

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
  FiPlus,
} from 'react-icons/fi';

export enum MenuItem {
  Home = 'home',
  TransferCoins = 'transfer-coins',
  Transactions = 'transactions',
  Contacts = 'contacts',
  Assets = 'assets',
  Settings = 'settings',
  LightTheme = 'light-theme',
}

interface IWalletLayoutProps {
  menuItem: MenuItem;
  children: ReactNode;
}

function WalletLayout(props: IWalletLayoutProps) {
  const { menuItem } = props;
  const navigate = useNavigate();
  const { nickname } = useParams();

  const { themeIcon, themeLabel, theme, handleSetTheme } =
    useOutletContext<IOutletContextType>();

  const {
    data: accounts = [],
    error,
    isLoading,
  } = useResource<AccountObject[]>('accounts');

  const hasAccounts = !isLoading && accounts;

  useEffect(() => {
    if (error) {
      navigate('/error');
    } else if (!hasAccounts) {
      navigate('/index');
    }
  }, [accounts, error, navigate]);

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
      active: isActive(MenuItem.TransferCoins),
      footer: false,
      onClickItem: () =>
        navigate(routeFor(`${nickname}/${MenuItem.TransferCoins}`)),
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
      label: themeLabel,
      icon: themeIcon,
      active: isActive(MenuItem.LightTheme),
      footer: true,
      onClickItem: () => handleSetTheme(),
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

  const selectedAccountKey: number = parseInt(
    Object.keys(accounts).find(
      (_, idx) => accounts[idx].nickname === nickname,
    ) || '0',
  );

  return (
    // TODO
    // remove ${theme}
    // this needs to be removed as soon we fix the steps to create an account
    <div className={`${theme} bg-primary`}>
      <SideMenu conf={menuConf} items={menuItems} />
      <div className="flex justify-center items-center h-screen text-f-primary">
        {props?.children}
      </div>
      <div className="absolute top-0 right-0 p-6">
        <Dropdown options={accountsItems} select={selectedAccountKey} />
      </div>
    </div>
  );
}

export default WalletLayout;
