import { SyntheticEvent, useRef, useState } from 'react';
import { promptRequest } from '../../events/events';
import { useLocation, useNavigate } from 'react-router-dom';
import { handleCancel } from '../../utils/utils';
import { IErrorObject, parseForm } from '../../utils';

import { Password, Button, Stepper } from '@massalabs/react-ui-kit';
import { Layout } from '../../layouts/Layout/Layout';

// TODO: create i18n and move this to translation file
const t = (key: string): string => {
  const errors: Record<string, string> = {
    'private-key-required': 'the Private key is required',
    'private-key-format': 'The private key doesn’t have the right format',
  };

  return errors[key] || key;
};

function PromptPrivateKey() {
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
    <Layout>
      <Stepper step={0} steps={['Private Key', 'Account name', 'Password']} />
      <form ref={form} onSubmit={handleSubmit}>
        <h1 className="mas-title pt-6">{req.Msg}</h1>
        <p className="mas-body pt-4">Enter your Private Key</p>
        <div className="pt-4">
          <Password
            defaultValue=""
            name="privateKey"
            placeholder={'Private key'}
            error={error?.privateKey}
          />
        </div>
        <div className="flex gap-4 pt-4">
          <Button variant={'secondary'} onClick={handleCancel}>
            Cancel
          </Button>
          <Button type="submit">Next</Button>
        </div>
      </form>
    </Layout>
  );
}

export default PromptPrivateKey;
