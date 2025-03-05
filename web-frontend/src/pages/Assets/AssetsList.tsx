import { useState } from 'react';

import { Token } from '@massalabs/react-ui-kit';

import Intl from '@/i18n/i18n';
import { Asset } from '@/models/AssetModel';
import { DeleteAssetModal } from '@/pages/Assets/DeleteAssets';
import { tokenIcon } from '@/utils/tokenIcon';

interface AssetsListProps {
  assets: Asset[] | undefined;
}

export function AssetsList(props: AssetsListProps) {
  const { assets } = props;

  const [tokenAddress, setTokenAddress] = useState<string>('');
  const [modal, setModal] = useState<boolean>(false);

  function handleDelete(address: string) {
    setTokenAddress(address);
    setModal(true);
  }

  return (
    <>
      {assets
        ?.filter((a) => a.balance !== undefined && a.balance !== '')
        .map((token: Asset, index: number) => (
          <Token
            logo={tokenIcon(token.symbol, 32)}
            name={token.name}
            symbol={token.symbol}
            decimals={token.decimals}
            balance={token.balance || ''}
            dollarValue={token.dollarValue}
            dollarValueError={
              token.dollarValue === undefined || token.dollarValue === ''
                ? Intl.t('errors.dollar-value')
                : undefined
            }
            key={index}
            address={token.address}
            disable={token?.isDefault}
            onDelete={() => {
              if (token.address) {
                handleDelete(token.address);
              }
            }}
          />
        ))}
      {modal && (
        <DeleteAssetModal
          tokenAddress={tokenAddress}
          closeModal={() => setModal(false)}
        />
      )}
    </>
  );
}
