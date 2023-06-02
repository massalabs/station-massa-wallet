import WalletLayout, {
  MenuItem,
} from '../../layouts/WalletLayout/WalletLayout';
export default function Home() {
  return (
    <WalletLayout menuItem={MenuItem.Home}>
      <div>
        <h1>Home Page</h1>
      </div>
    </WalletLayout>
  );
}
