import { useState } from "preact/hooks";
import { AbortAction, ApplyPassword, Hide } from "../../wailsjs/go/walletapp/WalletApp";
import { EventsOnce, WindowReloadApp } from "../../wailsjs/runtime/runtime";
import { h } from 'preact';
import { events, promptAction, promptRequest } from "../events/events";
import { useLocation, useNavigate } from "react-router-dom";

const PasswordPrompt = () => {

    const navigate = useNavigate();

    const { state } = useLocation();
    const req: promptRequest = state.req;

    const applyStr = (req: promptRequest) => {
        switch (req.Action) {
            case promptAction.deleteReq:
                return "Delete";
            default: return "Apply";
        }
    }

    const [applyResultMsg, setPasswordResult] = useState("do it!");
    const [password, setPassword] = useState("");

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
            navigate("/success", { state: { msg: "Password applied successfully!" } })
            setTimeout(hideAndReload, 2000);
        } else {
            setPasswordResult(result.Data);
            if (result.Error === "timeoutError") {
                setTimeout(hideAndReload, 2000);
            }
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
        <section class="PasswordPrompt">
            <div>
                {req.Msg}
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
            </div>
            <div>
                <button className="btn" onClick={handleCancel}>
                    Cancel
                </button>
                <button className="btn" onClick={applyPassword}>
                    {applyStr(req)}
                </button>
            </div>
            <div id="result" className="result">
                {applyResultMsg}
            </div>
        </section>
    );
};

export default PasswordPrompt;