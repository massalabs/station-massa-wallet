import './App.css'
import logo from "./assets/images/logo_massa.webp"
import { EventsOn } from "../wailsjs/runtime";
import { useState } from "preact/hooks";
import { h } from 'preact';
import PasswordPrompt from './pages/passwordPrompt';
import { events, promptRequest } from './events/events';

export function App(props: any) {
    const [eventData, setEventData] = useState<promptRequest>({} as promptRequest);
    const handlePromptRequest = (data: promptRequest) => {
        console.log("prompt request event received: ", data)
        setEventData(data)
    }
 
    EventsOn(events.promptRequest,handlePromptRequest)

    return (
        <>
            <div id="App">
            <img src={logo} id="logo" alt="logo" />
                <PasswordPrompt eventData={eventData}/>
            </div>
        </>
    )
}
