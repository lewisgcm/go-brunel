import {throwError} from 'rxjs';

export function handleResponse(response: Response) {
	return response.ok ? response.json() : response.json().then((b) => throwError(b));
}
