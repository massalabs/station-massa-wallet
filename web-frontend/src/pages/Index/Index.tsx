import { useEffect } from 'react';

import { Button } from '@massalabs/react-ui-kit/src/components/Button/Button';
import { FiAlertTriangle, FiArrowRight } from 'react-icons/fi';
import { Link, useNavigate } from 'react-router-dom';

import { Loading } from './Loading';
import { usePut } from '@/custom/api';
import Intl from '@/i18n/i18n';
import LandingPage from '@/layouts/LandingPage/LandingPage';
import { AccountObject } from '@/models/AccountModel';
import { fetchAccounts, routeFor } from '@/utils';

export default function Index() {
  const { error, okAccounts, corruptedAccounts, isLoading } = fetchAccounts();
  const corruptedAccountsNames = corruptedAccounts?.map(
    (account: AccountObject) => account.nickname,
  );
  const accountsStringified =
    corruptedAccountsNames
      ?.map(
        (name: string, index: number) =>
          `${name}${index === corruptedAccountsNames.length - 1 ? '' : ', '}`,
      )
      .join('') || '';
  const corruptedAccountsCount = corruptedAccountsNames
    ? corruptedAccountsNames.length
    : 0;
  const { mutate, isSuccess } = usePut<AccountObject>('accounts');
  const navigate = useNavigate();
  const hasAccounts = okAccounts?.length;
  const warningMessage = corruptedAccountsCount
    ? corruptedAccountsCount === 1
      ? Intl.t('account.corrupted-account', { account: accountsStringified })
      : Intl.t('account.corrupted-accounts', { accounts: accountsStringified })
    : undefined;

  useEffect(() => {
    if (error) {
      navigate(routeFor('error'));
    } else if (hasAccounts || isSuccess) {
      navigate(routeFor('account-select'));
    }
  }, [okAccounts, navigate, hasAccounts, error, isSuccess]);

  return (
    <LandingPage>
      {isLoading ? (
        <Loading />
      ) : !hasAccounts ? (
        <div className="flex flex-col justify-center items-center h-screen">
          <div className="flex flex-col justify-start items-start w-full h-full max-w-sm max-h-56">
            <h1 className="mas-banner text-f-primary pb-6">
              {Intl.t('index.title-first-part')}
              <br />
              {Intl.t('index.title-second-part')}
              <span className="text-brand">
                {Intl.t('index.title-third-part')}
              </span>
            </h1>
            <Link
              className="pb-4 w-full"
              to={routeFor('account-create-step-one')}
            >
              <Button posIcon={<FiArrowRight />}>
                {Intl.t('account.create.title')}
              </Button>
            </Link>
            <Button
              variant="secondary"
              onClick={() => {
                mutate({} as AccountObject);
              }}
            >
              {Intl.t('account.import')}
            </Button>
          </div>
          {corruptedAccountsCount > 0 ? (
            <div className="flex items-center text-f-primary w-[384px] h-fit py-6 gap-2 justify-center">
              <div className="min-w-fit flex items-center justify-center h-full text-s-warning">
                <FiAlertTriangle size={36} />
              </div>
              <p className="flex items-center justify-center h-full">
                {warningMessage}
              </p>
            </div>
          ) : null}
        </div>
      ) : null}
    </LandingPage>
  );
}
