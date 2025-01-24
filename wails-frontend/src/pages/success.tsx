import { walletapp } from '@wailsjs/go/models';
import { FiCheck } from 'react-icons/fi';
import { useLocation } from 'react-router';

import { promptRequest } from '@/events';
import Intl from '@/i18n/i18n';
import { Layout } from '@/layouts/Layout/Layout';

function Success() {
  const { state } = useLocation();
  const { req } = state;

  const { PromptRequestAction } = walletapp;

  function successMsg(req: promptRequest) {
    switch (req.Action) {
      case PromptRequestAction.delete:
        return Intl.t('success.delete');
      case PromptRequestAction.newPassword:
        return Intl.t('success.new-password');
      case PromptRequestAction.import:
        return Intl.t('success.import');
      case PromptRequestAction.sign:
        return Intl.t('success.sign');
      case PromptRequestAction.backup:
        return Intl.t('success.backup');
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
