import { useLocation, useNavigate } from 'react-router-dom';
import usePost from '../../custom/api/usePost';
import { AccountObject } from '../../models/AccountModel';
import { routeFor } from '../../utils';
import Intl from '../../i18n/i18n';

import { Button, Stepper } from '@massalabs/react-ui-kit';
import LandingPage from '../../layouts/LandingPage/LandingPage';
import { FiLock } from 'react-icons/fi';

export default function StepTwo() {
  const navigate = useNavigate();

  const { state } = useLocation();
  const nickname = state.nickname;

  const { mutate, isSuccess } = usePost<AccountObject>(`accounts/${nickname}`);

  if (isSuccess) {
    navigate(routeFor('account-create-step-three'), { state: { nickname } });
  }

  return (
    <LandingPage>
      <div className="flex flex-col justify-center items-center align-center h-screen">
        <div className="flex flex-col justify-center items-center w-fit h-fit max-w-sm">
          <div className="w-full max-w-xs mb-6">
            <Stepper
              step={1}
              steps={[
                Intl.t('account.create.step1.title'),
                Intl.t('account.create.step2.title'),
                Intl.t('account.create.step3.title'),
              ]}
            />
          </div>
          <div className="w-full">
            <h1 className="mas-banner text-neutral mb-6">
              {Intl.t('account.create.title')}
            </h1>
          </div>
          <div className="w-full mb-4">
            <p className="mas-body text-neutral">
              {Intl.t('account.create.step2.description')}
            </p>
          </div>
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
