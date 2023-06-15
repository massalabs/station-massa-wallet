import { SyntheticEvent, useRef, useState } from 'react';
import { promptRequest } from '../../events/events';
import { useLocation, useNavigate } from 'react-router-dom';
import { handleCancel } from '../../utils/utils';
import { IErrorObject, parseForm } from '../../utils';
import {
  IsNicknameUnique,
  IsNicknameValid,
} from '../../../wailsjs/go/walletapp/WalletApp';

import { Input, Button, Stepper } from '@massalabs/react-ui-kit';
import { Layout } from '../../layouts/Layout/Layout';

import { IMPORT_STEPS } from '../../const/stepper';

// TODO: create i18n and move this to translation file
const t = (key: string): string => {
  const errors: Record<string, string> = {
    'nickname-required': 'The account name is required',
    'nickname-already-exists': 'This account name already exists',
    'account-invalid-format':
      "The account name can't contain any special characters",
  };

  return errors[key] || key;
};

function PromptNickname() {
  const navigate = useNavigate();
  const form = useRef(null);

  const [error, setError] = useState<IErrorObject | null>(null);

  let { state } = useLocation();
  const req: promptRequest = state.req;

  async function validate(e: SyntheticEvent) {
    const formObject = parseForm(e);
    const { nickname } = formObject;

    if (!nickname) {
      setError({ nickname: t('nickname-required') });
      return false;
    }

    if (await IsNicknameUnique(nickname)) {
      setError({ nickname: t('nickname-already-exists') });
      return false;
    }

    const nicknameIsValid = await IsNicknameValid(nickname);
    if (!nicknameIsValid) {
      setError({ nickname: t('account-invalid-format') });
      return false;
    }

    return true;
  }

  async function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    const formObject = parseForm(e);
    const { nickname } = formObject;

    if (!(await validate(e))) return;

    state = { ...state, nickname };
    navigate('/new-password', { state });
  }

  return (
    <Layout>
      <Stepper step={1} steps={IMPORT_STEPS} />
      <form ref={form} onSubmit={handleSubmit}>
        <h1 className="mas-title pt-4">{req.Msg}</h1>
        <p className="mas-body pt-4">Define your account name</p>
        <div className="pt-4">
          <Input
            defaultValue=""
            name="nickname"
            placeholder={'Account name'}
            error={error?.nickname}
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

export default PromptNickname;
