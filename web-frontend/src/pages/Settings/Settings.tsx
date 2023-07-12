import { useNavigate, useParams } from 'react-router-dom';
import { routeFor } from '@/utils';
import Intl from '@/i18n/i18n';

import { WalletLayout, MenuItem } from '@/layouts/WalletLayout/WalletLayout';

import { Button, Identicon, Input } from '@massalabs/react-ui-kit';
import { FiEdit } from 'react-icons/fi';
import { SettingsOption } from './SettingsOption';

export default function Settings() {
  const navigate = useNavigate();
  const { nickname } = useParams();

  return (
    <WalletLayout menuItem={MenuItem.Settings}>
      <div className="flex flex-col justify-center items-center gap-9 w-1/2">
        <div className="bg-secondary rounded-2xl w-full max-w-2xl p-10">
          <p className="mas-body text-f-primary pb-5">
            {Intl.t('settings.title-profil')}
          </p>
          <div className="pb-7">
            <Identicon
              username={nickname || 'username'}
              size={120}
              className="bg-primary rounded-full"
            />
          </div>
          <div className="pb-5">
            <Input
              value={nickname || 'username'}
              icon={<FiEdit />}
              onClickIcon={() =>
                navigate(routeFor(`${nickname}/settings/update`))
              }
              disabled
            />
          </div>
          <Button disabled>{Intl.t('settings.buttons.update')}</Button>
        </div>
        <SettingsOption nickname={nickname || 'username'} />
      </div>
    </WalletLayout>
  );
}
