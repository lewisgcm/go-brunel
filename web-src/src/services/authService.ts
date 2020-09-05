import { injectable } from "inversify";
import jwtDecode from "jwt-decode";
import moment from "moment";

import { UserRole } from "./models";

interface JWT {
	exp: string;
	role: UserRole;
}

@injectable()
export class AuthService {
	isAuthenticated(): boolean {
		const token = window.localStorage.getItem("jwt");
		if (token) {
			try {
				return moment
					.unix(Number((jwtDecode(token) as JWT).exp))
					.isAfter(moment());
			} catch (e) {
				return false;
			}
		}
		return false;
	}

	getRole(): UserRole | undefined {
		const token = window.localStorage.getItem("jwt");
		if (token) {
			return (jwtDecode(token) as JWT).role;
		}
		return undefined;
	}

	setAuthentication(token: string) {
		window.localStorage.setItem("jwt", token);
	}

	getToken(): string | null {
		return window.localStorage.getItem("jwt");
	}

	getAuthHeaders(): HeadersInit {
		return {
			Authorization: `Bearer ${window.localStorage.getItem("jwt")}`,
		} as HeadersInit;
	}
}
