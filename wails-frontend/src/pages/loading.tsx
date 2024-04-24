import { Spinner } from '@massalabs/react-ui-kit';

import { Layout } from '@/layouts/Layout/Layout';

export const Loading = () => {
  return (
    <Layout>
      <div className="flex flex-col items-center justify-center">
        <div className="w-24 h-24 flex flex-col justify-center items-center rounded-full mb-6">
          <Spinner size={38} />
        </div>
      </div>
    </Layout>
  );
};
