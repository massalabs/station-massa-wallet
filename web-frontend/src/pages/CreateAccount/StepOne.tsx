import { Input } from '@massalabs/react-ui-kit/src/components/Input/Input';
import { Stepper } from '@massalabs/react-ui-kit/src/components/Stepper/Stepper';
import { Button } from '@massalabs/react-ui-kit/src/components/Button/Button';
import { FiArrowRight } from 'react-icons/fi';
import LandingPage from '../../layouts/LandingPage/LandingPage';

export default function StepOne() {
  return (
    <LandingPage>
      <div
        className={`flex flex-col justify-center items-center align-center h-screen`}
      >
        <div className="flex flex-col justify-center content-center items-center w-fit h-fit max-w-sm">
          <div className="w-full max-w-xs mb-6">
            <Stepper step={0} steps={['Username', 'Password', 'Back-up']} />{' '}
          </div>
          <h1 className="mas-banner text-neutral mb-6">Create an account</h1>
          <div className="w-full mb-4">
            <Input placeholder={'Enter an account name'} />
          </div>
          <div className="w-full">
            <Button posIcon={<FiArrowRight className="h-6 w-6" />}>Next</Button>
          </div>
        </div>
      </div>
    </LandingPage>
  );
}
