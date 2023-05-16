import LandingPage from '../../layouts/LandingPage/LandingPage';
import { Button } from '@massalabs/react-ui-kit/src/components/Button/Button';
import { FiArrowRight } from 'react-icons/fi';
import { Link, useNavigate } from 'react-router-dom';
import { goToErrorPage, routeFor } from '../../utils';
import useResource from '../../custom/api/useResource';
import { AccountObject } from '../../models/AccountModel';
import usePut from '../../custom/api/usePut';
import Intl from '../../i18n/i18n';

export default function Welcome() {
  const navigate = useNavigate();

  const { error, data = [] } = useResource<AccountObject[]>('accounts');

  if (error) goToErrorPage(navigate);

  if (data.length) {
    navigate(routeFor('account-select'));
  }

  const handleImport = usePut<AccountObject>('accounts');

  return (
    <LandingPage>
      <div className="flex flex-col justify-center items-center h-screen">
        <div className="w-fit h-fit max-w-lg">
          <h1 className="mas-banner text-f-primary">
            {Intl.t('welcome.title_first_part')}
            <br />
            {Intl.t('welcome.title_second_part')}
            <span className="text-brand">
              {Intl.t('welcome.title_third_part')}
            </span>
          </h1>
          <div className="pt-6">
            <Link to={routeFor('account-create')}>
              <Button posIcon={<FiArrowRight />}>
                {Intl.t('account.create')}
              </Button>
            </Link>
          </div>
          <div className="pt-3.5">
            <Button
              variant="secondary"
              onClick={() => {
                handleImport.mutate({} as AccountObject);
              }}
            >
              {Intl.t('account.import')}
            </Button>
          </div>
        </div>
      </div>
    </LandingPage>
  );
}
