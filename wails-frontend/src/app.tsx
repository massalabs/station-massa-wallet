import './App.css';
import { h } from 'preact';
import { EventsOn } from '../wailsjs/runtime';
import PasswordPrompt from './pages/passwordPrompt';
import { events, promptAction, promptRequest } from './events/events';
import { Routes, Route, useNavigate } from 'react-router-dom';
import Success from './pages/success';
import ImportMethods from './pages/importMethods';
import ImportFile from './pages/importFile';
import ImportPrivatekey from './pages/importPrivateKey';

export function App() {
  const navigate = useNavigate();

  const handlePromptRequest = (req: promptRequest) => {
    switch (req.Action) {
      case promptAction.deleteReq:
      case promptAction.newPasswordReq:
        navigate('/password', { state: { req } });
        return;
      case promptAction.importReq:
        navigate('/import-methods', { state: { req } });
        return;
      default:
    }
  };

  EventsOn(events.promptRequest, handlePromptRequest);

  return (
    <div id="app">
      <Routes>
        <Route path="/password" element={<PasswordPrompt />} />
        <Route path="/success" element={<Success />} />
        <Route path="/import-methods" element={<ImportMethods />} />
        <Route path="/import-file" element={<ImportFile />} />
        <Route path="/import-pkey" element={<ImportPrivatekey />} />
      </Routes>
    </div>
  );
}
