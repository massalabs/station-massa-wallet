import { Button } from '@massalabs/react-ui-kit';
import { Layout } from '../layouts/Layout/Layout';
import { useLocation, useNavigate } from 'react-router-dom';
import { events, promptRequest } from '../events/events';
import { FiAlertTriangle, FiTrash2 } from 'react-icons/fi';
import { handleApplyResult, handleCancel } from '../utils/utils';
import { EventsOnce } from '../../wailsjs/runtime/runtime';
import { SendPromptInput } from '../../wailsjs/go/walletapp/WalletApp';

function ConfirmDelete() {
  const navigate = useNavigate();
  const { state } = useLocation();
  const req: promptRequest = state.req;
  const password = state.password;

  function handleConfirm() {
    EventsOnce(events.promptResult, handleApplyResult(navigate, req));

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
        <Button variant={'secondary'} onClick={handleCancel}>
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
