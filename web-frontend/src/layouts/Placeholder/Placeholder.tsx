import Intl from '../../i18n/i18n';

function Placeholder({ ...props }) {
  const { message, icon } = props;
  return (
    <div
      className="flex flex-col justify-around items-center  
        h-96 w-[700px] p-10 gap-10
        relative left-20
        bg-secondary rounded-lg"
    >
      {icon}
      <h1 className="mas-banner text-center">
        {Intl.t('placeholder.placeholder-banner')}
      </h1>
      <p className="mas-buttons text-center">{message}</p>
    </div>
  );
}

export default Placeholder;
