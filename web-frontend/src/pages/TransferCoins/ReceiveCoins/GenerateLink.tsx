import {
  Input,
  MassaToken,
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
  Selector,
} from '@massalabs/react-ui-kit';
import { FiArrowUpRight } from 'react-icons/fi';

function GenerateLink({ ...props }) {
  const { modal, setModal } = props;
  return (
    <>
      <PopupModal fullMode={true} onClose={() => setModal(!modal)}>
        <PopupModalHeader>
          <h1 className="mas-banner">Generate a link</h1>
        </PopupModalHeader>
        <PopupModalContent>
          <div className="flex flex-col">
            <p>Token to send</p>
            <Input placeholder="Amount to ask" />
          </div>
          <div className="flex flex-col">
            <p>Account to provide</p>
            <Selector
              preIcon={<FiArrowUpRight size={24} />}
              content="Account 1"
              amount="0.00000000"
              posIcon={<MassaToken />}
            />
          </div>
          <div className="flex flex-col">
            <p>Link to share</p>
          </div>
        </PopupModalContent>
      </PopupModal>
    </>
  );
}

export default GenerateLink;
