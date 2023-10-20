import { useState } from 'react';

import { Button, Clipboard } from '@massalabs/react-ui-kit';
import QRCodeReact from 'qrcode.react';
import { FiLink } from 'react-icons/fi';

import GenerateLink from './GenerateLink';
import Intl from '@/i18n/i18n';

const VITE_BASE_APP = import.meta.env.VITE_BASE_APP;

function ReceiveCoins({ ...props }) {
  const { account } = props;

  const address = account?.address;
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
    <div
      className="flex flex-row items-center gap-5 mt-5"
      data-testid="receive-coins"
    >
      <QRCodeReact value={url} size={165} />
      <div className="flex flex-col w-full gap-3.5">
        <p>{Intl.t('receive-coins.account-address')}</p>
        <div className="h-16">
          <Clipboard
            displayedContent={address}
            rawContent={address}
            toggleHover={false}
            error={Intl.t('errors.no-content-to-copy')}
          />
        </div>
        <Button onClick={() => setModal(!modal)} preIcon={<FiLink size={24} />}>
          {Intl.t('receive-coins.receive-account')}
        </Button>
      </div>
      {modal ? <GenerateLink {...modalArgs} /> : null}
    </div>
  );
}

export default ReceiveCoins;
