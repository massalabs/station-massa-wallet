import LandingPage from '../../layouts/LandingPage/LandingPage';
import { Button } from '@massalabs/react-ui-kit/src/components/Button/Button';
import { FiArrowRight } from 'react-icons/fi';
import { useQuery } from '@tanstack/react-query';
import { getAllAccounts } from '../../api/account';
import { useNavigate } from 'react-router-dom';
import { navigateToImportAccount, routeFor } from '../../utils';

export default function Welcome() {
  const navigate = useNavigate();

  const accounts = useQuery({
    queryKey: ['accounts'],
    queryFn: getAllAccounts,
  });

  if (accounts.data?.length > 0) {
    navigate(routeFor('account-select'));
  }

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
              onClick={navigateToImportAccount(navigate)}
            >
              Import an existing account
            </Button>
          </div>
        </div>
      </div>
    </LandingPage>
  );
}
