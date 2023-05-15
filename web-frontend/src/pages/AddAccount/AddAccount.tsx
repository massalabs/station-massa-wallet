import LandingPage from '../../layouts/LandingPage/LandingPage';
import { FiArrowRight } from 'react-icons/fi';
import { Button } from '@massalabs/react-ui-kit/src/components/Button/Button';
import { routeFor } from '../../utils';
import { Link } from 'react-router-dom';
import { AccountObject } from '../../models/AccountModel';
import usePut from '../../custom/api/usePut';

export default function AddAccount() {
  const handleImport = usePut<AccountObject>('accounts');

  const defaultFlex = 'flex flex-col justify-center items-center align-center';
  return (
    <LandingPage>
      <div className={`${defaultFlex} h-screen`}>
        <div className="w-fit h-fit">
          <div className="mb-6">
            <h1 className="mas-banner text-neutral">Add an account</h1>
          </div>
          <div className="mb-4">
            <Link to={routeFor('')}>
              <Button posIcon={<FiArrowRight />}>Create an account</Button>
            </Link>
          </div>
          <div className="mb-4">
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
