import { useEffect, useState } from 'react';

import { MassaLogo, Mrc20, Token, toast } from '@massalabs/react-ui-kit';
import { useParams } from 'react-router-dom';

import { ITokenData, XMA } from '@/const/assets/assets';
import { useDelete } from '@/custom/api';
import { IToken } from '@/models/AccountModel';

export function AssetsList({ ...props }) {
  const { tokenArray } = props;
  const { nickname } = useParams();

  const [tokenAddress, setTokenAddress] = useState<string>('');

  const {
    mutate: mutateDelete,
    isSuccess: isSuccessDelete,
    isError: isErrorDelete,
    error: errorDelete,
  } = useDelete<IToken[]>(
    `accounts/${nickname}/assets?assetAddress=${tokenAddress}`,
  );

  useEffect(() => {
    if (isSuccessDelete) {
      toast.success('Token Deleted Successfully');
    } else if (isErrorDelete) {
      console.log(errorDelete);
    }
  }, [isSuccessDelete]);

  function handleDelete(address: string) {
    setTokenAddress(address);
    confirmDelete();
  }

  function confirmDelete() {
    mutateDelete({} as IToken[]);
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
    </>
  );
}
