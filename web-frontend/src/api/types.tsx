export type accountType = {
  nickname: string;
  address: string;
  candidateBalance: string;
  balance: string;
  keypair: {
    privateKey: string;
    publicKey: string;
    salt: string;
    nonce: string;
  };
};
