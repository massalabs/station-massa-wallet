export enum networks {
  mainnet = 'mainnet',
  buildnet = 'buildnet',
}

export const contracts = {
  [networks.mainnet]: {
    mnsContract: 'AS1q5hUfxLXNXLKsYQVXZLK7MPUZcWaNZZsK7e9QzqhGdAgLpUGT',
  },
  [networks.buildnet]: {
    mnsContract: 'AS12qKAVjU1nr66JSkQ6N4Lqu4iwuVc6rAbRTrxFoynPrPdP1sj3G',
  },
};

export const mnsExtension = '.massa';
