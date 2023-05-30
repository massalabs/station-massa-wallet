import { useLocation } from 'react-router';
import { promptAction, promptRequest } from '../events/events';
import { FiCheck } from 'react-icons/fi';
import I18n from '../i18n/i18n';

function Success() {
  const { state } = useLocation();
  const { req } = state;

  function successMsg(req: promptRequest) {
    switch (req.Action) {
      case promptAction.deleteReq:
        return I18n.t('success.delete');
      case promptAction.newPasswordReq:
        return I18n.t('success.new-password');
      case promptAction.importReq:
        return I18n.t('success.import');
      case promptAction.signReq:
        return I18n.t('success.sign');
      case promptAction.backupReq:
        return I18n.t('success.backup');
      case promptAction.transferReq:
        return I18n.t('success.transfer');
      default:
        return I18n.t('success.success');
    }
  }

  return (
    <div className="bg-primary flex flex-col justify-center items-center h-screen w-full">
      <div className="w-1/4 max-w-xs  min-w-fit flex flex-col justify-center items-center">
        <div className="w-12 h-12 bg-brand flex flex-col justify-center items-center rounded-full mb-6">
          <FiCheck className="w-6 h-6" />
        </div>
        <div>
          <p className="text-neutral mas-body">{successMsg(req)}</p>
        </div>
      </div>
    </div>
  );
}

export default Success;
