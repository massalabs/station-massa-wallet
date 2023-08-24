import React from 'react';

import ReactDOM from 'react-dom/client';
import {
  Route,
  createBrowserRouter,
  createRoutesFromElements,
  RouterProvider,
} from 'react-router-dom';
import './index.css';
import '@massalabs/react-ui-kit/src/global.css';

import PasswordPromptHandler from './pages/PasswordPromptHandler/PasswordPromptHandler';
import App from '@/app';
import BackupKeyPairs from '@/pages/backupKeyPairs';
import BackupMethods from '@/pages/backupMethods';
import ConfirmDelete from '@/pages/confirmDelete';
import Failure from '@/pages/failure';
import ImportFile from '@/pages/importFile';
import ImportMethods from '@/pages/importMethods';
import PromptNickname from '@/pages/ImportPrivateKey/PromptNickname';
import PromptPrivateKey from '@/pages/ImportPrivateKey/PromptPrivateKey';
import NewPassword from '@/pages/newPassword';
import Success from '@/pages/success';

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route>
      <Route path="/" element={<App />} />
      <Route path="/backup-methods" element={<BackupMethods />} />
      <Route path="/backup-pkey" element={<BackupKeyPairs />} />
      <Route path="/password" element={<PasswordPromptHandler />} />
      <Route path="/new-password" element={<NewPassword />} />
      <Route path="/success" element={<Success />} />
      <Route path="/failure" element={<Failure />} />
      <Route path="/import-methods" element={<ImportMethods />} />
      <Route path="/import-file" element={<ImportFile />} />
      <Route path="/import-key-pairs" element={<PromptPrivateKey />} />
      <Route path="/import-nickname" element={<PromptNickname />} />
      <Route path="/confirm-delete" element={<ConfirmDelete />} />
    </Route>,
  ),
);

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
);
