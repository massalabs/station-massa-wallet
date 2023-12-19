export const CHAIN_ID_TO_NETWORK_NAME = {
  77658377: 'MainNet',
  77658366: 'BuildNet',
  77658383: 'SecureNet',
  77658376: 'LabNet',
  77: 'Sandbox',
} as const; // type is inferred as the specific, unchangeable structure

export type ChainId = keyof typeof CHAIN_ID_TO_NETWORK_NAME;
