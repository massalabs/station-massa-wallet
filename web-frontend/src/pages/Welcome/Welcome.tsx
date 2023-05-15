import LandingPage from '../../layouts/LandingPage/LandingPage';
import { Button } from '@massalabs/react-ui-kit/src/components/Button/Button';
import { FiArrowRight } from 'react-icons/fi';
import { useNavigate } from 'react-router-dom';
import { goToErrorPage, routeFor } from '../../utils';
import useResource from '../../custom/api/useResource';
import { AccountObject } from '../../models/AccountModel';
import usePut from '../../custom/api/usePut';

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
            Welcome on
            <br />
            Massa<span className="text-brand">wallet</span>
          </h1>
          <div className="pt-6">
            <Button posIcon={<FiArrowRight />}>Create an account</Button>
          </div>
          <div className="pt-3.5">
            <Button
              variant="secondary"
              onClick={() => {
                handleImport.mutate({} as AccountObject);
              }}
            >
              Import an existing account
            </Button>
          </div>
        </div>
      </div>
    </LandingPage>
  );
}
