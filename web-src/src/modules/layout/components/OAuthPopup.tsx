import React, {useEffect, Fragment} from 'react';

interface Props {
    isOpen: boolean;
    provider: string;
    onDone: (token: string) => void;
}

export function OAuthPopup({isOpen, provider, onDone}: Props) {
	useEffect(
		() => {
			const url = `http://localhost:8085/api/user/login?provider=${provider}`;

			if (isOpen) {
				const popup = window.open(url, 'Login');

				if (popup) {
					const listener = window.onmessage = (e: MessageEvent) => {
						if (e.data.token) {
							popup.close();
							onDone(e.data.token);
							window.removeEventListener('message', listener);
						}
					};

					return () => {
						popup.close();
					};
				}
			}
		},
		[isOpen, provider, onDone],
	);

	return <Fragment />;
}
