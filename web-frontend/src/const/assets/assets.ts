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
