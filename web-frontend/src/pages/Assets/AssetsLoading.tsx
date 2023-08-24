import { MassaLogo, Token } from '@massalabs/react-ui-kit';

export function AssetsLoading() {
  return (
    <>
      <Token
        customClass="animate-pulse blur-md"
        logo={<MassaLogo size={40} />}
        name={'lorem'}
        symbol={'IPS'}
        decimals={7}
        balance={'000000000'}
      />
    </>
  );
}
