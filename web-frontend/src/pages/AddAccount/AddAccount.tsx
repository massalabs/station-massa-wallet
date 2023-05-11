import LandingPage from '../../layouts/LandingPage/LandingPage';
import { FiArrowRight } from 'react-icons/fi';
import { Button } from '@massalabs/react-ui-kit/src/components/Button/Button';

export default function AddAccount() {
  const topButtonArgs = {
    onClick: () => {
      console.log('clicked');
    },
    posIcon: <FiArrowRight />,
  };
  const botButtonArgs = {
    onClick: () => {
      console.log('clicked');
    },
  };
  const defaultFlex = 'flex flex-col justify-center items-center align-center';
  return (
    <LandingPage>
      <div className={`${defaultFlex} h-screen`}>
        <div className="w-fit h-fit">
          <div className="mb-6">
            <h1 className="mas-banner text-neutral">Add an account</h1>
          </div>
          <div className="mb-4">
            <Button {...topButtonArgs}>Create an account</Button>
          </div>
          <div className="mb-4">
            <Button variant="secondary" {...botButtonArgs}>
              Import an existing account
            </Button>
          </div>
        </div>
      </div>
    </LandingPage>
  );
}
