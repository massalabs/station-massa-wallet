import './App.css'
import { h } from 'preact';
import { EventsOn } from "../wailsjs/runtime";
import PasswordPrompt from './pages/passwordPrompt';
import { events, promptRequest } from './events/events';
import { Routes, Route, useNavigate } from "react-router-dom";
import Success from './pages/success';

export function App() {

    const navigate = useNavigate();

    const handlePromptRequest = (req: promptRequest) => {
        navigate("/password", { state: { req } } )
    }

    EventsOn(events.promptRequest, handlePromptRequest)

    return (
        <div id="App">
            <Routes>
                <Route path="/password" element={<PasswordPrompt />} />
                <Route path="/success" element={<Success />} />
            </Routes>
        </div>
    )
}
