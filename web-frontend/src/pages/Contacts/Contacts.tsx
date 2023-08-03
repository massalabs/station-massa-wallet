import { FiUsers } from 'react-icons/fi';

import Intl from '@/i18n/i18n';
import Placeholder from '@/layouts/Placeholder/Placeholder';
import { WalletLayout, MenuItem } from '@/layouts/WalletLayout/WalletLayout';

function Contacts() {
  return (
    <WalletLayout menuItem={MenuItem.Contacts}>
      <Placeholder
        message={Intl.t('placeholder.teaser-contacts-page')}
        icon={<FiUsers size={114} />}
      />
    </WalletLayout>
  );
}

export default Contacts;
