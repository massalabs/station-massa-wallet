import {
  Button,
  Input,
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
  RadioButton,
} from '@massalabs/react-ui-kit';
import { presetFees } from '../../../utils/MassaFormating';
import { useState } from 'react';

function Modal({ ...props }) {
  const feesTypes = Object.keys(presetFees);
  const { modal, setModal, fees, handleFees, handleFeesConfirm } = props;

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
    console.log('click');
    setRadioActive(!radioActive);
    console.log(radioActive);
  }

  return (
    <PopupModal fullMode={true} onClose={() => setModal(!modal)}>
      <PopupModalHeader>
        <div>
          <label className="mas-title">Advanced</label>
          <p className="mas-body">
            You pay gas fees to reward block validators and maximize your
            chances to see your transaction validated. It is a tip for people
            that support the blockchain network.
          </p>
        </div>
      </PopupModalHeader>
      <PopupModalContent>
        <div>
          <div className="flex flex-row items-center mas-buttons pb-3.5">
            <RadioButton
              defaultChecked={true}
              onClick={handleRadio}
              {...radioArgs}
            />
            <p>Preset :</p>
          </div>
          <div className="flex flex-row items-center w-full gap-4 pb-3.5">
            {feesTypes.map((type) => (
              <PresetFeeSelector key={type} name={type} />
            ))}
          </div>
          <div>
            <div className="flex flex-row items-center mas-buttons pb-3.5">
              <RadioButton
                defaultChecked={false}
                onClick={handleRadio}
                {...radioArgs}
              />
              <p>Custom fees :</p>
            </div>

            <Input
              type="text"
              placeholder="Gas fees amount"
              name="fees"
              defaultValue=""
              disabled={radioActive}
              onChange={(e) => handleFees(e.target.value)}
            />
            <div>
              <Button type="submit" onClick={handleFeesConfirm}>
                Confirm fees
              </Button>
            </div>
          </div>
        </div>
      </PopupModalContent>
    </PopupModal>
  );
}

export default Modal;
