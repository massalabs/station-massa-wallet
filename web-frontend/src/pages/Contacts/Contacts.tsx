import Placeholder from '../../layouts/Placeholder/Placeholder';
import WalletLayout, {
  MenuItem,
} from '../../layouts/WalletLayout/WalletLayout';

function Contacts() {
  return (
    <WalletLayout menuItem={MenuItem.Contacts}>
      <Placeholder />
    </WalletLayout>
  );
}

export default Contacts;
