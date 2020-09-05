import { Observable } from "rxjs";
import { webSocket, WebSocketSubject } from "rxjs/webSocket";
import { injectable } from "inversify";

import { EventTypes, EventType } from "./models";
import { AuthService } from "./authService";

type ExtractActionParameters<A, T> = A extends { Type: T } ? A : never;

@injectable()
export class SocketService {
	private _subject: WebSocketSubject<EventTypes>;

	constructor(private authService: AuthService) {
		const protocolPrefix =
			window.location.protocol === "https:" ? "wss:" : "ws:";

		const host =
			process.env.REACT_APP_WEBSOCKET_HOST || window.location.host;

		this._subject = webSocket(
			`${protocolPrefix}//${host}/api/event?token=${this.authService.getToken()}`
		);
	}

	events<T extends EventType>(
		type: T
	): Observable<ExtractActionParameters<EventTypes, T>> {
		return this._subject.multiplex(
			() => ({ subscribe: type }),
			() => ({ unsubscribe: type }),
			(message) => {
				return message.Type === type;
			}
		);
	}
}
