import { SyntheticEvent, useRef, useState } from 'react';
import { promptRequest } from '../../events/events';
import { useLocation, useNavigate } from 'react-router-dom';
import { handleCancel } from '../../utils/utils';
import { Password, Button } from '@massalabs/react-ui-kit';
import { IErrorObject, parseForm } from '../../utils';

// TODO: create i18n and move this to translation file
const t = (key: string): string => {
  const errors: Record<string, string> = {
    'private-key-required': 'the Private key is required',
    'private-key-format': 'The private key doesn’t have the right format',
  };

  return errors[key] || key;
};

const PromptPrivateKey = () => {
  const navigate = useNavigate();
  const form = useRef(null);

  const [error, setError] = useState<IErrorObject | null>(null);

  let { state } = useLocation();
  const req: promptRequest = state.req;

  async function validate(e: SyntheticEvent) {
    const formObject = parseForm(e);
    const { privateKey } = formObject;

    if (!privateKey) {
      setError({ privateKey: t('private-key-required') });
      return false;
    }

    if (!privateKey.startsWith('S')) {
      setError({ privateKey: t('private-key-format') });
      return false;
    }

    return true;
  }

  async function handleSubmit(e: SyntheticEvent) {
    const formObject = parseForm(e);
    const { privateKey } = formObject;

    if (!(await validate(e))) return;

    state = { ...state, privateKey };
    navigate('/import-nickname', { state });
  }

  return (
    <div className="bg-primary flex flex-col justify-center items-center h-screen w-full">
      <div className="w-1/4  max-w-sm  min-w-fit">
        <form ref={form} onSubmit={handleSubmit}>
          <div>
            <p className="mas-title text-neutral pb-4">{req.Msg}</p>
          </div>
          <div>
            <p className="mas-body text-neutral pb-4">Enter your Private Key</p>
          </div>
          <div className="pb-4">
            <Password
              defaultValue=""
              name="privateKey"
              placeholder={'Private key'}
              error={error?.privateKey}
            />
          </div>
          <div className="flex flex-row gap-4 w-full pb-4">
            <div className="min-w-fit w-full">
              <Button variant={'secondary'} onClick={handleCancel}>
                Cancel
              </Button>
            </div>
            <div className="min-w-fit w-full">
              <Button type="submit">Next</Button>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
};

export default PromptPrivateKey;
