import { useState, useRef, SyntheticEvent } from 'react';

import { Password, Button, Clipboard } from '@massalabs/react-ui-kit';
import { walletapp } from '@wailsjs/go/models';
import { SendPromptInput } from '@wailsjs/go/walletapp/WalletApp';
import { ClipboardSetText, EventsOnce } from '@wailsjs/runtime';
import { FiCopy, FiArrowRight } from 'react-icons/fi';
import { useLocation, useNavigate } from 'react-router-dom';

import { promptRequest, promptResult } from '@/events';
import Intl from '@/i18n/i18n';
import { Layout } from '@/layouts/Layout/Layout';
import { ErrorCode, IErrorObject, handleCancel } from '@/utils';
import { parseForm } from '@/utils/parseForm';

function EnterKey() {
  return (
    <>
      <Button variant="secondary" onClick={handleCancel}>
        {Intl.t('backup.buttons.cancel')}
      </Button>
      <Button posIcon={<FiArrowRight />} type="submit">
        {Intl.t('backup.buttons.next')}
      </Button>
    </>
  );
}

interface KeyPair {
  privateKey: string;
  publicKey: string;
}

interface CopyProps {
  privateKey: string;
}

function CopyKey({ privateKey }: CopyProps) {
  async function handleCopy() {
    await ClipboardSetText(privateKey);
  }

  return (
    <Button onClick={handleCopy} preIcon={<FiCopy />}>
      {Intl.t('backup.buttons.copy')}
    </Button>
  );
}

function BackupKeyPairs() {
  const form = useRef(null);
  const { state } = useLocation();
  const navigate = useNavigate();

  const req: promptRequest = state.req;
  const walletName: string = req.Msg;

  const [privateKey, setPrivateKey] = useState<string | undefined>(undefined);
  const [publicKey, setPublicKey] = useState<string | undefined>(undefined);

  const [error, setError] = useState<IErrorObject | null>(null);

  const { EventType } = walletapp;

  function validate(e: SyntheticEvent) {
    const formObject = parseForm(e);
    const { password } = formObject;

    if (!password || !password.length) {
      setError({ password: Intl.t('errors.PasswordRequired') });
      return false;
    }

    return true;
  }

  function handleResult(result: promptResult<KeyPair>) {
    let { Success, CodeMessage } = result;

    if (!Success) {
      const errorMessage = Intl.t(`errors.${CodeMessage}`);
      if (CodeMessage === ErrorCode.WrongPassword) {
        setError({ password: errorMessage });
        return;
      } else {
        navigate('/failure', {
          state: { req: { Msg: errorMessage } },
        });
      }
    }
    setPrivateKey(result.Data.privateKey);
    setPublicKey(result.Data.publicKey);
  }

  async function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    if (!validate(e)) return;

    const formObject = parseForm(e);
    const { password } = formObject;

    EventsOnce(EventType.promptResult, handleResult);
    await SendPromptInput(password);
  }

  return (
    <Layout>
      <form ref={form} onSubmit={handleSubmit}>
        <div className="flex items-end gap-3 mb-4">
          <h1 className="mas-title">{Intl.t('backup.title')}</h1>
          <span>/</span>
          <h2 className="mas-h2">{walletName}</h2>
        </div>
        {!privateKey && (
          <p className="mas-body pt-4">
            {Intl.t('backup.subtitle.private-key')}
          </p>
        )}
        <div className="pt-4">
          {privateKey ? (
            <>
              <div className="mb-3">
                <p className="mb-2 mas-body">
                  {Intl.t('backup.subtitle.public-key-title')}
                </p>
                <Clipboard
                  displayedContent={publicKey}
                  rawContent={publicKey as string}
                />
              </div>
              <div>
                <p className="mb-2 mas-body">
                  {Intl.t('backup.subtitle.private-key-title')}
                </p>
                <Password
                  defaultValue=""
                  name="privateKey"
                  placeholder={Intl.t('backup.inputs.private-key.placeholder')}
                  value={privateKey}
                />
              </div>
            </>
          ) : (
            <Password
              defaultValue=""
              name="password"
              placeholder={Intl.t('backup.inputs.password.placeholder')}
              error={error?.password}
            />
          )}
        </div>
        <div className="flex flex-row gap-4 pt-4">
          {privateKey ? <CopyKey privateKey={privateKey} /> : <EnterKey />}
        </div>
      </form>
    </Layout>
  );
}

export default BackupKeyPairs;
