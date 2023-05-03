import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter, Routes, Route } from 'react-router-dom';

import './index.css';
import '@massalabs/react-ui-kit/src/global.css';
import Welcome from './pages/Welcome.tsx';

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <BrowserRouter basename={import.meta.env.VITE_BASE_PATH}>
      <Routes>
        <Route path="index" element={<Welcome />} />
      </Routes>
    </BrowserRouter>
  </React.StrictMode>,
);
