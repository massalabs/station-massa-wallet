import { useState } from 'react';

import { FT1, MassaLogo, Token } from '@massalabs/react-ui-kit';

import { XMA } from '@/const/assets/assets';
import { Asset } from '@/models/AssetModel';
import { DeleteAssetModal } from '@/pages/Assets/DeleteAssets';

export function AssetsList({ ...props }) {
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
            token.symbol === XMA ? <MassaLogo size={40} /> : <FT1 size={40} />
          }
          name={token.name}
          symbol={token.symbol}
          decimals={token.decimals}
          balance={token.balance}
          key={index}
          disable={token?.symbol === XMA ? true : false}
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
