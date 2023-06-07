import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { usePut } from '../../custom/api';
import { useParams } from 'react-router-dom';
import { AccountObject } from '../../models/AccountModel';
import { routeFor } from '../../utils';
import Intl from '../../i18n/i18n';

import WalletLayout, {
  MenuItem,
} from '../../layouts/WalletLayout/WalletLayout';

import { Button, Identicon, Input } from '@massalabs/react-ui-kit';
import { SettingsOption } from './SettingsOption';

export default function SettingsEdit() {
  const navigate = useNavigate();
  const { nickname } = useParams();

  const [newNickname, setNewNickname] = useState<string>(nickname || '');

  const { mutate: mutableUpdate, isSuccess: isSuccessUpdate } =
    usePut<AccountObject>(`accounts/${nickname}`);

  if (isSuccessUpdate) {
    navigate(routeFor(`${newNickname}/settings`));
  }

  return (
    <WalletLayout menuItem={MenuItem.Settings}>
      <div className="flex flex-col justify-center items-center gap-9 w-1/2">
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
        <SettingsOption nickname={newNickname} />
      </div>
    </WalletLayout>
  );
}
