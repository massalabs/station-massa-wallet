import { useState, useRef, SyntheticEvent, useEffect } from 'react';

import { Button, Identicon, Input } from '@massalabs/react-ui-kit';
import { AxiosError } from 'axios';
import { useNavigate, useParams } from 'react-router-dom';

import { SettingsOption } from './SettingsOption';
import { usePut } from '@/custom/api';
import Intl from '@/i18n/i18n';
import { WalletLayout, MenuItem } from '@/layouts/WalletLayout/WalletLayout';
import { AccountObject } from '@/models/AccountModel';
import { routeFor, parseForm, parseErrors } from '@/utils';
import { isNicknameValid } from '@/validation/nickname';

interface IErrorObject {
  code?: string;
  message?: string;
}

export default function SettingsUpdate() {
  const navigate = useNavigate();
  const { nickname } = useParams();
  const form = useRef(null);

  const [newNickname, setNewNickname] = useState<string>(nickname ?? '');
  const [error, setError] = useState<IErrorObject | null>(null);

  const {
    mutate: mutableUpdate,
    isSuccess: isSuccessUpdate,
    error: errorUpdate,
  } = usePut<AccountObject>(`accounts/${nickname}`);

  useEffect(() => {
    if (errorUpdate) {
      setError({ message: parseErrors(errorUpdate as AxiosError) });
    } else if (isSuccessUpdate) {
      navigate(routeFor(`${newNickname}/settings`));
    }
  }, [isSuccessUpdate, errorUpdate, newNickname, navigate]);

  function validate(e: SyntheticEvent) {
    const formObject = parseForm(e);
    const { nickname } = formObject;

    if (!nickname || !isNicknameValid(nickname)) {
      setError({ message: Intl.t('errors.account-invalid-format') });
      return false;
    }

    return true;
  }

  function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    const formObject = parseForm(e);
    const { nickname } = formObject;

    if (!validate(e)) return;
    mutableUpdate({ nickname } as AccountObject);
  }

  return (
    <WalletLayout menuItem={MenuItem.Settings}>
      <form ref={form} onSubmit={handleSubmit} className="w-1/2">
        <div className="flex flex-col justify-center items-center gap-9">
          <div className="bg-secondary rounded-2xl w-full max-w-2xl p-10">
            <p className="mas-body text-f-primary pb-5">
              {Intl.t('settings.title-profile')}
            </p>
            <div className="pb-7">
              <Identicon
                username={newNickname}
                size={120}
                className="bg-primary rounded-full"
              />
            </div>
            <div className="pb-5">
              <Input
                name="nickname"
                value={newNickname}
                placeholder={Intl.t('settings.inputs.nickname')}
                error={error?.message}
                onChange={(e) => setNewNickname(e.target.value)}
              />
            </div>
            <Button type="submit">{Intl.t('settings.buttons.update')}</Button>
          </div>
          <SettingsOption nickname={nickname ?? ''} />
        </div>
      </form>
    </WalletLayout>
  );
}
