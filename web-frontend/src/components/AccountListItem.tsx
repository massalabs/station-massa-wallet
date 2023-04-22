import MenuActive from './MenuActive';
import AccountDefaultThumbnail from '../assets/account-default-thumbnail.svg';
import IcoAmountDark from '../assets/ico-amount-dark.svg';
import Amount from './Amount';

interface AccountListProps {
  name: string;
  amount: number;
}

export default function AccountListItem(props: AccountListProps) {
  return (
    <div className="flex flex-row justify-between align-center py-2 px-3 gap-48">
      <img src={AccountDefaultThumbnail} alt="account default thumbnail" />
      <MenuActive>{props.name}</MenuActive>
      <img src={IcoAmountDark} alt="picto amount" />
      <Amount amount={props.amount}></Amount>
    </div>
  );
}
