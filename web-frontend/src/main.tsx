import React from 'react';
import ReactDOM from 'react-dom/client';
import {
  RouterProvider,
  createBrowserRouter,
  createRoutesFromElements,
  Route,
} from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import '@massalabs/react-ui-kit/src/global.css';

import { ENV } from './const/env/env';
import './index.css';
import Welcome from './pages/Welcome/Welcome.tsx';
import SelectAccount from './pages/SelectAccount/SelectAccount.tsx';
import Error from './pages/Error.tsx';
import AddAccount from './pages/AddAccount/AddAccount.tsx';
import StepTwo from './pages/CreateAccount/StepTwo.tsx';
import StepOne from './pages/CreateAccount/StepOne.tsx';
import StepThree from './pages/CreateAccount/StepThree.tsx';
import mockServer from './mirage/server.ts';

// Add ENV.STANDALONE to the array to enable MirageJS
if ([ENV.DEV, ENV.TEST].includes(import.meta.env.VITE_ENV)) {
  mockServer(import.meta.env.VITE_ENV);
}

const queryClient = new QueryClient();

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path={import.meta.env.VITE_BASE_APP}>
      <Route path="index" element={<Welcome />} />
      <Route path="account-select" element={<SelectAccount />} />
      <Route path="account-create" element={<AddAccount />} />
      <Route path="account-create-step-one" element={<StepOne />} />
      <Route path="account-create-step-two" element={<StepTwo />} />
      <Route path="account-create-step-three" element={<StepThree />} />
      <Route path="error" element={<Error />} />
      <Route path="*" element={<Error />} />
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
