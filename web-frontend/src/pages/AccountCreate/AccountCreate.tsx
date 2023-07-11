import { Link, useNavigate } from 'react-router-dom';
import { routeFor } from '@/utils';
import { AccountObject } from '@/models/AccountModel';
import { usePut } from '@/custom/api';
import Intl from '@/i18n/i18n';

import { Button } from '@massalabs/react-ui-kit';
import LandingPage from '@/layouts/LandingPage/LandingPage';
import { FiArrowRight } from 'react-icons/fi';

export default function AccountCreate() {
  const navigate = useNavigate();

  const { mutate, isSuccess, data } = usePut<AccountObject>('accounts');

  if (isSuccess) {
    navigate(routeFor(`${data.nickname}/home`));
  }

  return (
    <LandingPage>
      <div className="flex flex-col justify-center items-center h-screen">
        <div className="flex flex-col justify-start items-center w-full h-full max-w-sm max-h-56">
          <h1 className="mas-banner text-neutral pb-6">
            {Intl.t('account-create.title')}
          </h1>
          <Link
            className="pb-4 w-full"
            to={routeFor('account-create-step-one')}
          >
            <Button posIcon={<FiArrowRight />}>
              {Intl.t('account-create.buttons.create')}
            </Button>
          </Link>
          <Button
            variant="secondary"
            onClick={() => mutate({} as AccountObject)}
          >
            {Intl.t('account-create.buttons.import')}
          </Button>
        </div>
      </div>
    </LandingPage>
  );
}
