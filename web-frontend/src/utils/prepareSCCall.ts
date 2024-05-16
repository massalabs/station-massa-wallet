import { ClientFactory } from '@massalabs/massa-web3';
import { providers } from '@massalabs/wallet-provider';

export async function prepareSCCall(nickname: string) {
  const massaProvider = (await providers())?.find(
    (p) => p.name() === 'MASSASTATION',
  );
  if (!massaProvider) {
    console.error('FATAL: Massa provider not found');
    return;
  }

  const account = (await massaProvider.accounts()).find(
    (a) => a.name() === nickname,
  );
  if (!account) {
    return;
  }

  return {
    client: await ClientFactory.fromWalletProvider(massaProvider, account),
  };
}
