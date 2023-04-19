// Success component

import { AbortAction, Hide } from "../../wailsjs/go/walletapp/WalletApp";
import { WindowReloadApp } from "../../wailsjs/runtime/runtime";

import { h } from 'preact';

const Failure = ({ msg }: { msg: string }) => {
	const handleClose = () => {
		AbortAction();
		Hide();
		WindowReloadApp();
	}

	const handleGoBack = () => {
		console.log("go back not implemented yet");
	}

	return (
			<div>
					{/* logo */}
					<p>FAILURE LOGO</p>

					{/* message */}
					{msg}

					{/* close button */}
					<button onClick={handleClose}>Close</button>

					{/* go back button */}
					<button onClick={handleGoBack}>Go back</button>
			</div>
	)
}

export default Failure;