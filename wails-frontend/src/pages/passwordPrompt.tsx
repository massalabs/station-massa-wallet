import { useState } from "preact/hooks";
import { AbortAction, ApplyPassword, Hide } from "../../wailsjs/go/walletapp/WalletApp";
import { EventsOnce, WindowReloadApp } from "../../wailsjs/runtime";
import { h } from 'preact';
import { events, promptRequest } from "../events/events";

type Props = {
    eventData: promptRequest
};


const PasswordPrompt = ({ eventData }: Props) => {

    const [applyResultMsg, setPasswordResult] = useState("do it!");
    const [password, setPassword] = useState("");

    const updatePasswordResult = (result: string) => setPasswordResult(result);

    const hideAndReload = () => {
        Hide();
        // Reload the wails frontend
        WindowReloadApp();
    };

    const handleCancel = () => {
        AbortAction();
        hideAndReload();
    }

    const handleApplyResult = (result: any) => {
        console.log("result", result);
        if (result.Success) {
            updatePasswordResult("Password applied successfully!");
            setTimeout(hideAndReload, 2000);
        } else {
            updatePasswordResult(result.Data);
        }
    };

    const applyPassword = () => {
        if (password.length) {
            EventsOnce(events.passwordResult, handleApplyResult);
            ApplyPassword(password);
        }
    };

    const updatePassword = (e: any) => setPassword(e.target.value);

    return (
        <div id="App">
            <div>
                {eventData.Msg}
            </div>
            <div id="result" className="result">
                Enter your password below to validate 👇
            </div>
            <div id="input" className="input-box">
                <input
                    id="name"
                    className="input"
                    onChange={updatePassword}
                    autoComplete="off"
                    name="input"
                    type="password"
                />
                <button className="btn" onClick={applyPassword}>
                    Ok
                </button>
                <button className="btn" onClick={handleCancel}>
                    Cancel
                </button>
            </div>
            <div id="result" className="result">
                {applyResultMsg}
            </div>
        </div>
    );
};

export default PasswordPrompt;