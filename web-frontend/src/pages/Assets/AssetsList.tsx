import { useState } from 'react';

import { Token, getAssetIcons } from '@massalabs/react-ui-kit';
import { useParams } from 'react-router-dom';

import { MAS } from '@/const/assets/assets';
import { useFTTransfer } from '@/custom/smart-contract/useFTTransfer';
import { Asset } from '@/models/AssetModel';
import { DeleteAssetModal } from '@/pages/Assets/DeleteAssets';
import { symbolDict } from '@/utils/tokenIcon';

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

  const { nickname } = useParams();
  const { isMainnet } = useFTTransfer(nickname || '');

  return (
    <>
      {assets?.map((token: Asset, index: number) => (
        <Token
          logo={getAssetIcons(
            symbolDict[token.symbol as keyof typeof symbolDict],
            true,
            isMainnet,
            32,
          )}
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
