import { h } from 'preact';
import { useLocation } from 'react-router';
import { promptAction, promptRequest } from '../events/events';

const Success = () => {
  const { state } = useLocation();
  const { req } = state;

  const successMsg = (req: promptRequest) => {
    switch (req.Action) {
      case promptAction.deleteReq:
        return 'The account has been deleted';
      case promptAction.newPasswordReq:
        return 'The password has been created';
      case promptAction.importReq:
        return 'The account has been imported';
      case promptAction.signReq:
        return 'The transaction has been signed';
      default:
        return 'Apply';
    }
  };

  return (
    <section class="Success">
      <div>Success !!</div>
      <div>{successMsg(req)}</div>
    </section>
  );
};

export default Success;
