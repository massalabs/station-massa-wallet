import { useLocation } from 'react-router';
import { FiX } from 'react-icons/fi';
import { promptRequest } from '../events/events';

const Failure = () => {
  const locationState = useLocation().state as { req: promptRequest };
  const req = locationState.req;

  return (
    <div className="bg-primary flex flex-col justify-center items-center h-screen w-full">
      <div className="w-1/4  max-w-xs  min-w-fit flex flex-col justify-center items-center">
        <div className="w-12 h-12 bg-s-error flex flex-col justify-center items-center rounded-full mb-6">
          <FiX className="w-6 h-6" />
        </div>
        <div>
          <p className="text-neutral mas-body">{req.Msg}</p>
        </div>
      </div>
    </div>
  );
};

export default Failure;
