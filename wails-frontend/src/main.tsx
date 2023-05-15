import React from 'react';
import ReactDOM from 'react-dom/client';
import {
  Route,
  createBrowserRouter,
  createRoutesFromElements,
  RouterProvider,
} from 'react-router-dom';
import './index.css';

import App from './app';
import PasswordPrompt from './pages/passwordPrompt';
import Success from './pages/success';
import ImportMethods from './pages/importMethods';
import ImportFile from './pages/importFile';
import ImportPrivatekey from './pages/importPrivateKey';

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route>
      <Route path="/" element={<App />} />
      <Route path="/password" element={<PasswordPrompt />} />
      <Route path="/success" element={<Success />} />
      <Route path="/import-methods" element={<ImportMethods />} />
      <Route path="/import-file" element={<ImportFile />} />
      <Route path="/import-pkey" element={<ImportPrivatekey />} />
    </Route>,
  ),
);

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
);
