import Placeholder from '../../layouts/Placeholder/Placeholder';
import WalletLayout, {
  MenuItem,
} from '../../layouts/WalletLayout/WalletLayout';

function Transactions() {
  return (
    <WalletLayout menuItem={MenuItem.Transactions}>
      <Placeholder />
    </WalletLayout>
  );
}

export default Transactions;
