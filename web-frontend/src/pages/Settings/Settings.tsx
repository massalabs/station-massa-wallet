import { Button, Identicon, Input } from '@massalabs/react-ui-kit';
import { FiEdit } from 'react-icons/fi';
import { useNavigate, useParams } from 'react-router-dom';

import { SettingsOption } from './SettingsOption';
import SettingsSignRules from './SettingsSignRule';
import Intl from '@/i18n/i18n';
import { WalletLayout, MenuItem } from '@/layouts/WalletLayout/WalletLayout';
import { routeFor } from '@/utils';

export default function Settings() {
  const navigate = useNavigate();
  const { nickname } = useParams();

  return (
    <WalletLayout menuItem={MenuItem.Settings}>
      <div className="w-full max-h-screen overflow-y-auto p-10">
        <div className="flex flex-col justify-center items-center gap-9">
          <div className="bg-secondary rounded-2xl w-full max-w-2xl p-10">
            <p className="mas-body text-f-primary pb-5">
              {Intl.t('settings.title-profile')}
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
          <div className="bg-secondary rounded-2xl w-full max-w-2xl p-10">
            <SettingsSignRules nickname={nickname || 'username'} />
          </div>
        </div>
      </div>
    </WalletLayout>
  );
}
