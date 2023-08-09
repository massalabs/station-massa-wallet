type keyPairObject = {
  nonce: string;
  privateKey: string;
  publicKey: string;
  salt: string;
};

export interface IToken {
  name: string;
  assetAddress: string;
  symbol: string;
  decimals: number;
  balance: string;
}

export type AccountObject = {
  address: string;
  balance: string;
  candidateBalance: string;
  keyPair: keyPairObject;
  nickname: string;
  assets: IToken[];
};

export type SendTransactionObject = {
  amount: string;
  recipientAddress: string;
  fee: string;
};
