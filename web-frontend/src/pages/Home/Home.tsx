import { useParams } from 'react-router-dom';
import { AccountObject } from '../../models/AccountModel';
import { formatStandard, Unit } from '../../utils/MassaFormating';
import { useResource } from '../../custom/api';
import Intl from '../../i18n/i18n';

import WalletLayout, {
  MenuItem,
} from '../../layouts/WalletLayout/WalletLayout';
import { Button, Balance } from '@massalabs/react-ui-kit';
import { FiArrowDownLeft, FiArrowUpRight } from 'react-icons/fi';

export default function Home() {
  const { nickname } = useParams();
  const { data: account } = useResource<AccountObject>(`accounts/${nickname}`);

  const unformattedBalance = account?.candidateBalance ?? '0';
  const balance = parseInt(unformattedBalance);
  const formattedBalance = formatStandard(balance, Unit.NanoMAS, 2);

  return (
    <WalletLayout menuItem={MenuItem.Home}>
      <div className="flex flex-col justify-center items-center gap-9 w-1/2">
        <div className="bg-secondary rounded-2xl w-full max-w-lg p-10">
          <p className="mas-body text-f-primary mb-2">
            {Intl.t('home.title-account-balance')}
          </p>
          <Balance size="lg" amount={formattedBalance} customClass="mb-6" />
          <div className="flex gap-7">
            <Button variant="secondary" preIcon={<FiArrowDownLeft />}>
              {Intl.t('home.buttons.receive')}
            </Button>
            <Button preIcon={<FiArrowUpRight />}>
              {Intl.t('home.buttons.send')}
            </Button>
          </div>
        </div>
      </div>
    </WalletLayout>
  );
}
