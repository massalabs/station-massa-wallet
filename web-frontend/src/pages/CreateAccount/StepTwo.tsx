import { useEffect } from 'react';

import { Button, Stepper } from '@massalabs/react-ui-kit';
import { FiLock } from 'react-icons/fi';
import { useLocation, useNavigate } from 'react-router-dom';

import { usePost } from '@/custom/api/usePost';
import Intl from '@/i18n/i18n';
import LandingPage from '@/layouts/LandingPage/LandingPage';
import { AccountObject } from '@/models/AccountModel';
import { routeFor } from '@/utils';

export default function StepTwo() {
  const navigate = useNavigate();
  const { state } = useLocation();
  const nickname = state.nickname;

  const { mutate, isSuccess, error } = usePost<AccountObject>(
    `accounts/${nickname}`,
  );

  useEffect(() => {
    if (error) {
      navigate(routeFor('error'));
    } else if (isSuccess) {
      navigate(routeFor('account-create-step-three'), { state: { nickname } });
    }
  }, [isSuccess]);

  return (
    <LandingPage>
      <div className="flex flex-col justify-center items-center h-screen">
        <div className="flex flex-col justify-start w-full h-full max-w-sm max-h-56">
          <div className="w-full pb-6">
            <Stepper
              step={1}
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
          <p className="mas-body text-f-primary w-full pb-4">
            {Intl.t('account.create.step2.description')}
          </p>
          <div className="w-full">
            <Button
              preIcon={<FiLock />}
              onClick={() => mutate({} as AccountObject)}
            >
              {Intl.t('account.create.buttons.define_password')}
            </Button>
          </div>
        </div>
      </div>
    </LandingPage>
  );
}
