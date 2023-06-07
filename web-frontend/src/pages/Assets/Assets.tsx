import Placeholder from '../../layouts/Placeholder/Placeholder';
import WalletLayout, {
  MenuItem,
} from '../../layouts/WalletLayout/WalletLayout';
import Intl from '../../i18n/i18n';
import { FiDisc } from 'react-icons/fi';

function Assets() {
  return (
    <WalletLayout menuItem={MenuItem.Assets}>
      <Placeholder
        message={Intl.t('placeholder.teaser-assets-page')}
        icon={<FiDisc size={114} />}
      />
    </WalletLayout>
  );
}

export default Assets;
