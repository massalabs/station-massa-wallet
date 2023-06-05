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
  const { fees, modal, setModal, handleFees } = props;

  function PresetFeeSelector({ name }: { name: string }) {
    return (
      <Button
        disabled={!radioActive} // Update the disabled prop here
        variant={fees === presetFees[name] ? 'primary' : 'secondary'}
        onClick={() => handleFees(presetFees[name])}
      >
        {name}
        <label className="text-tertiary text-xs flex pl-1 items-center">
          ({presetFees[name]} nMAS)
        </label>
      </Button>
    );
  }

  const radioArgs = {
    name: 'radio',
  };

  const [radioActive, setRadioActive] = useState<boolean>(true);

  return (
    <PopupModal fullMode={true} onClose={() => setModal(!modal)}>
      <PopupModalHeader>
        <div className="flex flex-col gap-3.5">
          <label className="mas-title">{Intl.t('sendcoins.advanced')}</label>
          <p className="mas-body">{Intl.t('sendcoins.advanced-message')}</p>
        </div>
      </PopupModalHeader>
      <PopupModalContent>
        <div className="flex flex-row items-center mas-buttons pb-3.5">
          <RadioButton
            defaultChecked={true}
            onClick={() => setRadioActive(!radioActive)}
            {...radioArgs}
          />
          <p>{Intl.t('sendcoins.preset')}</p>
        </div>
        <div className="flex flex-row items-center w-full gap-4 pb-3.5">
          {feesTypes.map((type) => (
            <PresetFeeSelector key={type} name={type} />
          ))}
        </div>

        <div className="flex flex-row items-center mas-buttons mb-3.5">
          <RadioButton
            defaultChecked={false}
            onClick={() => setRadioActive(!radioActive)}
            {...radioArgs}
          />
          <p>{Intl.t('sendcoins.custom-fees')}</p>
        </div>

        <Input
          type="text"
          placeholder="Gas fees amount"
          name="fees"
          defaultValue=""
          disabled={radioActive}
          onChange={(e) => handleFees(e.target.value)}
        />

        <Button type="submit" onClick={() => setModal(!modal)}>
          {Intl.t('sendcoins.confirm-fees')}
        </Button>
      </PopupModalContent>
    </PopupModal>
  );
}

export default Modal;
