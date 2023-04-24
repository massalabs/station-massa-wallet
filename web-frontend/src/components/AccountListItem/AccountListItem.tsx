import { MenuActive } from '@massalabs/react-ui-kit';
import AccountDefaultThumbnail from '../../assets/account-thumbnail-default.svg';
import { Amount } from '@massalabs/react-ui-kit';

export interface AccountListItemProps {
  name: string;
  amount: number;
}

export default function AccountListItem(props: AccountListItemProps) {
  return (
    <div className="flex flex-row justify-between align-center py-2 px-3 gap-48">
      <img src={AccountDefaultThumbnail} alt="account default thumbnail" />
      <MenuActive>{props.name}</MenuActive>
      <Amount amount={props.amount}></Amount>
    </div>
  );
}
