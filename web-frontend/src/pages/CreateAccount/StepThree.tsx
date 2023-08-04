import { useEffect } from 'react';

import { Button, Stepper } from '@massalabs/react-ui-kit';
import { FiArrowRight } from 'react-icons/fi';
import { Link, useLocation } from 'react-router-dom';
import { useNavigate } from 'react-router-dom';

import { usePost } from '@/custom/api';
import Intl from '@/i18n/i18n';
import LandingPage from '@/layouts/LandingPage/LandingPage';
import { AccountObject } from '@/models/AccountModel';
import { routeFor } from '@/utils';

export default function StepThree() {
  const navigate = useNavigate();
  const { state } = useLocation();

  const nickname: string = state.nickname;
  const { mutate, isSuccess, error } = usePost<AccountObject>(
    `accounts/${nickname}/backup`,
  );

  useEffect(() => {
    if (error) {
      navigate(routeFor('error'));
    } else if (isSuccess) {
      navigate(routeFor(`${nickname}/home`));
    }
  }, [isSuccess]);

  return (
    <LandingPage>
      <div className="flex flex-col justify-center items-center h-screen">
        <div className="flex flex-col justify-start w-full h-full max-w-sm max-h-56">
          <div className="w-full pb-6">
            <Stepper
              step={2}
              steps={[
                Intl.t('account.create.step1.title'),
                Intl.t('account.create.step2.title'),
                Intl.t('account.create.step3.title'),
              ]}
            />
          </div>
          <h1 className="mas-banner text-f-primary pb-6">
            {Intl.t('account.create.title')}
          </h1>
          <p className="mas-body text-f-primary pb-4">
            {Intl.t('account.create.step3.description')}
          </p>
          <div className="w-full pb-4">
            <Link to={routeFor(`${nickname}/home`)}>
              <Button posIcon={<FiArrowRight />}>
                {Intl.t('account.create.buttons.skip')}
              </Button>
            </Link>
          </div>
          <div className="pb-4 w-full">
            <Button variant="secondary" onClick={() => mutate({})}>
              {Intl.t('account.create.buttons.backup')}
            </Button>
          </div>
        </div>
      </div>
    </LandingPage>
  );
}
