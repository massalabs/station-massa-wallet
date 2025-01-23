import { ClipboardSetText } from '@wailsjs/runtime';
import { FiCopy } from 'react-icons/fi';

interface CopyProps {
  data: string;
}

export function CopyClip({ data }: CopyProps) {
  async function handleCopy() {
    await ClipboardSetText(data);
  }

  return <FiCopy onClick={handleCopy} style={{ cursor: 'pointer' }} />;
}
