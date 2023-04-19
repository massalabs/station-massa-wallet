// Success component

import { AbortAction, Hide } from "../../wailsjs/go/walletapp/WalletApp";
import { WindowReloadApp } from "../../wailsjs/runtime/runtime";

import { h } from 'preact';

const Success = ({ msg }: { msg: string }) => {
	const handleClose = () => {
		AbortAction();
		Hide();
		WindowReloadApp();
	}

		return (
				<div>
						{/* logo */}
						<p>SUCCESS LOGO</p>

						{/* message */}
						{msg}

						{/* close button */}
						<button onClick={handleClose}>Close</button>
				</div>
		)
}

export default Success;