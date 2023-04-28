import LandingPage from '../layouts/LandingPage/LandingPage';

export default function CreateAccount() {
  return (
    <LandingPage>
      <div className="">
        <h1 className="mas-banner">Add an account</h1>
        <button>Create an account</button>
        <button>Import an existing account</button>
      </div>
    </LandingPage>
  );
}
