import { SyntheticEvent, useRef, useState } from 'react';
import { promptRequest } from '../events/events';
import { useLocation, useNavigate } from 'react-router-dom';
import { handleCancel } from '../utils/utils';
import { Password, Input, Button } from '@massalabs/react-ui-kit';
import { IErrorObject, parseForm } from '../utils';
import {
  IsNicknameUnique,
  IsNicknameValid,
} from '../../wailsjs/go/walletapp/WalletApp';

// TODO: create i18n and move this to translation file
const t = (key: string): string => {
  const errors: Record<string, string> = {
    'nickname-required': 'Nickname is required',
    'nickname-already-exists': 'This account name already exists',
    'nickname-invalid-format':
      "The account name can't contain any special characters",
  };

  return errors[key] || key;
};

const ImportPrivatekey = () => {
  const navigate = useNavigate();
  const form = useRef(null);

  const [input, setInput] = useState<string | undefined>(undefined);
  const [privateKey, setPkey] = useState<string | undefined>(undefined);
  const [error, setError] = useState<IErrorObject | null>(null);

  let { state } = useLocation();
  const req: promptRequest = state.req;

  const applyStr = 'Next';

  const baselineStr = () =>
    privateKey ? 'Define your account name' : 'Enter your Private Key';

  const placeholder = () => (privateKey ? 'Account name' : 'Private key');

  async function validateNickname(e: SyntheticEvent) {
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
      setError({ nickname: t('nickname-invalid-format') });
      return false;
    }

    return true;
  }

  async function handleSubmit(e: SyntheticEvent) {
    if (!privateKey) {
      // if the private key is not defined, set the private key
      setPkey(input);
      // and reset the input
      setInput('');
    } else {
      const nicknameIsValid = await validateNickname(e);

      // if the private key is defined, validate the nickname
      if (!nicknameIsValid) return;
      // if the nickname is valid, navigate to the password page
      state = { ...state, pkey: privateKey, nickname: input };
      navigate('/new-password', { state });
    }
  }

  function updateInput(e: React.ChangeEvent<HTMLInputElement>) {
    setInput(e.target.value);
  }

  return (
    <div className="bg-primary flex flex-col justify-center items-center h-screen w-full">
      <div className="w-1/4  max-w-sm  min-w-fit">
        <form ref={form} onSubmit={handleSubmit}>
          <div>
            <p className="mas-title text-neutral pb-4">{req.Msg}</p>
          </div>
          <div>
            <p className="mas-body text-neutral pb-4">{baselineStr()}</p>
          </div>
          {privateKey ? (
            <div className="pb-4">
              <Input
                onChange={updateInput}
                autoComplete="off"
                name="nickname"
                placeholder={placeholder()}
                value={input}
                error={error?.nickname}
              />
            </div>
          ) : (
            <div className="pb-4">
              <Password
                onChange={updateInput}
                autoComplete="off"
                name="password"
                placeholder={placeholder()}
                value={input}
              />
            </div>
          )}
          <div className="flex flex-row gap-4 w-full pb-4">
            <div className="min-w-fit w-full">
              <Button variant={'secondary'} onClick={handleCancel}>
                Cancel
              </Button>
            </div>
            <div className="min-w-fit w-full">
              <Button type="submit">{applyStr}</Button>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
};

export default ImportPrivatekey;
