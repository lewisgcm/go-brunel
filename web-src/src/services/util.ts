export function handleResponse(response: Response) {
	return response.ok
		? response.json()
		: response.json().then((b) => {
				if (b.Error) {
					throw new Error(b.Error);
				}
				throw new Error(response.statusText);
		  });
}
