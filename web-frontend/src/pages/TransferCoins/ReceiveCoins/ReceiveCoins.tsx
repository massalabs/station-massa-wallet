import { Button } from '@massalabs/react-ui-kit';
import CopyAddress from './CopyAddress';
import { FiLink } from 'react-icons/fi';
import Intl from '../../../i18n/i18n';

function ReceiveCoins({ ...props }) {
  return (
    <div className="flex flex-row items-center gap-5">
      {/* Placeholder image that will be replaced by a QR code */}
      <img src="https://placehold.jp/165x165.png"></img>
      <div className="flex flex-col w-full gap-3.5">
        <p>{Intl.t('receive.account-address')}</p>
        <CopyAddress {...props} />
        <Button preIcon={<FiLink size={24} />}>
          {Intl.t('receive.account-receive')}
        </Button>
      </div>
    </div>
  );
}

export default ReceiveCoins;
