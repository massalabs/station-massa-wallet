import { getAssetIcons } from '@massalabs/react-ui-kit';

export function tokenIcon(symbol: string, size: number) {
  return getAssetIcons(symbol, undefined, false, size) as JSX.Element;
}
