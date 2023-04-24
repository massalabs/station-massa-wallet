import { Banner } from '@massalabs/react-ui-kit';
import LandingPage from '../layouts/LandingPage/LandingPage';
import { Body } from '@massalabs/react-ui-kit';

export default function CreateAccount() {
  return (
    <LandingPage>
      <div>
        <Banner>Hey!</Banner>
        <h1 className="text-banner">Hey banner</h1>
        <Body>Select an account</Body>
        <button>Create an account</button>
        <button>Import an existing account</button>
      </div>
    </LandingPage>
  );
}
