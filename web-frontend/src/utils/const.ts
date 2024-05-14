// TODO: add network store to be able to dynamically fetch correct contract

export enum MnsContract {
  mainnet = 'mainnet',
  buildnet = 'buildnet',
}

export const contracts = {
  [MnsContract.mainnet]: {
    MNScontract: '',
  },
  [MnsContract.buildnet]: {
    MNSContract: 'AS1CpitsdLu4dtbQrqAzhThygL2ytGyacFED1ogr2HsxZxfNy8qQ',
  },
};
