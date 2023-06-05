import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { usePost, useDelete, usePut } from '../../custom/api';
import { useParams } from 'react-router-dom';
import { AccountObject } from '../../models/AccountModel';
import Intl from '../../i18n/i18n';
import { routeFor } from '../../utils';

import WalletLayout, {
  MenuItem,
} from '../../layouts/WalletLayout/WalletLayout';

import { Button, Identicon, Input } from '@massalabs/react-ui-kit';
import { FiTrash2 } from 'react-icons/fi';

export default function Settings() {
  const navigate = useNavigate();
  const { nickname } = useParams();

  const { mutate: mutableBackup } = usePost<AccountObject>(
    `accounts/${nickname}/backup`,
  );

  const { mutate: mutableDelete, isSuccess: isSuccessDelete } =
    useDelete<AccountObject>(`accounts/${nickname}`);

  if (isSuccessDelete) {
    navigate(routeFor('index'));
  }

  const { mutate: mutableUpdate } = usePut<AccountObject>(
    `accounts/${nickname}`,
  );

  const [newNickname, setNewNickname] = useState<string>(nickname || '');

  return (
    <WalletLayout menuItem={MenuItem.Settings}>
      <div className="flex flex-col justify-center items-center gap-9">
        <div className="bg-secondary rounded-2xl w-full max-w-2xl p-10">
          <p className="mas-body text-f-primary pb-5">
            {Intl.t('settings.title-profil')}
          </p>
          <div className="pb-7">
            <Identicon
              username={newNickname}
              size={120}
              className="bg-primary rounded-full"
            />
          </div>
          <div className="pb-5">
            <Input
              value={newNickname}
              onChange={(e) => setNewNickname(e.target.value)}
            />
          </div>
          <Button
            onClick={() =>
              mutableUpdate({ nickname: newNickname } as AccountObject)
            }
          >
            {Intl.t('settings.buttons.update')}
          </Button>
        </div>
        <div className="bg-secondary rounded-2xl w-full max-w-2xl p-10">
          <p className="mas-body text-f-primary pb-5">
            {Intl.t('settings.title-security')}
          </p>
          <div className="flex gap-5">
            <Button onClick={() => mutableBackup({} as AccountObject)}>
              {Intl.t('settings.buttons.backup')}
            </Button>
            <Button
              preIcon={<FiTrash2 />}
              onClick={() => mutableDelete({} as AccountObject)}
              variant="danger"
            >
              {Intl.t('settings.buttons.delete')}
            </Button>
          </div>
        </div>
      </div>
    </WalletLayout>
  );
}
