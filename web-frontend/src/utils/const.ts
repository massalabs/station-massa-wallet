export enum networks {
  mainnet = 'mainnet',
  buildnet = 'buildnet',
}

export const contracts = {
  [networks.mainnet]: {
    mnsContract: '',
  },
  [networks.buildnet]: {
    mnsContract: 'AS1CpitsdLu4dtbQrqAzhThygL2ytGyacFED1ogr2HsxZxfNy8qQ',
  },
};

export const mnsExtension = '.massa';
