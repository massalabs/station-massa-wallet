import { Button } from '@massalabs/react-ui-kit';
import { useLocation, useNavigate } from 'react-router-dom';

import { promptRequest } from '@/events';
import { Layout } from '@/layouts/Layout/Layout';

const ImportMethods = () => {
  const navigate = useNavigate();

  const { state } = useLocation();
  const req: promptRequest = state.req;

  const handleYaml = () => navigate('/import-file', { state: { req } });
  const handlePkey = () => navigate('/import-key-pairs', { state: { req } });

  return (
    <Layout>
      <h1 className="mas-title">{req.Msg}</h1>
      <p className="mas-body pt-4">Choose an import method</p>
      <div className="pt-4">
        <Button variant="secondary" onClick={handleYaml}>
          I have a .yaml file
        </Button>
      </div>
      <div className="pt-4">
        <Button onClick={handlePkey}>I have a private key</Button>
      </div>
    </Layout>
  );
};

export default ImportMethods;
