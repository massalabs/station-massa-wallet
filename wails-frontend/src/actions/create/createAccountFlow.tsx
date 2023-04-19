// React component for the create account flow

import { useState } from "preact/hooks";
import { h } from 'preact';
import { promptRequest } from "../../events/events";
import PasswordPrompt from "../../components/passwordPrompt";
import Success from "../../components/Success";
import Failure from "../../components/Failure";

type Props = {
    eventData: promptRequest
};

export enum statuses {
    prompt = 0,
    success = 1,
    failure = 2,
}


const CreateAccountFlow = ({ eventData }: Props) => {
    // step state
    const [status, setStatus] = useState(statuses.prompt);

    switch (status) {
        case statuses.prompt:
            return (
                <PasswordPrompt eventData={eventData} />
            )
        case statuses.success:
            return (
                <Success msg={"The password has been created"} />
            )
        case statuses.failure:
            return (
                <Failure msg={"The password could not be created"} />
            )
        default:
            break;
    }

    return (
            <div>
                {eventData.Msg}
            </div>
    );
};

export default CreateAccountFlow;
