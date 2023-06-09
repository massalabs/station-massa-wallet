import { useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { AccountObject } from '../../models/AccountModel';
import { formatStandard, Unit, maskAddress } from '../../utils/MassaFormating';
import { useResource } from '../../custom/api';
import { routeFor } from '../../utils';
import Intl from '../../i18n/i18n';

import WalletLayout, {
  MenuItem,
} from '../../layouts/WalletLayout/WalletLayout';
import { Button, Balance } from '@massalabs/react-ui-kit';
import { FiArrowDownLeft, FiArrowUpRight, FiCopy } from 'react-icons/fi';

export default function Home() {
  const navigate = useNavigate();
  const { nickname } = useParams();
  const {
    error,
    data: account,
    status,
  } = useResource<AccountObject>(`accounts/${nickname}`);

  const isLoadingData = status === 'loading';

  useEffect(() => {
    if (error) {
      navigate('/error');
    } else if (!account && !isLoadingData) {
      navigate(routeFor('account-select'));
    }
  }, [account, navigate]);

  const unformattedBalance = account?.candidateBalance ?? '0';
  const balance = parseInt(unformattedBalance);
  const formattedBalance = formatStandard(balance, Unit.NanoMAS, 2);
  const address = account?.address ?? '';
  const formattedAddress = maskAddress(address);

  return (
    <WalletLayout menuItem={MenuItem.Home}>
      <div className="flex flex-col justify-center items-center gap-5 w-1/2">
        <div className="bg-secondary rounded-2xl w-full max-w-lg p-10">
          <p className="mas-body text-f-primary mb-2">
            {Intl.t('home.title-account-balance')}
          </p>
          <Balance size="lg" amount={formattedBalance} customClass="mb-6" />
          <div className="flex gap-7">
            <Button
              variant="secondary"
              preIcon={<FiArrowDownLeft />}
              onClick={() =>
                navigate(routeFor(`${nickname}/send-coins?tabIndex=1`))
              }
            >
              {Intl.t('home.buttons.receive')}
            </Button>
            <Button
              preIcon={<FiArrowUpRight />}
              onClick={() =>
                navigate(routeFor(`${nickname}/send-coins?tabIndex=0`))
              }
            >
              {Intl.t('home.buttons.send')}
            </Button>
          </div>
        </div>
        <div className="bg-secondary rounded-2xl w-full max-w-lg p-10">
          <p className="mas-body text-f-primary mb-6">
            {Intl.t('home.title-account-address')}
          </p>
          <div
            className="flex flex-row items-center mas-body2 justify-between
              w-full h-12 px-3 rounded bg-primary cursor-pointer"
            onClick={() => navigator.clipboard.writeText(address)}
          >
            <u>{formattedAddress}</u>
            <FiCopy size={24} />
          </div>
        </div>
      </div>
    </WalletLayout>
  );
}
