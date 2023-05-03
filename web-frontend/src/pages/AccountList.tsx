import LandingPage from '../layouts/LandingPage/LandingPage';

export default function CreateAccount() {
  return (
    <LandingPage>
      <div>
        <h1 className="mas-banner">Hey!</h1>
        <h1 className="text-banner">Hey banner</h1>
        <p className="mas-body">Select an account</p>
        <button>Create an account</button>
        <button>Import an existing account</button>
      </div>
    </LandingPage>
  );
}
