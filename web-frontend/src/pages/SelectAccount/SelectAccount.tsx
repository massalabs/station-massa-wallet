import { FiPlus } from 'react-icons/fi';
import { Link, useNavigate } from 'react-router-dom';
import { routeFor } from '../../utils';
import { useResource } from '../../custom/api';
import { AccountObject } from '../../models/AccountModel';
import { toMAS } from '@massalabs/massa-web3';
import Intl from '../../i18n/i18n';

import LandingPage from '../../layouts/LandingPage/LandingPage';
import {
  Button,
  Selector,
  MassaLogo,
  Identicon,
} from '@massalabs/react-ui-kit';
import { formatStandard } from '../../utils/MassaFormating';

export default function SelectAccount() {
  const navigate = useNavigate();

  const { data: accounts = [] } = useResource<AccountObject[]>('accounts');

  // If no account, redirect to welcome page
  if (!accounts.length) {
    navigate(routeFor('index'));
  }

  function getFormattedBalance(account: AccountObject): string {
    return formatStandard(toMAS(account.candidateBalance).toNumber());
  }

  return (
    <LandingPage>
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
    </LandingPage>
  );
}
