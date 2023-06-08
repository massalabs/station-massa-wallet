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
  const { fees, modal, setModal, handleFees, handleConfirm } = props;
  const [presetGasFees, setPresetGasFees] = useState<boolean>(true);
  const [customGasFees, setCustomFees] = useState<boolean>(false);

  function PresetFeeSelector({ name }: { name: string }) {
    const isDisabled = fees !== presetFees[name] || customGasFees;
    const disabledButton = 'text-neutral bg-secondary';

    return (
      <Button
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
  }

  function handleCustomGas() {
    setCustomFees(true);
    setPresetGasFees(false);
  }

  const presetArgs = {
    checked: presetGasFees,
  };

  const customArgs = {
    checked: customGasFees,
  };

  return (
    <PopupModal
      fullMode={true}
      onClose={() => setModal(!modal)}
      customClass="!w-1/2 min-w-[775px]"
    >
      <PopupModalHeader>
        <div className="flex flex-col gap-3.5">
          <label className="mas-title mb-6">
            {Intl.t('sendcoins.advanced')}
          </label>
          <p className="mas-body mb-6">
            {Intl.t('sendcoins.advanced-message')}
          </p>
        </div>
      </PopupModalHeader>
      <PopupModalContent>
        <div className="pb-10">
          <div className="flex flex-row items-center mas-buttons mb-3">
            <div className="pl-3">
              <RadioButton
                defaultChecked={true}
                onClick={() => handlePresetGas()}
                {...presetArgs}
              />
            </div>
            <p>{Intl.t('sendcoins.preset')}</p>
          </div>
          <div className="flex flex-row items-center w-full gap-4 mb-6">
            {feesTypes.map((type) => (
              <PresetFeeSelector key={type} name={type} />
            ))}
          </div>

          <div className="flex flex-row items-center mas-buttons mb-3">
            <div className="pl-3">
              <RadioButton
                defaultChecked={false}
                onClick={() => handleCustomGas()}
                {...customArgs}
              />
            </div>
            <p>{Intl.t('sendcoins.custom-fees')}</p>
          </div>
          <Input
            type="text"
            placeholder="Gas fees amount"
            name="fees"
            defaultValue=""
            disabled={!customGasFees}
            onChange={(e) => handleFees(e.target.value)}
          />

          <Button
            customClass="mt-6"
            type="submit"
            onClick={(e) => handleConfirm(e)}
          >
            {Intl.t('sendcoins.confirm-fees')}
          </Button>
        </div>
      </PopupModalContent>
    </PopupModal>
  );
}

export default Modal;
