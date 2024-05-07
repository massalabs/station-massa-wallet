import { useState, FormEvent } from 'react';

import {
  Button,
  Money,
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
  RadioButton,
} from '@massalabs/react-ui-kit';

import Intl from '@/i18n/i18n';
import { parseForm } from '@/utils/parseForm';

interface InputsErrors {
  fees?: string;
}

interface FeesForm {
  fees: string;
}

export const PRESET_LOW = '0.01';
export const PRESET_STANDARD = '0.03';
export const PRESET_HIGH = '0.1';

export const presetFees: { [key: string]: string } = {
  low: PRESET_LOW,
  standard: PRESET_STANDARD,
  high: PRESET_HIGH,
};

export interface AdvancedProps {
  onClose: () => void;
  fees: string;
  setFees: (fees: string) => void;
}

export default function Advanced(props: AdvancedProps) {
  const { onClose, fees: currentFees, setFees: setCurrentFees } = props;

  const isOneOfPressedFees = Object.values(presetFees).includes(currentFees);
  const initialFees = !isOneOfPressedFees ? currentFees : PRESET_LOW;
  const initialPresetFees = !isOneOfPressedFees
    ? PRESET_LOW
    : currentFees || PRESET_LOW;
  const initialCustomFees = !isOneOfPressedFees;

  const [error, setError] = useState<InputsErrors | null>(null);
  const [fees, setFees] = useState<string>(initialFees);
  const [feesField, setFeesField] = useState<string>('');
  const [customFees, setCustomFees] = useState<boolean>(initialCustomFees);
  const [presetFee, setPresetFee] = useState<string>(initialPresetFees);

  function handleGasFeesOption(isCustomFees: boolean) {
    setCustomFees(isCustomFees);
    setError(null);
    setFees(PRESET_LOW);
    setPresetFee(PRESET_LOW);
  }

  function validate(formObject: FeesForm) {
    const { fees } = formObject;
    setError(null);
    if (customFees && !fees) {
      setError({ fees: Intl.t('errors.send-coins.no-fees') });
      return false;
    }

    return true;
  }

  function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const formObject = parseForm(e) as FeesForm;

    if (!validate(formObject)) return;

    setCurrentFees(fees);
    onClose();
  }

  function FeesSelector({ name }: { name: string }) {
    const isDisabled = presetFee !== presetFees[name];
    function handleClick() {
      setPresetFee(presetFees[name]);
      setFees(presetFees[name]);
    }
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
        onClick={handleClick}
      >
        {name}
        <label className="text-tertiary text-xs flex pl-1 items-center cursor-pointer">
          ({presetFees[name].toString()} MAS)
        </label>
      </Button>
    );
  }

  function onFeeChange(event: { value: string }) {
    setFees(event.value);
    setFeesField(event.value);
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
            <Money
              placeholder={Intl.t('send-coins.custom-fees')}
              name="fees"
              variant="MAS"
              value={(customFees && currentFees) || feesField}
              disabled={!customFees}
              onValueChange={(event) => onFeeChange(event)}
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
