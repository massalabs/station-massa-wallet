import { FT1, MassaLogo, Token } from '@massalabs/react-ui-kit';

import { ITokenData, XMA } from '@/const/assets/assets';

export function AssetsList({ ...props }) {
  const { tokenArray } = props;

  return (
    <>
      {tokenArray?.map((token: ITokenData, index: number) => (
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
        />
      ))}
    </>
  );
}
