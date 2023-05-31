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

import App from './app';
import PasswordPrompt from './pages/passwordPrompt';
import Success from './pages/success';
import ImportMethods from './pages/importMethods';
import ImportFile from './pages/importFile';
import PromptPrivateKey from './pages/ImportPrivateKey/PromptPrivateKey';
import PromptNickname from './pages/ImportPrivateKey/PromptNickname';
import Failure from './pages/failure';
import BackupMethods from './pages/backupMethods';
import BackupKeyPairs from './pages/backupKeyPairs';
import NewPassword from './pages/newPassword';

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route>
      <Route path="/" element={<App />} />
      <Route path="/backup-methods" element={<BackupMethods />} />
      <Route path="/backup-pkey" element={<BackupKeyPairs />} />
      <Route path="/password" element={<PasswordPrompt />} />
      <Route path="/new-password" element={<NewPassword />} />
      <Route path="/success" element={<Success />} />
      <Route path="/failure" element={<Failure />} />
      <Route path="/import-methods" element={<ImportMethods />} />
      <Route path="/import-file" element={<ImportFile />} />
      <Route path="/import-key-pairs" element={<PromptPrivateKey />} />
      <Route path="/import-nickname" element={<PromptNickname />} />
    </Route>,
  ),
);

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
);
