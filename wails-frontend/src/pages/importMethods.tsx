import { h } from 'preact';
import { promptRequest } from '../events/events';
import { useLocation, useNavigate } from 'react-router-dom';

const ImportMethods = () => {
  const navigate = useNavigate();

  const { state } = useLocation();
  const req: promptRequest = state.req;

  const title = 'Import';
  const baselineStr = 'Choose an import method';

  const handleYml = () => navigate('/import-file', { state: { req } });
  const handlePkey = () => navigate('/import-methods', { state: { req } });

  return (
    <section class="PasswordPrompt">
      <div>{title}</div>
      <div className="baseline">{baselineStr}</div>
      <div className="flex flex-col">
        <button className="btn" onClick={handleYml}>
          I have a .yml file
        </button>
        <button className="btn" onClick={handlePkey}>
          I have a private key
        </button>
      </div>
    </section>
  );
};

export default ImportMethods;
