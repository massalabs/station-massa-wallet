import { Button } from '@massalabs/react-ui-kit';
import { FiTrash2 } from 'react-icons/fi';
import { useNavigate } from 'react-router-dom';

import { usePost, useDelete } from '@/custom/api';
import Intl from '@/i18n/i18n';
import { AccountObject } from '@/models/AccountModel';
import { routeFor } from '@/utils';

interface SettingsOptionProps {
  nickname: string;
}

export function SettingsOption(props: SettingsOptionProps) {
  const { nickname } = props;

  const navigate = useNavigate();

  const { mutate: mutableBackup } = usePost<AccountObject>(
    `accounts/${nickname}/backup`,
  );

  const { mutate: mutableDelete, isSuccess: isSuccessDelete } =
    useDelete<AccountObject>(`accounts/${nickname}`);

  if (isSuccessDelete) {
    navigate(routeFor('index'));
  }

  return (
    <div className="bg-secondary rounded-2xl w-full max-w-2xl p-10">
      <h2 className="mas-h2 text-f-primary pb-5">
        {Intl.t('settings.title-security')}
      </h2>
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
  );
}
