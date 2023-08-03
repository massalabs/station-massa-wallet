import { WalletLayout, MenuItem } from '@/layouts/WalletLayout/WalletLayout';

import Intl from '@/i18n/i18n';

import { FiPlus } from 'react-icons/fi';

import { Button } from '@massalabs/react-ui-kit';

import { useState } from 'react';

import { useResource } from '@/custom/api';

import { useParams } from 'react-router-dom';

import { IToken } from '@/models/AccountModel';

import { AssetsImportModal } from './AssetsImportModal';

import { AssetsList } from './AssetsList';

import { AssetsLoading } from './AssetsLoading';

function Assets() {
  const [modal, setModal] = useState(false);
  const { nickname } = useParams();
  const { data: tokenArray = [], isLoading: isGetLoading } = useResource<
    IToken[]
  >(`accounts/${nickname}/assets`);

  return (
    <WalletLayout menuItem={MenuItem.Assets}>
      <div className="bg-secondary p-10 h-fit w-2/3 rounded-lg">
        <div className="flex items-center justify-between w-full mb-6">
          <div> {Intl.t('assets.title')}</div>
          <Button
            customClass="w-fit"
            preIcon={<FiPlus size={114} />}
            onClick={() => {
              setModal(true);
            }}
          >
            {Intl.t('assets.import')}
          </Button>
        </div>
        <div className="flex flex-col w-full h-fit bg-primary rounded-lg gap-4 p-8">
          {isGetLoading ? (
            <AssetsLoading />
          ) : (
            <AssetsList tokenArray={tokenArray} />
          )}
        </div>
        {modal && <AssetsImportModal setModal={setModal} />}
      </div>
    </WalletLayout>
  );
}

export default Assets;
