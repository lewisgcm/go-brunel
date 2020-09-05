import { createStore, combineReducers, applyMiddleware, compose } from "redux";
import thunk from "redux-thunk";

import { reducer as layout } from "./modules/layout";

const composeEnhancers =
	((window as any).__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ as typeof compose) ||
	compose;

const reducer = combineReducers({
	layout,
});

export const store = createStore(
	reducer,
	composeEnhancers(applyMiddleware(thunk))
);
