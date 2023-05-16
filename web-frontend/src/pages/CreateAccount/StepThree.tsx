import { Stepper } from '@massalabs/react-ui-kit/src/components/Stepper/Stepper';
import { Button } from '@massalabs/react-ui-kit/src/components/Button/Button';
import { FiArrowRight } from 'react-icons/fi';
import LandingPage from '../../layouts/LandingPage/LandingPage';

export default function StepThree() {
  return (
    <LandingPage>
      <div
        className={`flex flex-col justify-center items-center align-center h-screen`}
      >
        <div className="flex flex-col justify-center items-center w-fit h-fit max-w-sm">
          <div className=" w-full max-w-xs mb-6">
            <Stepper step={2} steps={['Username', 'Password', 'Backup']} />
          </div>
          <div className="w-full">
            <h1 className="mas-banner text-neutral mb-6">Create an account</h1>
          </div>
          <div className="w-full mb-4">
            <p className="mas-body text-neutral">
              Security is very important, we recommend backing up your account
              using a key pair or a .yaml file.
              <br />
              Would you like to do it now or later?
            </p>
          </div>
          <div className="w-full mb-4">
            <Button posIcon={<FiArrowRight className="h-6 w-6" />}>Skip</Button>
          </div>
          <div className="w-full">
            <Button variant="secondary">Save now</Button>
          </div>
        </div>
      </div>
    </LandingPage>
  );
}
