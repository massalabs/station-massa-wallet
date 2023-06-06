import Placeholder from '../../layouts/Placeholder/Placeholder';
import WalletLayout, {
  MenuItem,
} from '../../layouts/WalletLayout/WalletLayout';

function Assets() {
  return (
    <WalletLayout menuItem={MenuItem.Assets}>
      <Placeholder />
    </WalletLayout>
  );
}

export default Assets;
