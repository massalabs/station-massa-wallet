import { useLocation } from 'react-router';
import { promptAction, promptRequest } from '../events/events';
import { FiX } from 'react-icons/fi';

const Failure = () => {
  const { state } = useLocation();
  const { req } = state;

  const successMsg = (req: promptRequest) => {
    switch (req.Action) {
      case promptAction.deleteReq:
        return 'The account could not been deleted';
      case promptAction.newPasswordReq:
        return 'The password could not be created';
      case promptAction.importReq:
        return 'The account could not be imported';
      case promptAction.signReq:
        return 'The transaction has been signed';
      default:
        return 'Failed';
    }
  };

  return (
    <div className="bg-primary flex flex-col justify-center items-center h-screen w-full">
      <div className="w-1/4  max-w-xs  min-w-fit flex flex-col justify-center items-center">
        <div className="w-12 h-12 bg-s-error flex flex-col justify-center items-center rounded-full mb-6">
          <FiX className="w-6 h-6" />
        </div>
        <div>
          <p className="text-neutral mas-body">{successMsg(req)}</p>
        </div>
      </div>
    </div>
  );
};

export default Failure;
