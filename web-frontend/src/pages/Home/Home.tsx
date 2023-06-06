import Placeholder from '../../layouts/Placeholder/Placeholder';
import WalletLayout, {
  MenuItem,
} from '../../layouts/WalletLayout/WalletLayout';
import Intl from '../../i18n/i18n';
import { FiHome } from 'react-icons/fi';
export default function Home() {
  return (
    <WalletLayout menuItem={MenuItem.Home}>
      <Placeholder
        message={Intl.t('placeholder.teaser-home-page')}
        icon={<FiHome size={114} />}
      />
    </WalletLayout>
  );
}
