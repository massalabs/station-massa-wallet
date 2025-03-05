import React from 'react';

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import ReactDOM from 'react-dom/client';
import {
  RouterProvider,
  createBrowserRouter,
  createRoutesFromElements,
  Route,
  Navigate,
} from 'react-router-dom';

import '@massalabs/react-ui-kit/src/global.css';
import './index.css';
import { ENV } from '@/const/env/env';
import { mockServer, mockServerWithCypress } from '@/mirage';
import AccountCreate from '@/pages/AccountCreate/AccountCreate.tsx';
import AccountSelect from '@/pages/AccountSelect/AccountSelect.tsx';
import Assets from '@/pages/Assets/Assets.tsx';
import Base from '@/pages/Base/Base.tsx';
import Contacts from '@/pages/Contacts/Contacts.tsx';
import StepOne from '@/pages/CreateAccount/StepOne.tsx';
import StepThree from '@/pages/CreateAccount/StepThree.tsx';
import StepTwo from '@/pages/CreateAccount/StepTwo.tsx';
import Error from '@/pages/Error.tsx';
import Home from '@/pages/Home/Home.tsx';
import Index from '@/pages/Index/Index.tsx';
import Settings from '@/pages/Settings/Settings.tsx';
import SettingsUpdate from '@/pages/Settings/SettingsUpdate.tsx';
import Transactions from '@/pages/Transactions/Transactions.tsx';
import SendRedirect from '@/pages/TransferCoins/SendCoins/SendRedirect.tsx';
import TransferCoins from '@/pages/TransferCoins/TransferCoins.tsx';

const baseURL = import.meta.env.VITE_BASE_APP;
const baseENV = import.meta.env.VITE_ENV;

// Add ENV.STANDALONE to the array to enable MirageJS
if ([ENV.DEV, ENV.TEST].includes(baseENV)) {
  mockServer(import.meta.env.VITE_ENV);
  mockServerWithCypress();
}

// save the mode in context. It is used in wallet-provider
window.massaWallet = {
  standalone: import.meta.env.VITE_ENV === ENV.STANDALONE,
};

const queryClient = new QueryClient();

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path={baseURL} element={<Base />}>
      {/* routes for onboarding and account creation */}
      <Route path="index" element={<Index />} />
      <Route path="account-select" element={<AccountSelect />} />
      <Route path="account-create" element={<AccountCreate />} />
      <Route path="account-create-step-one" element={<StepOne />} />
      <Route path="account-create-step-two" element={<StepTwo />} />
      <Route path="account-create-step-three" element={<StepThree />} />
      <Route path="send-redirect" element={<SendRedirect />} />

      {/* routes for wallet */}
      <Route path=":nickname/home" element={<Home />} />
      <Route path=":nickname/transactions" element={<Transactions />} />
      <Route path=":nickname/transfer-coins" element={<TransferCoins />} />
      <Route path=":nickname/assets" element={<Assets />} />
      <Route path=":nickname/settings" element={<Settings />} />
      <Route path=":nickname/settings/update" element={<SettingsUpdate />} />
      <Route path=":nickname/contacts" element={<Contacts />} />

      {/* routes for errors */}
      <Route path="error" element={<Error />} />
      <Route path="*" element={<Error />} />
      <Route
        path={`${baseURL}/`}
        element={<Navigate to={`${baseURL}/index`} />}
      />
    </Route>,
  ),
);

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} fallbackElement={<Error />} />
    </QueryClientProvider>
  </React.StrictMode>,
);
