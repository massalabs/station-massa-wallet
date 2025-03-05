import { Button } from '@massalabs/react-ui-kit';
import { walletapp } from '@wailsjs/go/models';
import { SendPromptInput } from '@wailsjs/go/walletapp/WalletApp';
import { EventsOnce } from '@wailsjs/runtime/runtime';
import { FiAlertTriangle, FiTrash2 } from 'react-icons/fi';
import { useLocation, useNavigate } from 'react-router-dom';

import { promptRequest } from '@/events';
import { Layout } from '@/layouts/Layout/Layout';
import { handleApplyResult, handleCancel } from '@/utils';

function ConfirmDelete() {
  const navigate = useNavigate();
  const { state } = useLocation();
  const req: promptRequest = state.req;
  const password = state.password;

  const { EventType } = walletapp;

  function handleConfirm() {
    EventsOnce(EventType.promptResult, handleApplyResult(navigate, req));

    SendPromptInput(password);
  }

  return (
    <Layout>
      <h1 className="mas-title">{req.Msg}</h1>
      <div className="pt-4 flex gap-4">
        <FiAlertTriangle size={42} className="text-s-warning" strokeWidth="1" />
        <p className="mas-body">
          you have some coins left, are your
          <br />
          sure you want to delete?
        </p>
      </div>
      <div className="pt-4 flex gap-4">
        <Button variant="secondary" onClick={handleCancel}>
          Cancel
        </Button>
        <Button preIcon={<FiTrash2 />} onClick={handleConfirm}>
          Delete
        </Button>
      </div>
    </Layout>
  );
}

export default ConfirmDelete;
