import { useState, useRef, SyntheticEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import Intl from '../../i18n/i18n';
import { routeFor } from '../../utils';
import { parseForm } from '../../utils/parseForm';

import { Input, Stepper, Button } from '@massalabs/react-ui-kit';
import { FiArrowRight } from 'react-icons/fi';
import LandingPage from '../../layouts/LandingPage/LandingPage';
import { useResource } from '../../custom/api';
import { AccountObject } from '../../models/AccountModel';
import { isAlreadyExists, isNicknameValid } from '../../validation/nickname';

interface IErrorObject {
  nickname: string;
}

export default function StepOne() {
  const navigate = useNavigate();

  const { data: accounts = [] } = useResource<AccountObject[]>('accounts');

  const form = useRef(null);

  const [error, setError] = useState<IErrorObject | null>(null);

  function validate(e: SyntheticEvent) {
    const formObject = parseForm(e);
    const { nickname } = formObject;

    if (!nickname) {
      setError({ nickname: Intl.t('errors.nickname-required') });
      return false;
    }

    if (isAlreadyExists(nickname, accounts)) {
      setError({ nickname: Intl.t('errors.nickname-already-exists') });
      return false;
    }

    if (!isNicknameValid(nickname)) {
      setError({ nickname: Intl.t('errors.nickname-invalid-format') });
      return false;
    }

    return true;
  }

  function sendNickname(e: SyntheticEvent) {
    const formObject = parseForm(e);
    const { nickname } = formObject;

    navigate(routeFor('account-create-step-two'), { state: { nickname } });
  }

  function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();

    if (!validate(e)) return;
    sendNickname(e);
  }

  return (
    <LandingPage>
      <form ref={form} onSubmit={handleSubmit}>
        <div className="flex flex-col justify-center items-center h-screen">
          <div className="flex flex-col justify-start w-full h-full max-w-sm max-h-56">
            <div className="w-full pb-6">
              <Stepper
                step={0}
                steps={[
                  Intl.t('account.create.step1.title'),
                  Intl.t('account.create.step2.title'),
                  Intl.t('account.create.step3.title'),
                ]}
              />
            </div>
            <h1 className="mas-banner text-neutral pb-6">
              {Intl.t('account.create.title')}
            </h1>
            <div className="w-full pb-4">
              <Input
                defaultValue=""
                name="nickname"
                placeholder={Intl.t('account.create.input.nickname')}
                error={error?.nickname}
              />
            </div>
            <div className="w-full">
              <Button posIcon={<FiArrowRight />} type="submit">
                {Intl.t('account.create.buttons.next')}
              </Button>
            </div>
          </div>
        </div>
      </form>
    </LandingPage>
  );
}
