export const XMA = 'XMA';

export interface ITokenData {
  logo?: React.ReactNode;
  name: string;
  symbol: string;
  assetAddress: string;
  decimals: number;
  balance: string;
  customClass?: string;
  disable?: boolean;
  onDelete?: () => void;
}

export interface IDeleteAssetsErrors {
  phrase?: string;
}

export const assetDeleteErrors = {
  success: 204,
  badRequest: 400,
  invalidAddress: 404,
  serverError: 500,
};

export const deleteConfirm = 'DELETE';
