import Placeholder from '../../layouts/Placeholder/Placeholder';
import WalletLayout, {
  MenuItem,
} from '../../layouts/WalletLayout/WalletLayout';
import Intl from '../../i18n/i18n';

function Contacts() {
  return (
    <WalletLayout menuItem={MenuItem.Contacts}>
      <Placeholder message={Intl.t('placeholder.teaser-contacts-page')} />
    </WalletLayout>
  );
}

export default Contacts;
