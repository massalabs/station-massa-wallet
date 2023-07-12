import {
  Button,
  Currency,
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
  RadioButton,
} from '@massalabs/react-ui-kit';
import Intl from '@/i18n/i18n';
import { useState, FormEvent } from 'react';
import { parseForm } from '@/utils/parseForm';

interface InputsErrors {
  fees?: string;
}

const PRESET_LOW = '1';
const PRESET_STANDARD = '1000';
const PRESET_HIGH = '5000';

export const presetFees: { [key: string]: string } = {
  low: PRESET_LOW,
  standard: PRESET_STANDARD,
  high: PRESET_HIGH,
};

function Advanced({ ...props }) {
  const { onClose, fees: currentFees, setFees: setCurrentFees } = props;

  const isOneOfPressedFees = Object.values(presetFees).includes(currentFees);
  const initialFees = !isOneOfPressedFees ? currentFees : 0;
  const initialPressetFees = !isOneOfPressedFees
    ? '0'
    : currentFees || PRESET_STANDARD;
  const initialCustomFees = !isOneOfPressedFees || false;

  const [error, setError] = useState<InputsErrors | null>(null);
  const [fees, setFees] = useState<number | string | undefined>(initialFees);
  const [customFees, setCustomFees] = useState<boolean>(initialCustomFees);
  const [presetFee, setPresetFee] = useState<string>(initialPressetFees);

  function handleGasFeesOption(isCustomFees: boolean) {
    setCustomFees(isCustomFees);
    setError(null);
    setFees(0);

    isCustomFees ? setPresetFee('0') : setPresetFee(PRESET_STANDARD);
  }

  function validate(formObject: any) {
    const { fees } = formObject;
    setError(null);
    if (customFees && !fees) {
      setError({ fees: Intl.t('errors.send-coins.no-gas-fees') });
      return false;
    }

    if (customFees && fees <= 0) {
      setError({ fees: Intl.t('errors.send-coins.gas-fees-to-low') });
      return false;
    }

    return true;
  }

  function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const formObject = parseForm(e);

    if (!validate(formObject)) return;

    let pickedFee = Number(fees) > 0 ? fees : presetFee;

    setCurrentFees(pickedFee);
    onClose?.(pickedFee);
  }

  function FeesSelector({ name }: { name: string }) {
    const isDisabled = presetFee !== presetFees[name];
    return (
      <Button
        disabled={customFees}
        name={name}
        customClass={
          isDisabled
            ? 'w-1/3 h-12 text-neutral bg-secondary'
            : 'w-1/3 h-12 bg-neutral'
        }
        variant="toggle"
        onClick={() => setPresetFee(presetFees[name])}
      >
        {name}
        <label className="text-tertiary text-xs flex pl-1 items-center cursor-pointer">
          ({presetFees[name]} nMAS)
        </label>
      </Button>
    );
  }

  return (
    <PopupModal
      fullMode={true}
      onClose={onClose}
      customClass="!w-1/2 min-w-[775px]"
    >
      <PopupModalHeader>
        <div className="flex flex-col gap-3.5">
          <label className="mas-title mb-6">
            {Intl.t('send-coins.advanced')}
          </label>
          <p className="mas-body mb-6">
            {Intl.t('send-coins.advanced-message')}
          </p>
        </div>
      </PopupModalHeader>
      <PopupModalContent>
        <div className="pb-10">
          <form onSubmit={handleSubmit}>
            <div className="flex flex-row items-center mas-buttons mb-3">
              <RadioButton
                checked={!customFees}
                onChange={() => handleGasFeesOption(false)}
                name="gas"
              />
              <p
                className="h-full ml-3 pb-1 cursor-pointer"
                onClick={() => handleGasFeesOption(false)}
              >
                {Intl.t('send-coins.preset')}
              </p>
            </div>
            <div className="flex flex-row items-center w-full gap-4 mb-6">
              {Object.keys(presetFees).map((type) => (
                <FeesSelector key={type} name={type} />
              ))}
            </div>

            <div className="flex flex-row items-center mas-buttons mb-3">
              <RadioButton
                checked={customFees}
                onChange={() => handleGasFeesOption(true)}
                name="gas"
              />
              <p
                className="h-full ml-3 pb-1 cursor-pointer"
                onClick={() => handleGasFeesOption(true)}
              >
                {Intl.t('send-coins.custom-fees')}:
              </p>
            </div>
            <Currency
              placeholder={Intl.t('send-coins.custom-fees')}
              name="fees"
              variant="nMAS"
              value={fees === 0 ? '' : fees}
              disabled={!customFees}
              onValueChange={(value) => setFees(value)}
              error={error?.fees}
            />

            <Button customClass="mt-6" type="submit">
              {Intl.t('send-coins.confirm-fees')}
            </Button>
          </form>
        </div>
      </PopupModalContent>
    </PopupModal>
  );
}

export default Advanced;
