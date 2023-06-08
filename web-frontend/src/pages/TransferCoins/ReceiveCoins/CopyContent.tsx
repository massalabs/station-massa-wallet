import { FiCopy } from 'react-icons/fi';
interface CopyAddressProps {
  content: string;
  formattedContent: string;
}
function CopyContent(props: CopyAddressProps) {
  const { content, formattedContent } = props;

  function handleCopyClick() {
    navigator.clipboard.writeText(content);
  }

  return (
    <>
      <div
        className="flex flex-row items-center mas-body2 justify-between w-full 
        h-12 px-3 rounded bg-secondary cursor-pointer"
        onClick={() => handleCopyClick()}
      >
        <u>{formattedContent ?? content}</u>
        <FiCopy size={24} />
      </div>
    </>
  );
}

export default CopyContent;
