import { useEffect } from 'react';

import { Dropdown, formatAmount, IOption } from '@massalabs/react-ui-kit';
import { useParams } from 'react-router-dom';

import { MAS } from '@/const/assets/assets';
import { useResource } from '@/custom/api';
import Intl from '@/i18n/i18n';
import { Asset } from '@/models/AssetModel';
import { tokenIcon } from '@/utils/tokenIcon';

interface AssetSelectorProps {
  selectedAsset: Asset | undefined;
  setSelectedAsset: (asset: Asset) => void;
  selectSymbol?: string;
}

export function AssetSelector(props: AssetSelectorProps) {
  const { selectedAsset, setSelectedAsset, selectSymbol } = props;
  const { nickname } = useParams();

  const { data: assets, isLoading: isAssetsLoading } = useResource<Asset[]>(
    `accounts/${nickname}/assets`,
    false,
  );

  useEffect(() => {
    if (!assets || assets.length === 0) return;

    if (selectSymbol) {
      const requested = assets.find((asset) => asset.symbol === selectSymbol);
      if (requested) {
        setSelectedAsset(requested);
        return;
      }
    }

    if (!selectedAsset) {
      const defaultAsset = assets.find((a) => a.symbol === MAS) ?? assets[0];
      setSelectedAsset(defaultAsset);
    }
  }, [assets, setSelectedAsset, selectedAsset, selectSymbol]);

  // Defer rendering the Dropdown until assets are loaded: react-ui-kit's
  // Dropdown captures options[select] in its initial useState and only re-syncs
  // when `select` or `defaultItem` change, so mounting it with an empty
  // options array leaves the selection blank even after assets arrive.
  if (isAssetsLoading || !assets || assets.length === 0) {
    return <div className="pb-3.5 h-14" />;
  }

  const options: IOption[] = assets.map((asset) => {
    const formattedBalance = formatAmount(
      asset.balance || '',
      asset.decimals,
    ).full;
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
      icon: tokenIcon(asset.symbol, 28),
      onClick: () => setSelectedAsset(asset),
    };
  });

  const selectedAssetIndex = selectedAsset
    ? assets.findIndex(
        (a) =>
          a.address === selectedAsset.address &&
          a.symbol === selectedAsset.symbol,
      )
    : -1;
  const selectedAssetKey = selectedAssetIndex >= 0 ? selectedAssetIndex : 0;

  return (
    <Dropdown
      select={selectedAssetKey}
      size="md"
      options={options}
      className="pb-3.5"
      fullWidth={false}
    />
  );
}
