import { useLocation } from 'react-router-dom';
import { usePost, useDelete } from '../../custom/api';
import { AccountObject } from '../../models/AccountModel';

import WalletLayout from '../../layouts/WalletLayout/WalletLayout';
import { Button } from '@massalabs/react-ui-kit';
import { FiTrash2 } from 'react-icons/fi';

export default function Settings() {
  const { state } = useLocation();
  const nickname: string = state.nickname;

  const { mutate: mutableBackup } = usePost<AccountObject>(
    `accounts/${nickname}/backup`,
  );

  const { mutate: mutableDelete } = useDelete<AccountObject>(
    `accounts/${nickname}`,
  );

  return (
    <WalletLayout>
      <div className="bg-secondary rounded-2xl w-full max-w-2xl p-10 m-auto">
        <p className="mas-body text-f-primary">Account security</p>
        <div className="flex gap-5 pt-5">
          <Button onClick={() => mutableBackup({} as AccountObject)}>
            Back up account
          </Button>
          <Button
            preIcon={<FiTrash2 />}
            onClick={() => mutableDelete({} as AccountObject)}
            variant="danger"
          >
            Delete account
          </Button>
        </div>
      </div>
    </WalletLayout>
  );
}
