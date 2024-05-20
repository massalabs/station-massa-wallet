import { useEffect } from 'react';

import {
  Dropdown,
  formatAmount,
  getAssetIcons,
  IOption,
} from '@massalabs/react-ui-kit';
import { useParams } from 'react-router-dom';

import { useResource } from '@/custom/api';
import { useFTTransfer } from '@/custom/smart-contract/useFTTransfer';
import Intl from '@/i18n/i18n';
import { Asset } from '@/models/AssetModel';
import { symbolDict } from '@/utils/tokenIcon';

interface AssetSelectorProps {
  selectedAsset: Asset | undefined;
  setSelectedAsset: (asset: Asset) => void;
  selectSymbol?: string;
}

export function AssetSelector(props: AssetSelectorProps) {
  const { selectedAsset, setSelectedAsset, selectSymbol } = props;
  const { nickname } = useParams();
  const { isMainnet } = useFTTransfer();

  const { data: assets, isLoading: isAssetsLoading } = useResource<Asset[]>(
    `accounts/${nickname}/assets`,
    false,
  );

  useEffect(() => {
    if (selectSymbol) {
      const selectedAsset = assets?.find(
        (asset) => asset.symbol === selectSymbol,
      );
      if (selectedAsset) setSelectedAsset(selectedAsset);
    }
    if (!selectedAsset && assets && assets?.length > 0) {
      setSelectedAsset(assets?.[0]);
    }
  }, [assets, setSelectedAsset, selectedAsset, selectSymbol]);

  let options: IOption[] = [];

  if (assets) {
    options = assets.map((asset) => {
      const formattedBalance = formatAmount(
        asset.balance || '',
        asset.decimals,
      ).amountFormattedFull;
      return {
        itemPreview: asset.symbol,
        item: (
          <div>
            <p>{asset.symbol}</p>
            <p className="mas-caption">
              {Intl.t('send-coins.balance')} {formattedBalance}
            </p>
          </div>
        ),
        icon: getAssetIcons(
          symbolDict[asset.symbol as keyof typeof symbolDict],
          true,
          isMainnet,
          28,
        ),
        onClick: () => setSelectedAsset(asset),
      };
    });
  }

  const selectedAssetKey: number = selectedAsset
    ? assets?.indexOf(selectedAsset) || 0
    : 0;

  return (
    <Dropdown
      select={selectedAssetKey}
      readOnly={isAssetsLoading}
      size="md"
      options={options}
      className="pb-3.5"
      fullWidth={false}
    />
  );
}
