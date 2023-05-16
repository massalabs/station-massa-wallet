import React from 'react';
import ReactDOM from 'react-dom/client';
import {
  RouterProvider,
  createBrowserRouter,
  createRoutesFromElements,
  Route,
} from 'react-router-dom';
import './index.css';
import '@massalabs/react-ui-kit/src/global.css';
import Welcome from './pages/Welcome/Welcome.tsx';
import SelectAccount from './pages/SelectAccount/SelectAccount.tsx';
import Error from './pages/Error.tsx';
import AddAccount from './pages/AddAccount/AddAccount.tsx';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import StepTwo from './pages/CreateAccount/StepTwo.tsx';
import mockServer from './mirage/server.ts';

if (import.meta.env.VITE_ENV === 'dev') {
  mockServer();
}
import StepOne from './pages/CreateAccount/StepOne.tsx';
import StepThree from './pages/CreateAccount/StepThree.tsx';

const queryClient = new QueryClient();

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path={import.meta.env.VITE_BASE_APP}>
      <Route path="welcome" element={<Welcome />} />
      <Route path="account-select" element={<SelectAccount />} />
      <Route path="account-create" element={<AddAccount />} />
      <Route path="account-create-step2" element={<StepTwo />} />
      <Route path="error" element={<Error />} />
      <Route path="account-new" element={<AddAccount />} />
      <Route path="account-create-step1" element={<StepOne />} />
      <Route path="account-create-step-three" element={<StepThree />} />
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
