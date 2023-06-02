import { Link, useNavigate } from 'react-router-dom';
import { goToErrorPage, routeFor } from '../../utils';
import { useResource, usePut } from '../../custom/api';
import { AccountObject } from '../../models/AccountModel';
import Intl from '../../i18n/i18n';

import { Button } from '@massalabs/react-ui-kit/src/components/Button/Button';
import LandingPage from '../../layouts/LandingPage/LandingPage';
import { FiArrowRight } from 'react-icons/fi';

export default function Welcome() {
  const navigate = useNavigate();

  const { error, data = [] } = useResource<AccountObject[]>('accounts');

  if (error) goToErrorPage(navigate);

  const { mutate, isSuccess } = usePut<AccountObject>('accounts');

  if (data.length || isSuccess) {
    navigate(routeFor('account-select'));
  }

  return (
    <LandingPage>
      <div className="flex flex-col justify-center items-center h-screen">
        <div className="flex flex-col justify-start items-start w-full h-full max-w-sm max-h-56">
          <h1 className="mas-banner text-f-primary pb-6">
            {Intl.t('welcome.title-first-part')}
            <br />
            {Intl.t('welcome.title-second-part')}
            <span className="text-brand">
              {Intl.t('welcome.title-third-part')}
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
      </div>
    </LandingPage>
  );
}
