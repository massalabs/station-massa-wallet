import LandingPage from '../../layouts/LandingPage/LandingPage';
import { FiUser, FiPlus } from 'react-icons/fi';
import { Button } from '@massalabs/react-ui-kit/src/components/Button/Button';
import { AccountSelectorButton } from '@massalabs/react-ui-kit/src/components/AccountSelector/AccountSelector';
import { MassaToken } from '@massalabs/react-ui-kit/src/components/Icons/Svg/SvgComponent/MassaToken';

export default function SelectAccount() {
  const selectProps = {
    avatar: <FiUser className="text-neutral h-6 w-6" />,
    accountName: 'account #',
    icon: <MassaToken size={24} />,
    amount: '0,000.00',
  };
  const buttonProps = {
    onClick: () => {
      console.log('clicked');
    },
    preIcon: <FiPlus />,
  };
  const defaultFlex = 'flex flex-col justify-center items-center align-center';
  return (
    <LandingPage>
      <div
        className={`${defaultFlex} 
                      h-screen`}
      >
        <div className="w-1/2">
          <div className="mb-6">
            <p className="mas-banner text-neutral"> Hey !</p>
          </div>
          <div className="mb-6">
            <label className="mas-body text-neutral" htmlFor="account-select">
              Select an account
            </label>
          </div>
          <div id="account-select" className="w-full flex flex-col">
            <div className="mb-4">
              <AccountSelectorButton {...selectProps} />
            </div>
            <div className="mb-4">
              <AccountSelectorButton {...selectProps} />
            </div>
            <div>
              <Button variant="secondary" {...buttonProps}>
                Add an account
              </Button>
            </div>
          </div>
        </div>
      </div>
    </LandingPage>
  );
}
