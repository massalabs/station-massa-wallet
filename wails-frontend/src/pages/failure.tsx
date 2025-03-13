import { FiX } from 'react-icons/fi';
import { useLocation } from 'react-router';

import { promptRequest } from '@/events';
import { Layout } from '@/layouts/Layout/Layout';

const Failure = () => {
  const locationState = useLocation().state as { req: promptRequest };
  const req = locationState.req;

  return (
    <Layout>
      <div className="flex flex-col items-center justify-center">
        <div className="w-12 h-12 bg-s-error flex flex-col justify-center items-center rounded-full mb-6">
          <FiX className="w-6 h-6" />
        </div>
        <p className="text-neutral text-center mas-body">{req.Msg}</p>
      </div>
    </Layout>
  );
};

export default Failure;
