import { useState } from 'react';

import { MassaLogo, Mrc20, Token } from '@massalabs/react-ui-kit';

import { DeleteAssetModal } from './DeleteAssetModal';
import { ITokenData, XMA } from '@/const/assets/assets';

export function AssetsList({ ...props }) {
  const { tokenArray, refetch } = props;

  const [tokenAddress, setTokenAddress] = useState<string>('');
  const [modalOpen, setModalOpen] = useState<boolean>(false);

  function handleDelete(address: string) {
    setTokenAddress(address);
    setModalOpen(true);
  }

  return (
    <>
      {tokenArray?.map((token: ITokenData, index: number) => (
        <Token
          logo={
            token.symbol === XMA ? <MassaLogo size={40} /> : <Mrc20 size={40} />
          }
          name={token.name}
          symbol={token.symbol}
          decimals={token.decimals}
          balance={token.balance}
          key={index}
          disable={token?.symbol === XMA ? true : false}
          onDelete={() => {
            handleDelete(token.assetAddress);
          }}
        />
      ))}
      {modalOpen && (
        <DeleteAssetModal
          tokenAddress={tokenAddress}
          setModalOpen={setModalOpen}
          refetch={refetch}
        />
      )}
    </>
  );
}
