import Intl from '@/i18n/i18n';
import LandingPage from '@/layouts/LandingPage/LandingPage';

export function Loading() {
  return (
    <LandingPage>
      <div className="flex flex-col justify-center items-center h-screen">
        <div className="flex flex-col justify-center items-start w-full h-full max-w-lg">
          <h1 className="mas-banner text-f-primary pb-6 animate-pulse blur-sm">
            {Intl.t('account.header.title')}
          </h1>
          <label
            className="mas-body text-info pb-6 animate-pulse blur-sm"
            htmlFor="account-select"
          >
            {Intl.t('account.select')}
          </label>
          <div id="account-select" className="pb-4 w-full">
            <div className="flex items-center space-x-4 animate-pulse">
              <div className="flex-1 py-1">
                <div className="h-14 bg-c-disabled-2 rounded-lg mb-3"></div>
                <div className="h-14 bg-c-disabled-2 rounded-lg mb-3"></div>
                <div className="h-14 bg-tertiary rounded-lg mb-3"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </LandingPage>
  );
}
