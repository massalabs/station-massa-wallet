import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { AccountObject } from '../../models/AccountModel';
import { formatStandard, Unit, maskAddress } from '../../utils/massaFormat';
import { useResource } from '../../custom/api';
import { routeFor } from '../../utils';
import Intl from '../../i18n/i18n';
import { TAB_SEND, TAB_RECEIVE } from '../../const/tabs/tabs';
import { Loading } from './Loading';

import WalletLayout, {
  MenuItem,
} from '../../layouts/WalletLayout/WalletLayout';
import { Button, Balance, Clipboard } from '@massalabs/react-ui-kit';
import { FiArrowDownLeft, FiArrowUpRight } from 'react-icons/fi';

export default function Home() {
  const navigate = useNavigate();
  const { nickname } = useParams();
  const {
    error,
    data: account,
    isLoading,
  } = useResource<AccountObject>(`accounts/${nickname}`);

  const breakpoint = 1920; // Screen breakpoint value
  const unformattedBalance = account?.candidateBalance ?? '0';
  const balance = parseInt(unformattedBalance);
  const formattedBalance = formatStandard(balance, Unit.NanoMAS);
  const address = account?.address ?? '';
  const formattedAddress = maskAddress(address);

  const [screenWidth, setScreenWidth] = useState(window.innerWidth);
  const [displayedContent, setDisplayedContent] = useState(formattedAddress);

  // Function seems unnecessary,
  // but it is needed to update the state of screenWidth
  function handleResize() {
    setScreenWidth(window.innerWidth);
  }

  useEffect(() => {
    if (error) {
      navigate(routeFor('error'));
    } else if (!account && !isLoading) {
      navigate(routeFor('account-select'));
    }
  }, [account, navigate]);

  // Call handleResize once after component mounts to log formattedAddress on page load
  useEffect(() => {
    window.addEventListener('resize', handleResize);
    handleResize();
    return () => {
      window.removeEventListener('resize', handleResize);
    };
  }, []);

  useEffect(() => {
    screenWidth <= breakpoint
      ? setDisplayedContent(formattedAddress)
      : setDisplayedContent(address);
  }, [screenWidth, breakpoint, address, formattedAddress]);

  return (
    <WalletLayout menuItem={MenuItem.Home}>
      {isLoading ? (
        <Loading />
      ) : (
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
                  navigate(
                    routeFor(`${nickname}/transfer-coins?tab=${TAB_RECEIVE}`),
                  )
                }
              >
                {Intl.t('home.buttons.receive')}
              </Button>
              <Button
                preIcon={<FiArrowUpRight />}
                onClick={() =>
                  navigate(
                    routeFor(`${nickname}/transfer-coins?tab=${TAB_SEND}`),
                  )
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
            <Clipboard
              displayedContent={displayedContent}
              rawContent={address}
              error={Intl.t('errors.no-content-to-copy')}
              className="flex flex-row items-center mas-body2 justify-between
              w-full h-12 px-3 rounded bg-primary cursor-pointer"
            />
          </div>
        </div>
      )}
    </WalletLayout>
  );
}
