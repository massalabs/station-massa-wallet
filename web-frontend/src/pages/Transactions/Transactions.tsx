import Placeholder from '../../layouts/Placeholder/Placeholder';
import WalletLayout, {
  MenuItem,
} from '../../layouts/WalletLayout/WalletLayout';
import Intl from '../../i18n/i18n';

function Transactions() {
  return (
    <WalletLayout menuItem={MenuItem.Transactions}>
      <Placeholder message={Intl.t('placeholder.teaser-transaction-page')} />
    </WalletLayout>
  );
}

export default Transactions;
