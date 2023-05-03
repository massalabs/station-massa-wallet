import AccountDefaultThumbnail from '../../assets/account-thumbnail-default.svg';

export interface AccountListItemProps {
  name: string;
  amount: number;
}

export default function AccountListItem(props: AccountListItemProps) {
  return (
    <div className="flex flex-row justify-between align-center py-2 px-3 gap-48">
      <img src={AccountDefaultThumbnail} alt="account default thumbnail" />
      <p className="body">{props.name}</p>
    </div>
  );
}
