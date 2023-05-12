import LandingPage from '../../layouts/LandingPage/LandingPage';
import { FiArrowRight } from 'react-icons/fi';
import { Button } from '@massalabs/react-ui-kit/src/components/Button/Button';
import { navigateToImportAccount, routeFor } from '../../utils';
import { Link, useNavigate } from 'react-router-dom';

export default function AddAccount() {
  const navigate = useNavigate();

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
