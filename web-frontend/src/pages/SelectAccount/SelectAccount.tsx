import LandingPage from '../../layouts/LandingPage/LandingPage';
import { FiUser, FiPlus } from 'react-icons/fi';
import { Button } from '@massalabs/react-ui-kit/src/components/Button/Button';
import { AccountSelectorButton } from '@massalabs/react-ui-kit/src/components/AccountSelector/AccountSelector';
import { MassaToken } from '@massalabs/react-ui-kit/src/components/Icons/Svg/SvgComponent/MassaToken';
import { useQuery } from '@tanstack/react-query';
import { getAllAccounts } from '../../api/account';
import { accountType } from '../../api/types';

export default function SelectAccount() {
  const accounts = useQuery({
    queryKey: ['accounts'],
    queryFn: getAllAccounts,
  });

  const accountBaseProps = {
    avatar: <FiUser className="text-neutral h-6 w-6" />,
    icon: <MassaToken size={24} />,
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
            {accounts.data?.map((account: accountType) => (
              <div className="mb-4">
                <AccountSelectorButton
                  {...accountBaseProps}
                  accountName={account.nickname}
                  amount={account.candidateBalance}
                />
              </div>
            ))}
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
