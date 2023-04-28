import LandingPage from '../layouts/LandingPage/LandingPage';

export default function Welcome() {
  return (
    <LandingPage>
      <div className="text-center">
        <h1 className="mas-banner">
          Welcome on
          <br />
          Massa<span className="text-brand">wallet</span>
        </h1>
        <button>Create an account</button>
        <button>Import an existing account</button>
      </div>
    </LandingPage>
  );
}
