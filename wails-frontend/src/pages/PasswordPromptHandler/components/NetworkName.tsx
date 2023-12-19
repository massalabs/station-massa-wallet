import { MassaLogo } from "@massalabs/react-ui-kit";

export interface NetworkNameProps {
  networkName: string;
}

export function NetworkName({ networkName }: NetworkNameProps) {
  return (
    <div
      data-testid="tag"
      className="flex justify-between items-center bg-tertiary mas-caption
        rounded-full w-fit px-3 py-1 text-f-primary mb-4"
    >
      <MassaLogo size={16} />
      <span className="ml-2">{networkName}</span>
    </div>
  );
}
