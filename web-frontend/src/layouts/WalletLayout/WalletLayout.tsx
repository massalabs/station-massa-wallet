import { ReactNode, useEffect } from 'react';

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
import {
  useNavigate,
  useParams,
  useLocation,
  useOutletContext,
  createSearchParams,
} from 'react-router-dom';

import Intl from '../../i18n/i18n';
import { fetchAccounts, routeFor } from '../../utils';

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

interface IOutletContextType {
  themeIcon: JSX.Element;
  themeLabel: string;
  theme: string;
  handleSetTheme: () => void;
}

export function WalletLayout(props: IWalletLayoutProps) {
  const { menuItem } = props;
  const navigate = useNavigate();
  const location = useLocation();
  const { nickname } = useParams();

  const { themeIcon, themeLabel, theme, handleSetTheme } =
    useOutletContext<IOutletContextType>();

  const { okAccounts: accounts = [], isLoading, error } = fetchAccounts();

  const hasAccounts = !isLoading && accounts;

  useEffect(() => {
    if (error) {
      navigate(routeFor('error'));
    } else if (!hasAccounts) {
      navigate(routeFor('index'));
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
      label: Intl.t('menu.home'),
      icon: <FiHome data-testid="side-menu-wallet-icon" />,
      active: isActive(MenuItem.Home),
      footer: false,
      onClickItem: () => navigate(routeFor(`${nickname}/${MenuItem.Home}`)),
    },
    {
      label: Intl.t('menu.transactions'),
      icon: <FiList data-testid="side-menu-transactions-icon" />,
      active: isActive(MenuItem.Transactions),
      footer: false,
      onClickItem: () =>
        navigate(routeFor(`${nickname}/${MenuItem.Transactions}`)),
    },
    {
      label: Intl.t('menu.send-receive'),
      icon: <FiArrowUpRight data-testid="side-menu-sendreceive-icon" />,
      active: isActive(MenuItem.TransferCoins),
      footer: false,
      onClickItem: () =>
        navigate(routeFor(`${nickname}/${MenuItem.TransferCoins}`)),
    },
    {
      label: Intl.t('menu.contacts'),
      icon: <FiUsers data-testid="side-menu-contacts-icon" />,
      active: isActive(MenuItem.Contacts),
      footer: false,
      onClickItem: () => navigate(routeFor(`${nickname}/${MenuItem.Contacts}`)),
    },
    {
      label: Intl.t('menu.assets'),
      icon: <FiDisc data-testid="side-menu-assets-icon" />,
      active: isActive(MenuItem.Assets),
      footer: false,
      onClickItem: () => navigate(routeFor(`${nickname}/${MenuItem.Assets}`)),
    },
    {
      label: Intl.t('menu.settings'),
      icon: <FiSettings data-testid="side-menu-settings-icon" />,
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
    onClick: () => {
      const lastUrl = location.pathname.split('/').pop();

      return navigate(routeFor(`${account.nickname}/${lastUrl}`));
    },
  }));

  accountsItems.push({
    icon: <FiPlus size={32} />,
    item: Intl.t('account.add'),
    onClick: () =>
      navigate(
        nickname
          ? {
              pathname: routeFor('account-create'),
              search: createSearchParams({ from: nickname }).toString(),
            }
          : routeFor('account-create'),
      ),
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
        <div className="w-64">
          <Dropdown options={accountsItems} select={selectedAccountKey} />
        </div>
      </div>
    </div>
  );
}
