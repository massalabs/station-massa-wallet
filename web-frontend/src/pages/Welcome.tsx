import LandingPage from '../layouts/LandingPage/LandingPage';
import { Button } from '@massalabs/react-ui-kit/src/components/Button/Button';
import { FiArrowRight } from 'react-icons/fi';

export default function Welcome() {
  return (
    <LandingPage>
      <div className="w-fit h-fit">
        <h1 className="mas-banner text-f-primary">
          Welcome on
          <br />
          Massa<span className="text-brand">wallet</span>
        </h1>
        <div className="pt-6">
          <Button posIcon={<FiArrowRight />}>Create an account</Button>
        </div>
        <div className="pt-3.5">
          <Button variant="secondary">Import an existing account</Button>
        </div>
      </div>
    </LandingPage>
  );
}
