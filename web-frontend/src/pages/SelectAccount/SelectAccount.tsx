import { FiUser, FiPlus } from 'react-icons/fi';
import { Link, useNavigate } from 'react-router-dom';
import { routeFor } from '../../utils';
import useResource from '../../custom/api/useResource';
import { AccountObject } from '../../models/AccountModel';
import { toMAS } from '@massalabs/massa-web3';
import Intl from '../../i18n/i18n';

import LandingPage from '../../layouts/LandingPage/LandingPage';
import { Button, AccountSelector, MassaToken } from '@massalabs/react-ui-kit';

export default function SelectAccount() {
  const navigate = useNavigate();

  const { data: accounts = [] } = useResource<AccountObject[]>('accounts');

  // If no account, redirect to welcome page
  if (!accounts.length) {
    navigate(routeFor('index'));
  }

  const defaultFlex = 'flex flex-col justify-center items-center align-center';

  return (
    <LandingPage>
      <div className={`${defaultFlex} h-screen`}>
        <div className="w-1/2">
          <div className="mb-6">
            <p className="mas-banner text-neutral">
              {Intl.t('account.header.title')}
            </p>
          </div>
          <div className="mb-6">
            <label className="mas-body text-neutral" htmlFor="account-select">
              {Intl.t('account.select')}
            </label>
          </div>
          <div id="account-select" className="w-full flex flex-col">
            {accounts.map((account: AccountObject) => (
              <Link
                to={routeFor('dashboard')}
                state={{ nickname: account.nickname }}
              >
                <div className="mb-4" key={account.nickname}>
                  <AccountSelector
                    avatar={<FiUser className="text-neutral h-6 w-6" />}
                    icon={<MassaToken size={24} />}
                    accountName={account.nickname}
                    amount={toMAS(account.candidateBalance).toString()}
                  />
                </div>
              </Link>
            ))}
            <div>
              <Link to={routeFor('account-create')}>
                <Button variant="secondary" preIcon={<FiPlus />}>
                  {Intl.t('account.add')}
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </div>
    </LandingPage>
  );
}
