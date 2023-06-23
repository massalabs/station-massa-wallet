import { Link } from 'react-router-dom';
import Intl from '../i18n/i18n';
import { routeFor } from '../utils';
import LandingPage from '../layouts/LandingPage/LandingPage';

export default function Error() {
  return (
    <LandingPage>
      <div
        id="error-page"
        className="flex flex-col justify-center items-center h-screen text-f-primary"
      >
        <h1 className="mas-banner">
          {Intl.t('errors.unexpected-error.title')}
        </h1>
        <p className="mas-bod">
          {Intl.t('errors.unexpected-error.description')}
        </p>
        <Link to={routeFor('index')} className="underline">
          {Intl.t('errors.unexpected-error.link')}
        </Link>
      </div>
    </LandingPage>
  );
}
