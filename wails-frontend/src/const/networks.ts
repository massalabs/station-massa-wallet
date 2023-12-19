export const MAINNET = 'MainNet';
export const BUILDNET = 'BuildNet';
export const SECURENET = 'SecureNet';
export const LABNET = 'LabNet';
export const SANDBOX = 'Sandbox';

export const CHAIN_ID_TO_NETWORK_NAME = {
  77658377: MAINNET,
  77658366: BUILDNET,
  77658383: SECURENET,
  77658376: LABNET,
  77: SANDBOX,
} as const; // type is inferred as the specific, unchangeable structure

export type ChainId = keyof typeof CHAIN_ID_TO_NETWORK_NAME;
