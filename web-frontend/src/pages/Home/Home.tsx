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
import Intl from '@/i18n/i18n';
import { WalletLayout, MenuItem } from '@/layouts/WalletLayout/WalletLayout';
import { AccountObject } from '@/models/AccountModel';
import { routeFor } from '@/utils';

export default function Home() {
  const navigate = useNavigate();
  const { nickname } = useParams();

  const {
    error,
    data: account,
    isLoading,
  } = useResource<AccountObject>(`accounts/${nickname}`);

  const { reverseResolveDns, domainNameList, resetDomainList } = useMNS();
  const accountAddress = account?.address;

  useEffect(() => {
    if (!isLoading && accountAddress) {
      resetDomainList();
      reverseResolveDns(accountAddress);
    }
  }, [reverseResolveDns, accountAddress, resetDomainList, isLoading]);

  useEffect(() => {
    if (error) {
      navigate(routeFor('error'));
    } else if (!account && !isLoading) {
      navigate(routeFor('account-select'));
    }
  }, [account, navigate, error, isLoading]);

  const unformattedBalance = account?.candidateBalance ?? '0';
  const formattedBalance = formatAmount(unformattedBalance).preview;

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
            {accountAddress && (
              <div className="flex w-full justify-between items-center">
                <Clipboard
                  displayedContent={accountAddress}
                  rawContent={accountAddress}
                  error={Intl.t('errors.no-content-to-copy')}
                  className="flex flex-row items-center mas-body2 justify-between
              w-full h-12 px-3 rounded bg-primary cursor-pointer"
                />
              </div>
            )}
          </div>
          {domainNameList && domainNameList.length > 0 && (
            <div className="bg-secondary rounded-2xl w-full max-w-lg p-10 gap-6 flex flex-col justify-between">
              <div className="flex flex-col">
                <p className="mas-body text-f-primary">
                  {Intl.t('home.title-mns')}
                </p>
                <p className="mas-caption">{Intl.t('home.desc-mns')}</p>
              </div>
              <div className="flex flex-col w-full gap-4 max-h-[120px] overflow-y-scroll ">
                {domainNameList?.map((domainName) => (
                  <Mns mns={domainName} />
                ))}
              </div>
            </div>
          )}
        </div>
      )}
    </WalletLayout>
  );
}
