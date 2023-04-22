import Banner from '../components/Banner/Banner';
import LandingPage from '../layouts/LandingPage/LandingPage';

export default function Welcome() {
  return (
    <LandingPage>
      <div className="">
        <Banner>
          Welcome on
          <br />
          Massa<span className="text-green">wallet</span>
        </Banner>
        <button>Create an account</button>
        <button>Import an existing account</button>
      </div>
    </LandingPage>
  );
}
