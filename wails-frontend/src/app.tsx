import './App.css'
import { h } from 'preact';
import { EventsOn } from "../wailsjs/runtime";
import { useState } from "preact/hooks";
import { events, promptRequest, promptAction } from './events/events';
import CreateAccountFlow from './actions/create/createAccountFlow';

export function App() {
    const [eventData, setEventData] = useState<promptRequest>({} as promptRequest);
    const handlePromptRequest = (data: promptRequest) => {
        console.log("prompt request event received: ", data)
        setEventData(data)
    }

    EventsOn(events.promptRequest, handlePromptRequest)

    let flow;

    switch (eventData.Action) {
        case promptAction.newPasswordReq:
            flow = <CreateAccountFlow eventData={eventData}/>
            break;
    
        default:
            break;
    }

    return (
        <>
            <div id="App">
                {eventData.Action === promptAction.createAccountReq &&
                    <CreateAccountFlow eventData={eventData}/>
                }
                {eventData.Action === promptAction.exportReq &&
                    <ExportFlow eventData={eventData}/>
                }
            </div>
        </>
    )
}
