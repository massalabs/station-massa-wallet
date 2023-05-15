import { Stepper } from '@massalabs/react-ui-kit/src/components/Stepper/Stepper';
import { Button } from '@massalabs/react-ui-kit/src/components/Button/Button';
import { FiLock } from 'react-icons/fi';
import LandingPage from '../../layouts/LandingPage/LandingPage';

export default function StepTwo() {
  return (
    <LandingPage>
      <div
        className={`flex flex-col justify-center items-center align-center h-screen`}
      >
        <div className="flex flex-col justify-center items-center w-fit h-fit max-w-sm">
          <div className=" w-full max-w-xs mb-6">
            <Stepper step={1} steps={['Username', 'Password', 'Back-up']} />
          </div>
          <div className="w-full">
            <h1 className="mas-banner text-neutral mb-6">Create an account</h1>
          </div>
          <div className="w-full mb-4">
            <p className="mas-body text-neutral">
              Now, define your password. To maximize security, a pop up will
              open directly on your device via MassaStation.
            </p>
          </div>
          <div className="w-full">
            <Button preIcon={<FiLock className="h-6 w-6" />}>
              Define a password
            </Button>
          </div>
        </div>
      </div>
    </LandingPage>
  );
}
