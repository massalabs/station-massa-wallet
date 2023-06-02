import { Link, useLocation } from 'react-router-dom';
import Intl from '../../i18n/i18n';
import { routeFor } from '../../utils';

import { Button, Stepper } from '@massalabs/react-ui-kit';
import { FiArrowRight } from 'react-icons/fi';
import LandingPage from '../../layouts/LandingPage/LandingPage';
import { useNavigate } from 'react-router-dom';
import { AccountObject } from '../../models/AccountModel';
import usePost from '../../custom/api/usePost';

export default function StepThree() {
  const navigate = useNavigate();

  const { state } = useLocation();
  const nickname: string = state.nickname;

  const { mutate, isSuccess } = usePost<AccountObject>(
    `accounts/${nickname}/backup`,
  );

  if (isSuccess) {
    navigate(routeFor('home'), { state: { nickname } });
  }

  return (
    <LandingPage>
      <div className="flex flex-col justify-center items-center align-center h-screen">
        <div className="flex flex-col justify-center items-center w-fit h-fit max-w-sm text-neutral">
          <div className="w-full max-w-xs">
            <Stepper
              step={2}
              steps={[
                Intl.t('account.create.step1.title'),
                Intl.t('account.create.step2.title'),
                Intl.t('account.create.step3.title'),
              ]}
            />
          </div>
          <h1 className="mas-banner pt-6">{Intl.t('account.create.title')}</h1>
          <p className="mas-body pt-6">
            {Intl.t('account.create.step3.description')}
          </p>
          <div className="pt-4 w-full">
            <Link to={routeFor('home')} state={{ nickname }}>
              <Button posIcon={<FiArrowRight />}>
                {Intl.t('account.create.buttons.skip')}
              </Button>
            </Link>
          </div>
          <div className="pt-4 w-full">
            <Button
              variant="secondary"
              onClick={() => mutate({} as AccountObject)}
            >
              {Intl.t('account.create.buttons.backup')}
            </Button>
          </div>
        </div>
      </div>
    </LandingPage>
  );
}
