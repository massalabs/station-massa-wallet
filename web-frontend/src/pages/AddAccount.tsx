import Banner from '../components/Banner';
import LandingPage from '../layouts/LandingPage/LandingPage';

export default function CreateAccount() {
  return (
    <LandingPage>
      <div className="">
        <Banner>Add an account</Banner>
        <button>Create an account</button>
        <button>Import an existing account</button>
      </div>
    </LandingPage>
  );
}
