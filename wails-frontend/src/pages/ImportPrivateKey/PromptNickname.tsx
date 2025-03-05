import { SyntheticEvent, useRef, useState } from 'react';

import { Input, Button, Stepper } from '@massalabs/react-ui-kit';
import {
  IsNicknameUnique,
  IsNicknameValid,
} from '@wailsjs/go/walletapp/WalletApp';
import { useLocation, useNavigate } from 'react-router-dom';

import { IMPORT_STEPS } from '@/const/stepper';
import { promptRequest } from '@/events';
import Intl from '@/i18n/i18n';
import { Layout } from '@/layouts/Layout/Layout';
import { IErrorObject, parseForm, handleCancel } from '@/utils';

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
      setError({ nickname: Intl.t('errors.nickname-required') });
      return false;
    }

    if (await IsNicknameUnique(nickname)) {
      setError({ nickname: Intl.t('errors.DuplicateNickname-001') });
      return false;
    }

    const nicknameIsValid = await IsNicknameValid(nickname);
    if (!nicknameIsValid) {
      setError({ nickname: Intl.t('errors.account-invalid-format') });
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
            placeholder={Intl.t('import.steps.account')}
            error={error?.nickname}
          />
        </div>
        <div className="flex gap-4 pt-4">
          <Button variant="secondary" onClick={handleCancel}>
            Cancel
          </Button>
          <Button type="submit">Next</Button>
        </div>
      </form>
    </Layout>
  );
}

export default PromptNickname;
