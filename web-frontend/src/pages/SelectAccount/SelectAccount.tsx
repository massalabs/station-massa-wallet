import { FiPlus, FiArrowUpRight } from 'react-icons/fi';
import { Link, useNavigate } from 'react-router-dom';
import { routeFor } from '../../utils';
import useResource from '../../custom/api/useResource';
import { AccountObject } from '../../models/AccountModel';
import { toMAS } from '@massalabs/massa-web3';
import Intl from '../../i18n/i18n';

import LandingPage from '../../layouts/LandingPage/LandingPage';
import { Button, Selector, MassaLogo } from '@massalabs/react-ui-kit';

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
          <h1 className="mas-banner text-f-primary">
            {Intl.t('account.header.title')}
          </h1>
          <label className="mas-body text-info pt-6" htmlFor="account-select">
            {Intl.t('account.select')}
          </label>
          <div id="account-select" className="pt-5">
            {accounts.map((account: AccountObject) => (
              <Link
                to={routeFor('home')}
                state={{ nickname: account.nickname }}
              >
                <div className="pb-4" key={account.nickname}>
                  <Selector
                    preIcon={<FiArrowUpRight />}
                    posIcon={<MassaLogo size={24} />}
                    content={account.nickname}
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
