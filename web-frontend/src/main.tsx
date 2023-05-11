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

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path={import.meta.env.VITE_BASE_PATH}>
      <Route path="index" element={<Welcome />} />
      <Route path="account-select" element={<SelectAccount />} />
      <Route path="account-new" element={<AddAccount />} />
      <Route path="*" element={<Error />} />
    </Route>,
  ),
);

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <RouterProvider router={router} fallbackElement={<Error />} />
  </React.StrictMode>,
);
