import { useEffect } from 'react';

import {
  Button,
  Balance,
  Clipboard,
  formatAmount,
  Mns,
} from '@massalabs/react-ui-kit';
import { FiArrowDownLeft, FiArrowUpRight } from 'react-icons/fi';
import { useNavigate, useParams } from 'react-router-dom';

import { Loading } from './Loading';
import { TAB_SEND, TAB_RECEIVE } from '@/const/tabs/tabs';
import { useResource } from '@/custom/api';
import { useMNS } from '@/custom/useMNS';
import { usePrepareScCall } from '@/custom/usePrepareScCall';
import Intl from '@/i18n/i18n';
import { WalletLayout, MenuItem } from '@/layouts/WalletLayout/WalletLayout';
import { AccountObject } from '@/models/AccountModel';
import { maskAddress, routeFor } from '@/utils';

export default function Home() {
  const navigate = useNavigate();
  const { nickname } = useParams();
  const {
    error,
    data: account,
    isLoading,
  } = useResource<AccountObject>(`accounts/${nickname}`);
  const { client } = usePrepareScCall();
  const { reverseResolveDns, domainName } = useMNS(client);
  const accountAddress = account?.address ?? '';

  useEffect(() => {
    reverseResolveDns(accountAddress);
  }, [reverseResolveDns, nickname]);

  useEffect(() => {
    if (error) {
      navigate(routeFor('error'));
    } else if (!account && !isLoading) {
      navigate(routeFor('account-select'));
    }
  }, [account, navigate, error, isLoading]);

  const unformattedBalance = account?.candidateBalance ?? '0';
  const balance = parseInt(unformattedBalance);
  const formattedBalance = formatAmount(
    balance.toString(),
  ).amountFormattedPreview;

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
                data-testid="receive-button"
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
                data-testid="send-button"
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

          <div className="bg-secondary rounded-2xl w-full max-w-lg p-10 flex flex-col justify-between">
            <p className="mas-body text-f-primary mb-6">
              {Intl.t('home.title-account-address')}
            </p>
            <div className="flex w-full justify-between items-center">
              {domainName && <Mns mns={domainName} />}
              <Clipboard
                displayedContent={maskAddress(accountAddress)}
                rawContent={accountAddress}
                error={Intl.t('errors.no-content-to-copy')}
                className="flex flex-row items-center mas-body2 justify-between
              w-fit h-12 px-3 rounded bg-primary cursor-pointer"
              />
            </div>
          </div>
        </div>
      )}
    </WalletLayout>
  );
}
