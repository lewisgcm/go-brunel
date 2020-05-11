import React, {useEffect, Fragment} from 'react';

interface Props {
	isOpen: boolean;
	provider: string;
	onDone: (token: string) => void;
	onError: (error: string) => void;
}

export function OAuthPopup({isOpen, provider, onDone, onError}: Props) {
	useEffect(
		() => {
			const url = `${process.env.REACT_APP_OAUTH_BASE_URL}/api/user/login?provider=${provider}`;

			if (isOpen) {
				const popup = window.open(url, 'Login');

				if (popup) {
					const listener = window.onmessage = (e: MessageEvent) => {
						if (e.data.token) {
							popup.close();
							onDone(e.data.token);
							window.removeEventListener('message', listener);
						} else if (e.data.error) {
							popup.close();
							onError(e.data.error);
							window.removeEventListener('message', listener);
						}
					};

					return () => {
						popup.close();
					};
				}
			}
		},
		[isOpen, provider, onDone, onError],
	);

	return <Fragment />;
}
