import { useState } from 'react';

import { Button } from '@massalabs/react-ui-kit';
import { maskNickname } from '@massalabs/react-ui-kit/src/lib/massa-react/utils';
import { SendPromptInput } from '@wailsjs/go/walletapp/WalletApp';
import { EventsOnce } from '@wailsjs/runtime/runtime';
import { useLocation, useNavigate } from 'react-router-dom';

import { backupMethods, events, promptRequest } from '@/events/events';
import Intl from '@/i18n/i18n';
import { Layout } from '@/layouts/Layout/Layout';
import { handleApplyResult } from '@/utils';

function BackupMethods() {
  const { state } = useLocation();
  const navigate = useNavigate();

  const [errorMsg, setErrorMsg] = useState('');

  const req: promptRequest = state.req;
  const walletName: string = req.Msg;

  async function handleDownloadYaml() {
    EventsOnce(
      events.promptResult,
      handleApplyResult(navigate, req, setErrorMsg, true),
    );

    await SendPromptInput(backupMethods.ymlFileBackup, req.CorrelationID);
  }

  async function handleKeyPairs() {
    await SendPromptInput(backupMethods.privateKeyBackup, req.CorrelationID);

    navigate('/backup-pkey', { state: { req } });
  }

  return (
    <Layout>
      <div className="flex items-end gap-3">
        <h1 className="mas-title">{Intl.t('backup.title')}</h1>
        <span>/</span>
        <h2 className="mas-h2">{maskNickname(walletName)}</h2>
      </div>
      <p className="mas-body pt-4">{Intl.t('backup.subtitle.choose-method')}</p>
      <div className="flex flex-col gap-4 pt-4">
        <Button variant="secondary" onClick={handleDownloadYaml}>
          {Intl.t('backup.buttons.download')}
        </Button>
        {errorMsg && <p className="mas-body text-s-error">{errorMsg}</p>}
        <Button onClick={handleKeyPairs}>
          {Intl.t('backup.buttons.show')}
        </Button>
      </div>
    </Layout>
  );
}

export default BackupMethods;
