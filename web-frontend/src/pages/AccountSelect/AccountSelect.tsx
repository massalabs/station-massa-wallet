import { useEffect } from 'react';

import { toMAS } from '@massalabs/massa-web3';
import {
  Button,
  Selector,
  MassaLogo,
  Identicon,
} from '@massalabs/react-ui-kit';
import { FiPlus } from 'react-icons/fi';
import { Link, useNavigate } from 'react-router-dom';

import { Loading } from './Loading';
import { useResource } from '@/custom/api';
import Intl from '@/i18n/i18n';
import LandingPage from '@/layouts/LandingPage/LandingPage';
import { AccountObject } from '@/models/AccountModel';
import { routeFor, formatStandard } from '@/utils';

export default function AccountSelect() {
  const navigate = useNavigate();

  const { okAccounts: accounts, isLoading, error } = fetchAccounts();
  const hasAccounts = accounts?.length;

  useEffect(() => {
    if (error) {
      navigate(routeFor('error'));
    } else if (!isLoading && !hasAccounts) {
      navigate(routeFor('index'));
    }
  }, [accounts, navigate, hasAccounts, error, isLoading]);

  function getFormattedBalance(account: AccountObject): string {
    return formatStandard(toMAS(account.candidateBalance).toNumber());
  }

  return (
    <LandingPage>
      {isLoading ? (
        <Loading />
      ) : hasAccounts ? (
        <div className="flex flex-col justify-center items-center h-screen">
          <div className="flex flex-col justify-center items-start w-full h-full max-w-lg">
            <h1 className="mas-banner text-f-primary pb-6">
              {Intl.t('account.header.title')}
            </h1>
            <label className="mas-body text-info pb-6" htmlFor="account-select">
              {Intl.t('account.select')}
            </label>
            <div id="account-select" className="pb-4 w-full">
              {accounts.map((account: AccountObject, index: number) => (
                <Link
                  key={index}
                  className="w-full"
                  to={routeFor(`${account.nickname}/home`)}
                >
                  <div className="pb-4" key={account.nickname}>
                    <Selector
                      preIcon={
                        <Identicon username={account.nickname} size={32} />
                      }
                      posIcon={<MassaLogo size={24} />}
                      content={account.nickname}
                      amount={getFormattedBalance(account)}
                    />
                  </div>
                </Link>
              ))}
              <Link to={routeFor('account-create')}>
                <Button variant="secondary" preIcon={<FiPlus />}>
                  {Intl.t('account.add')}
                </Button>
              </Link>
            </div>
          </div>
        </div>
      ) : null}
    </LandingPage>
  );
}
