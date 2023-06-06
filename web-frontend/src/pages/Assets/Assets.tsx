import Placeholder from '../../layouts/Placeholder/Placeholder';
import WalletLayout, {
  MenuItem,
} from '../../layouts/WalletLayout/WalletLayout';
import Intl from '../../i18n/i18n';

function Assets() {
  return (
    <WalletLayout menuItem={MenuItem.Assets}>
      <Placeholder message={Intl.t('placeholder.teaser-assets-page')} />
    </WalletLayout>
  );
}

export default Assets;
