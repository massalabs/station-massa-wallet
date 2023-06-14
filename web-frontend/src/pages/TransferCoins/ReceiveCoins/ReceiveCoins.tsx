import { useState } from 'react';
import { FiLink } from 'react-icons/fi';
import QRCodeReact from 'qrcode.react';
import { maskAddress } from '../../../utils/massaFormating';
import CopyContent from './CopyContent';
import GenerateLink from './GenerateLink';
import Intl from '../../../i18n/i18n';
import { Button } from '@massalabs/react-ui-kit';

const VITE_BASE_APP = import.meta.env.VITE_BASE_APP;

function ReceiveCoins({ ...props }) {
  const { account } = props;

  const address = account?.address;
  const formattedAddress = maskAddress(address);
  const baseURL = window.location.origin;
  const presetURL = `${baseURL}${VITE_BASE_APP}/send-redirect/?to=${address}`;

  const [modal, setModal] = useState(false);
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
        <p>{Intl.t('receive-coins.account-address')}</p>
        <CopyContent content={address} formattedContent={formattedAddress} />
        <Button onClick={() => setModal(!modal)} preIcon={<FiLink size={24} />}>
          {Intl.t('receive-coins.receive-account')}
        </Button>
      </div>
      {modal ? <GenerateLink {...modalArgs} /> : null}
    </div>
  );
}

export default ReceiveCoins;
