import LandingPage from '../../layouts/LandingPage/LandingPage';
import { FiUser, FiPlus } from 'react-icons/fi';
import { Button } from '@massalabs/react-ui-kit/src/components/Button/Button';
import { AccountSelector } from '@massalabs/react-ui-kit/src/components/AccountSelector/AccountSelector';
import { MassaToken } from '@massalabs/react-ui-kit/src/components/Icons/Svg/SvgComponent/MassaToken';
import { toMAS } from '@massalabs/massa-web3';
import { Link, useNavigate } from 'react-router-dom';
import { goToErrorPage, routeFor } from '../../utils';
import useResource from '../../custom/api/useResource';
import { AccountObject } from '../../models/AccountModel';

export default function SelectAccount() {
  // If no account, redirect to welcome page
  const navigate = useNavigate();

  const { error, data = [] } = useResource<AccountObject[]>('accounts');

  if (error) goToErrorPage(navigate);

  if (!data.length) {
    navigate(routeFor('welcome'));
  }

  const defaultFlex = 'flex flex-col justify-center items-center align-center';
  return (
    <LandingPage>
      <div className={`${defaultFlex} h-screen`}>
        <div className="w-1/2">
          <div className="mb-6">
            <p className="mas-banner text-neutral"> Hey!</p>
          </div>
          <div className="mb-6">
            <label className="mas-body text-neutral" htmlFor="account-select">
              Select an account
            </label>
          </div>
          <div id="account-select" className="w-full flex flex-col">
            {data.map((account: AccountObject) => (
              <div className="mb-4" key={account.nickname}>
                <AccountSelector
                  avatar={<FiUser className="text-neutral h-6 w-6" />}
                  icon={<MassaToken size={24} />}
                  accountName={account.nickname}
                  amount={toMAS(account.candidateBalance).toString()}
                />
              </div>
            ))}
            <div>
              <Link to={routeFor('account-create')}>
                <Button variant="secondary" preIcon={<FiPlus />}>
                  Add an account
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </div>
    </LandingPage>
  );
}
