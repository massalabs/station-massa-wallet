import MenuActive from '../MenuActive';
import AccountDefaultThumbnail from '../../assets/account-thumbnail-default.svg';
import IcoAmountDark from '../../assets/ico-amount-dark.svg';
import Amount from '../Amount';

export interface AccountListItemProps {
  name: string;
  amount: number;
}

export default function AccountListItem(props: AccountListItemProps) {
  return (
    <div className="flex flex-row justify-between align-center py-2 px-3 gap-48">
      <img src={AccountDefaultThumbnail} alt="account default thumbnail" />
      <MenuActive>{props.name}</MenuActive>
      <img src={IcoAmountDark} alt="picto amount" />
      <Amount amount={props.amount}></Amount>
    </div>
  );
}
