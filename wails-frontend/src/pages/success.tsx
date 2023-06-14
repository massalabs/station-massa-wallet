import { useLocation } from 'react-router';
import { promptAction, promptRequest } from '../events/events';
import { FiCheck } from 'react-icons/fi';
import Intl from '../i18n/i18n';
import { Layout } from '../layouts/Layout/Layout';

function Success() {
  const { state } = useLocation();
  const { req } = state;

  function successMsg(req: promptRequest) {
    switch (req.Action) {
      case promptAction.deleteReq:
        return Intl.t('success.delete');
      case promptAction.newPasswordReq:
        return Intl.t('success.new-password');
      case promptAction.importReq:
        return Intl.t('success.import');
      case promptAction.signReq:
        return Intl.t('success.sign');
      case promptAction.backupReq:
        return Intl.t('success.backup');
      case promptAction.transferReq:
        return Intl.t('success.transfer');
      default:
        return Intl.t('success.success');
    }
  }

  return (
    <Layout>
      <div className="flex flex-col items-center justify-center">
        <div className="w-12 h-12 bg-brand flex flex-col justify-center items-center rounded-full mb-6">
          <FiCheck className="w-6 h-6" />
        </div>
        <p className="text-neutral mas-body">{successMsg(req)}</p>
      </div>
    </Layout>
  );
}

export default Success;
