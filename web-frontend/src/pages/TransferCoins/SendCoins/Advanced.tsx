import {
  Button,
  Input,
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
  RadioButton,
} from '@massalabs/react-ui-kit';
import { presetFees } from '../../../utils/MassaFormating';
import Intl from '../../../i18n/i18n';
import { useState } from 'react';

function Modal({ ...props }) {
  const feesTypes = Object.keys(presetFees);
  const {
    fees,
    modal,
    setModal,
    handleFees,
    handleConfirm,
    setErrorAdvanced,
    errorAdvanced,
  } = props;
  const [presetGasFees, setPresetGasFees] = useState<boolean>(true);
  const [customGasFees, setCustomFees] = useState<boolean>(false);

  function PresetFeeSelector({ name }: { name: string }) {
    const isDisabled = fees !== presetFees[name] || customGasFees;
    const disabledButton = 'text-neutral bg-secondary';

    const disabledButtonArgs = {
      disabled: customGasFees,
    };

    return (
      <Button
        {...disabledButtonArgs}
        customClass={
          isDisabled ? `w-1/3 h-12 ${disabledButton}` : `w-1/3 h-12 bg-neutral`
        }
        variant="toggle"
        onClick={() => handleFees(presetFees[name])}
      >
        {name}
        <label className="text-tertiary text-xs flex pl-1 items-center">
          ({presetFees[name]} nMAS)
        </label>
      </Button>
    );
  }

  function handlePresetGas() {
    setCustomFees(false);
    setPresetGasFees(true);
    handleFees('1000');
    setErrorAdvanced(null);
  }

  function handleCustomGas() {
    setCustomFees(true);
    setPresetGasFees(false);
    handleFees('');
  }

  const presetArgs = {
    checked: presetGasFees,
  };

  const customArgs = {
    checked: customGasFees,
  };

  function handleClose() {
    setModal(!modal);
    setErrorAdvanced(null);
    handleFees('1000');
  }

  function handleOpen() {
    handleFees('1000');
  }

  return (
    <PopupModal
      fullMode={true}
      onClose={handleClose}
      onOpen={handleOpen}
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
          <div className="flex flex-row items-center mas-buttons mb-3">
            <RadioButton
              defaultChecked={true}
              onClick={() => handlePresetGas()}
              {...presetArgs}
            />
            <p className="h-full ml-3 pb-1">{Intl.t('send-coins.preset')}</p>
          </div>
          <div className="flex flex-row items-center w-full gap-4 mb-6">
            {feesTypes.map((type) => (
              <PresetFeeSelector key={type} name={type} />
            ))}
          </div>

          <div className="flex flex-row items-center mas-buttons mb-3">
            <RadioButton
              defaultChecked={false}
              onClick={() => handleCustomGas()}
              {...customArgs}
            />
            <p className="h-full ml-3 pb-1">
              {Intl.t('send-coins.custom-fees')}
            </p>
          </div>
          <Input
            type="text"
            placeholder="Gas fees amount (nMAS)"
            name="fees"
            value={!customGasFees ? '' : fees}
            disabled={!customGasFees}
            onChange={(e) => handleFees(e.target.value)}
            error={errorAdvanced?.amount}
          />

          <Button
            customClass="mt-6"
            type="submit"
            onClick={(e) => handleConfirm(e)}
          >
            {Intl.t('send-coins.confirm-fees')}
          </Button>
        </div>
      </PopupModalContent>
    </PopupModal>
  );
}

export default Modal;
