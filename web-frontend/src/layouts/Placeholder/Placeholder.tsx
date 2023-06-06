import { FiList } from 'react-icons/fi';
import Intl from '../../i18n/i18n';

function Placeholder() {
  return (
    <div
      className="flex flex-col justify-around items-center  
        h-[350px] w-[680px] p-10 gap-10
        relative left-20
        bg-secondary rounded-lg"
    >
      <FiList size={114} />
      <h1 className="mas-banner text-center">
        {Intl.t('placeholder.placeholder-banner')}
      </h1>
      <p className="mas-buttons text-center">
        {Intl.t('placeholder.placeholder-message')}
      </p>
    </div>
  );
}

export default Placeholder;
