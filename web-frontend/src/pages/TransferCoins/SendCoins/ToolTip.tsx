import { FiHelpCircle } from 'react-icons/fi';

function ToolTip({ ...props }) {
  const { showTooltip, content } = props;
  return (
    <div>
      <FiHelpCircle />
      {showTooltip && (
        <div className="flex flex-col w-96 absolute z-10 t-10 l-10 bg-tertiary p-3 rounded-lg text-neutral">
          {content}
        </div>
      )}
    </div>
  );
}

export default ToolTip;
