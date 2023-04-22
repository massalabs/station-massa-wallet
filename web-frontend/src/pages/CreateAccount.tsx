import Banner from '../components/Banner';
import LandingPage from '../layouts/LandingPage/LandingPage';

export default function CreateAccount() {
  return (
    <LandingPage>
      <div className="">
        {/* stepper header */}
        <Banner>Create an account</Banner>
        {/* stepper content */}
      </div>
    </LandingPage>
  );
}
