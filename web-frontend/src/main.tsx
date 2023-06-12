import React from 'react';

import ReactDOM from 'react-dom/client';
import {
  RouterProvider,
  createBrowserRouter,
  createRoutesFromElements,
  Route,
  Navigate,
} from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import '@massalabs/react-ui-kit/src/global.css';

import { ENV } from './const/env/env';
import './index.css';
import mockServer from './mirage/server.js';
import Index from './pages/Index/Index.tsx';
import SelectAccount from './pages/SelectAccount/SelectAccount.tsx';
import Error from './pages/Error.tsx';
import AddAccount from './pages/AddAccount/AddAccount.tsx';
import StepTwo from './pages/CreateAccount/StepTwo.tsx';
import StepOne from './pages/CreateAccount/StepOne.tsx';
import StepThree from './pages/CreateAccount/StepThree.tsx';
import SendCoins from './pages/TransferCoins/SendCoins.tsx';
import Home from './pages/Home/Home.tsx';
import ReceiveCoins from './pages/TransferCoins/ReceiveCoins/ReceiveCoins.tsx';
import Settings from './pages/Settings/Settings.tsx';
import SettingsEdit from './pages/Settings/SeetingsEdit.tsx';
import Transactions from './pages/Transactions/Transactions.tsx';
import Contacts from './pages/Contacts/Contacts.tsx';
import Assets from './pages/Assets/Assets.tsx';
import SendRedirect from './pages/TransferCoins/SendCoins/SendRedirect.tsx';
import Base from './pages/Base/Base.tsx';

// Add ENV.STANDALONE to the array to enable MirageJS
if ([ENV.DEV, ENV.TEST].includes(import.meta.env.VITE_ENV)) {
  mockServer(import.meta.env.VITE_ENV);
}

const queryClient = new QueryClient();

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path={import.meta.env.VITE_BASE_APP} element={<Base />}>
      {/* routes for onboarding and account creation */}
      <Route path="index" element={<Index />} />
      <Route path="account-select" element={<SelectAccount />} />
      <Route path="account-create" element={<AddAccount />} />
      <Route path="account-create-step-one" element={<StepOne />} />
      <Route path="account-create-step-two" element={<StepTwo />} />
      <Route path="account-create-step-three" element={<StepThree />} />
      <Route path="send-redirect" element={<SendRedirect />} />

      {/* routes for wallet */}
      <Route path=":nickname/home" element={<Home />} />
      <Route path=":nickname/transactions" element={<Transactions />} />
      <Route path=":nickname/send-coins" element={<SendCoins />} />
      <Route path=":nickname/receive" element={<ReceiveCoins />} />
      <Route path=":nickname/assets" element={<Assets />} />
      <Route path=":nickname/settings" element={<Settings />} />
      <Route path=":nickname/settings/update" element={<SettingsEdit />} />
      <Route path=":nickname/contacts" element={<Contacts />} />

      {/* routes for errors */}
      <Route path="error" element={<Error />} />
      <Route path="*" element={<Error />} />
      <Route path="/" element={<Navigate to="/index" />} />
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
