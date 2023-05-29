import { promptRequest } from '../events/events';
import { useLocation, useNavigate } from 'react-router-dom';
import { Button } from '@massalabs/react-ui-kit';

const ImportMethods = () => {
  const navigate = useNavigate();

  const { state } = useLocation();
  const req: promptRequest = state.req;

  const baselineStr = 'Choose an import method';

  const handleYaml = () => navigate('/import-file', { state: { req } });
  const handlePkey = () => navigate('/import-key-pairs', { state: { req } });

  return (
    <div className="bg-primary flex flex-col justify-center items-center h-screen w-full">
      <div className="w-1/4  max-w-xs  min-w-fit">
        <div>
          <div>
            <p className="mas-title text-neutral pb-4">{req.Msg}</p>
          </div>
          <div>
            <p className="mas-body text-neutral pb-4">{baselineStr}</p>
          </div>
        </div>
        <div className="flex flex-col">
          <div className="pb-4 w-full">
            <Button variant={'secondary'} onClick={handleYaml}>
              I have a .yaml file
            </Button>
          </div>
          <div className="w-full">
            <Button onClick={handlePkey}>I have a private key</Button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ImportMethods;
