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
  const { fees, handleModal, handleFees } = props;

  function PresetFeeSelector({ name }: { name: string }) {
    return (
      <div className="w-full">
        <Button
          disabled={!radioActive} // Update the disabled prop here
          variant={fees === presetFees[name] ? 'primary' : 'secondary'}
          onClick={() => handleFees(presetFees[name])}
        >
          {name}
          <label className="text-info text-xs flex pl-1 items-center">
            ({presetFees[name]} nMAS)
          </label>
        </Button>
      </div>
    );
  }

  const radioArgs = {
    name: 'radio',
  };

  const [radioActive, setRadioActive] = useState<boolean>(true);

  function handleRadio() {
    setRadioActive(!radioActive);
  }

  return (
    <PopupModal fullMode={true} onClose={() => handleModal()}>
      <PopupModalHeader>
        <div>
          <label className="mas-title">{Intl.t('sendcoins.advanced')}</label>
          <p className="mas-body">{Intl.t('sendcoins.advanced-message')}</p>
        </div>
      </PopupModalHeader>
      <PopupModalContent>
        <div className="flex flex-row items-center mas-buttons pb-3.5">
          <RadioButton
            defaultChecked={true}
            onClick={handleRadio}
            {...radioArgs}
          />
          <p>{Intl.t('sendcoins.preset')}</p>
        </div>
        <div className="flex flex-row items-center w-full gap-4 pb-3.5">
          {feesTypes.map((type) => (
            <PresetFeeSelector key={type} name={type} />
          ))}
        </div>

        <div className="flex flex-row items-center mas-buttons pb-3.5">
          <RadioButton
            defaultChecked={false}
            onClick={handleRadio}
            {...radioArgs}
          />
          <p>{Intl.t('sendcoins.custom-fees')}</p>
        </div>

        <div className="pt-3.5">
          <Input
            type="text"
            placeholder="Gas fees amount"
            name="fees"
            defaultValue=""
            disabled={radioActive}
            onChange={(e) => handleFees(e.target.value)}
          />
        </div>
        <Button type="submit" onClick={() => handleModal()}>
          {Intl.t('sendcoins.confirm-fees')}
        </Button>
      </PopupModalContent>
    </PopupModal>
  );
}

export default Modal;
