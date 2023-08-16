import { Token } from './AssetModel';

interface keyPairObject {
  nonce: string;
  privateKey: string;
  publicKey: string;
  salt: string;
}

export interface AccountObject {
  address: string;
  balance: string;
  candidateBalance: string;
  keyPair: keyPairObject;
  nickname: string;
  assets: Token[];
  status: string;
}

export type SendTransactionObject = {
  amount: string;
  recipientAddress: string;
  fee: string;
};
