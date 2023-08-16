import { useState } from 'react';

import { Button } from '@massalabs/react-ui-kit';
import { FiPlus } from 'react-icons/fi';
import { useParams } from 'react-router-dom';

import { useResource } from '@/custom/api';
import Intl from '@/i18n/i18n';
import { WalletLayout, MenuItem } from '@/layouts/WalletLayout/WalletLayout';
import { Asset } from '@/models/AssetModel';
import { AssetsLoading, AssetsImportModal, AssetsList } from '@/pages/Assets';

function Assets() {
  const [modal, setModal] = useState(false);
  const { nickname } = useParams();
  const { data: assets, isLoading: isGetLoading } = useResource<Asset[]>(
    `accounts/${nickname}/assets`,
  );

  return (
    <WalletLayout menuItem={MenuItem.Assets}>
      <div className="bg-secondary p-10 h-fit w-2/3 rounded-lg z-0">
        <div className="flex items-center justify-between w-full mb-6">
          <div> {Intl.t('assets.title')}</div>
          <Button
            customClass="w-fit"
            preIcon={<FiPlus size={114} />}
            onClick={() => {
              setModal(true);
            }}
          >
            {Intl.t('assets.import.import-button')}
          </Button>
        </div>
        <div className="flex flex-col w-full h-fit bg-primary rounded-lg gap-4 p-8">
          {isGetLoading ? <AssetsLoading /> : <AssetsList assets={assets} />}
        </div>
        {modal && <AssetsImportModal closeModal={() => setModal(false)} />}
      </div>
    </WalletLayout>
  );
}

export default Assets;
