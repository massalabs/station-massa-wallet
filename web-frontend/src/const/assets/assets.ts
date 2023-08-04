export const XMA = 'XMA';

export interface ITokenData {
  logo?: React.ReactNode;
  name: string;
  symbol: string;
  decimals: number;
  balance: string;
  customClass?: string;
  disable?: boolean;
  onDelete?: () => void;
}

export const assetImportErrors = {
  success: 201,
  badRequest: 400,
  invalidAddress: 422,
  notFound: 404,
  serverError: 500,
};
