import { backupMethods, events, promptRequest } from '../events/events';
import { useLocation, useNavigate } from 'react-router-dom';
import { Layout } from '../layouts/Layout/Layout';
import { Button } from '@massalabs/react-ui-kit';
import { EventsOnce } from '../../wailsjs/runtime/runtime';
import { handleApplyResult } from '../utils';
import { SendPromptInput } from '../../wailsjs/go/walletapp/WalletApp';
import { useState } from 'react';

function BackupMethods() {
  const navigate = useNavigate();

  const [errorMsg, setErrorMsg] = useState('');

  const { state } = useLocation();
  const req: promptRequest = state.req;

  async function handleDownloadYaml() {
    setErrorMsg('');
    EventsOnce(
      events.promptResult,
      handleApplyResult(navigate, req, setErrorMsg, false),
    );

    await SendPromptInput(backupMethods.ymlFileBackup);
  }

  async function handleKeyPairs() {
    await SendPromptInput(backupMethods.privateKeyBackup);
    navigate('/backup-pkey', { state: { req } });
  }

  return (
    <Layout>
      <h1 className="mas-title">{req.Msg}</h1>
      <p className="mas-body pt-4">Choose a back up method</p>
      <div className="flex flex-col gap-4 pt-4">
        <Button variant="secondary" onClick={handleDownloadYaml}>
          Download a .yaml file
        </Button>
        {errorMsg && <p className="mas-body text-s-error">{errorMsg}</p>}
        <Button onClick={handleKeyPairs}>Show key pairs</Button>
      </div>
    </Layout>
  );
}

export default BackupMethods;
