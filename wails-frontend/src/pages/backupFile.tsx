/* eslint-disable new-cap */
import { useState } from 'react';
import { promptRequest } from '../events/events';
import { useLocation } from 'react-router-dom';
import { handleCancel } from '../utils/utils';
import { Button } from '@massalabs/react-ui-kit/src/components/Button/Button';
import { Password } from '@massalabs/react-ui-kit/src/components/Password/Password';
import { FiCheck } from 'react-icons/fi';

function BackupFile() {
  const [backup, setBackup] = useState<boolean>(false);
  const [input, setInput] = useState<string | undefined>(undefined);
  const [isValid, setIsValid] = useState<boolean>(false);

  const { state } = useLocation();
  const req: promptRequest = state.req;

  function baselineStr() {
    if (!backup) {
      return 'Choose which .yaml file you want';
    } else {
      ('Enter your password to decrypt your file');
    }
  }

  function handleEncryptedYml() {
    console.log('handleEncryptedYml');
  }

  function handleDecryptedYml() {
    console.log('handleDecryptedYml');
    setBackup(!backup);
  }

  function handleDownload() {
    console.log('handleDownload');
  }

  function handleApply() {
    console.log('handleApply');
    setIsValid(!isValid);
  }

  function updateInput(e: React.ChangeEvent<HTMLInputElement>) {
    setInput(e.target.value);
  }

  function passwordConfirm() {
    return (
      <>
        <div>
          <div>
            <p className="mas-title text-neutral pb-4">{req.Msg}</p>
          </div>
          <div className="mas-body text-neutral pb-4">{baselineStr()}</div>
        </div>
        <div id="input" className="pb-4">
          <Password
            id="input"
            onChange={updateInput}
            autoComplete="off"
            name="input"
            placeholder={'Password'}
            value={input}
          />
        </div>
        <div className="flex flex-row gap-4 w-full pb-4">
          <div className="min-w-fit w-full">
            <Button variant={'secondary'} onClick={handleCancel}>
              Cancel
            </Button>
          </div>
          <div className="min-w-fit w-full">
            <Button onClick={handleApply}>Confirm</Button>
          </div>
        </div>
      </>
    );
  }

  function chooseYamlDownload() {
    return (
      <>
        <div>
          <div>
            <p className="mas-title text-neutral pb-4">{req.Msg}</p>
          </div>
          <div className="mas-body text-neutral pb-4">{baselineStr()}</div>
        </div>
        <div className="flex flex-col">
          <div className="pb-4 w-full">
            <Button variant={'secondary'} onClick={handleEncryptedYml}>
              Download an encrypted .yaml file
            </Button>
          </div>
          <div className="w-full">
            <Button onClick={handleDecryptedYml}>
              Download a decrypted .yaml file
            </Button>
          </div>
        </div>
      </>
    );
  }

  return (
    <div className="bg-primary flex flex-col justify-center items-center h-screen w-full">
      <div className="w-1/4  max-w-xs  min-w-fit">
        {isValid ? (
          <>
            <div className="bg-primary flex flex-col justify-center items-center h-screen w-full">
              <div className="w-1/4  max-w-xs  min-w-fit flex flex-col justify-center items-center">
                <div className="w-12 h-12 bg-brand flex flex-col justify-center items-center rounded-full mb-6">
                  <FiCheck className="w-6 h-6" />
                </div>
                <div>
                  <p className="text-neutral mas-body">
                    The file has been decrypted
                  </p>
                </div>
                <div className="flex flex-row gap-4 w-full pb-4">
                  <div className="min-w-fit w-full">
                    <Button onClick={handleDownload}>Download the .yaml</Button>
                  </div>
                </div>
              </div>
            </div>
          </>
        ) : backup ? (
          passwordConfirm()
        ) : (
          chooseYamlDownload()
        )}
      </div>
    </div>
  );
}

export default BackupFile;
