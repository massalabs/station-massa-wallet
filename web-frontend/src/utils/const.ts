export enum networks {
  mainnet = 'mainnet',
  buildnet = 'buildnet',
}

export const contracts = {
  [networks.mainnet]: {
    mnsContract: 'AS12mdKsjAqcWC5DjabWZp7tG9s5wgkrwDGDuh6xUCSc53SrmfY9r',
  },
  [networks.buildnet]: {
    mnsContract: 'AS12qKAVjU1nr66JSkQ6N4Lqu4iwuVc6rAbRTrxFoynPrPdP1sj3G',
  },
};

export const mnsExtension = '.massa';
