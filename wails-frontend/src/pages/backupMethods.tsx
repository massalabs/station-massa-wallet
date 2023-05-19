import { promptRequest } from '../events/events';
import { useLocation, useNavigate } from 'react-router-dom';

import { Layout } from '../layouts/Layout/Layout';
import { Button } from '@massalabs/react-ui-kit';

function BackupMethods() {
  const navigate = useNavigate();

  const { state } = useLocation();
  const req: promptRequest = state.req;

  function handleDownloadYml() {
    // TODO: download a .yaml file
    console.log('Download a .yaml file');
  }

  function handleKeyPairs() {
    navigate('/backup-pkey', { state: { req } });
  }

  return (
    <Layout>
      <h1 className="mas-title">{req.Msg}</h1>
      <p className="mas-body pt-4">Choose a back up method</p>
      <div className="flex flex-col gap-4 pt-4">
        <Button variant="secondary" onClick={handleDownloadYml}>
          Download a .yaml file
        </Button>
        <Button onClick={handleKeyPairs}>Show key pairs</Button>
      </div>
    </Layout>
  );
}

export default BackupMethods;
