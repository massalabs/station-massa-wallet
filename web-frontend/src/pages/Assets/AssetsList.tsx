import { ITokenData, XMA } from '@/const/assets/assets';
import { MassaLogo, Mrc20, Token } from '@massalabs/react-ui-kit';

export function AssetsList({ ...props }) {
  const { tokenArray } = props;
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
        />
      ))}
    </>
  );
}
