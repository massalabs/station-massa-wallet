import { useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { FiLink } from 'react-icons/fi';
import QRCodeReact from 'qrcode.react';
import { maskAddress } from '../../../utils/MassaFormating';
import { AccountObject } from '../../../models/AccountModel';
import { useResource } from '../../../custom/api';
import { routeFor } from '../../../utils';
import CopyContent from './CopyContent';
import GenerateLink from './GenerateLink';
import Intl from '../../../i18n/i18n';
import { Button } from '@massalabs/react-ui-kit';

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
  const [modal, setModal] = useState(false);
  const presetURL = `${baseURL}${VITE_BASE_APP}/send-redirect/?to=${address}`;
  const [url, setURL] = useState(presetURL);

  const modalArgs = {
    account,
    presetURL,
    url,
    setURL,
    setModal,
  };

  return (
    <div className="flex flex-row items-center gap-5 mt-5">
      <QRCodeReact value={url} size={165} />
      <div className="flex flex-col w-full gap-3.5">
        <p>{Intl.t('receive.account-address')}</p>
        <CopyContent content={address} formattedContent={formattedAddress} />
        <Button onClick={() => setModal(!modal)} preIcon={<FiLink size={24} />}>
          {Intl.t('receive.receive-account')}
        </Button>
      </div>
      {modal ? <GenerateLink {...modalArgs} /> : null}
    </div>
  );
}

export default ReceiveCoins;
