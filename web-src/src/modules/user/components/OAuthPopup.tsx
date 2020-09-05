import React, { useEffect, Fragment } from "react";

interface Props {
	isOpen: boolean;
	provider: string;
	onDone: (token: string) => void;
	onError: (error: string) => void;
	onAbort: () => void;
}

function OAuthPopupComponent({
	isOpen,
	provider,
	onDone,
	onError,
	onAbort,
}: Props) {
	useEffect(() => {
		const url = `${process.env.REACT_APP_OAUTH_BASE_URL}/api/user/login?provider=${provider}`;

		if (isOpen) {
			const popup = window.open(url, "Login");
			let messageReceived = false;

			if (popup) {
				const listener = (window.onmessage = (e: MessageEvent) => {
					if (e.data.token && !messageReceived === true) {
						window.removeEventListener("message", listener);
						onDone(e.data.token);
						messageReceived = true;
					} else if (e.data.error && !messageReceived === true) {
						window.removeEventListener("message", listener);
						onError(e.data.error);
						messageReceived = true;
					}
				});

				const interval = setInterval(() => {
					if (popup.closed) {
						clearInterval(interval);
						onAbort();
					}
				}, 5);

				return () => {
					clearInterval(interval);
					popup.close();
				};
			}
		}
	}, [isOpen, provider, onDone, onError, onAbort]);

	return <Fragment />;
}

export const OAuthPopup = React.memo(
	OAuthPopupComponent,
	(a, b) => a.isOpen === b.isOpen
) as typeof OAuthPopupComponent;
