import { Button } from '@massalabs/react-ui-kit';
import CopyAddress from './CopyAddress';
import { FiLink } from 'react-icons/fi';
import Intl from '../../../i18n/i18n';
import QRCodeReact from 'qrcode.react';
import { maskAddress } from '../../../utils/MassaFormating';
import { AccountObject } from '../../../models/AccountModel';
import { useResource } from '../../../custom/api';
import { useNavigate, useParams } from 'react-router-dom';
import { routeFor } from '../../../utils';

function ReceiveCoins() {
  const navigate = useNavigate();
  const { nickname } = useParams();
  const { data: account } = useResource<AccountObject>(`accounts/${nickname}`);

  if (account === undefined) {
    navigate(routeFor(`${nickname}/home`));
    return null;
  }

  const address = account.address;
  const formattedAddress = maskAddress(address);
  const VITE_BASE_APP = import.meta.env.VITE_BASE_APP;
  const baseURL = window.location.origin;

  const presetURL = `${baseURL}${VITE_BASE_APP}/send-redirect/?to=${address}`;

  return (
    <div className="flex flex-row items-center gap-5">
      <QRCodeReact value={presetURL} size={165} />
      <div className="flex flex-col w-full gap-3.5">
        <p>{Intl.t('receive.account-address')}</p>
        <CopyAddress address={address} formattedAddress={formattedAddress} />
        <Button preIcon={<FiLink size={24} />}>
          {Intl.t('receive.account-receive')}
        </Button>
      </div>
    </div>
  );
}

export default ReceiveCoins;
