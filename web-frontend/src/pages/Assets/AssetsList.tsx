import { useState } from 'react';

import { FT1, MassaLogo, Token } from '@massalabs/react-ui-kit';

import { MAS } from '@/const/assets/assets';
import { Asset } from '@/models/AssetModel';
import { DeleteAssetModal } from '@/pages/Assets/DeleteAssets';

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
      {assets?.map((token: Asset, index: number) => (
        <Token
          logo={
            token.symbol === MAS ? <MassaLogo size={40} /> : <FT1 size={40} />
          }
          name={token.name}
          symbol={token.symbol}
          decimals={token.decimals}
          balance={token.balance}
          key={index}
          disable={token?.symbol === MAS ? true : false}
          onDelete={() => {
            handleDelete(token.address);
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
