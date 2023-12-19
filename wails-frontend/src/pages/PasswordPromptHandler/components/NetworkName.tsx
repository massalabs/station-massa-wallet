import { MassaLogo } from "@massalabs/react-ui-kit";

import { MAINNET } from "@/const/networks";

export interface NetworkNameProps {
  networkName: string;
}

export function NetworkName({ networkName }: NetworkNameProps) {
  let secondaryColor = undefined;
  let primaryColor = undefined;

  if (networkName === MAINNET) {
    primaryColor = "#FF0000";
    secondaryColor = "#FFFFFF";
  } else {
    // colors from the design on figma labelled as Buildnet
    primaryColor = "#FFFFFF";
    secondaryColor = "#151A26";
  }

  return (
    <div
      data-testid="tag"
      className="flex justify-between items-center bg-tertiary mas-caption
        rounded-full w-fit px-3 py-1 text-f-primary mb-4"
    >
      <MassaLogo size={16} primaryColor={primaryColor} secondaryColor={secondaryColor} />
      <span className="ml-2">{networkName}</span>
    </div>
  );
}
