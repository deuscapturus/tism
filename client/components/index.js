import { auth } from "./auth.js";
import { app } from "./app.js";
import { Token } from "../models/Token.js";

var index = {
	view: function() {
		return m("div", [
			m(auth),
			m(app, {auth: Token.info, token: Token.current, err: Token.error}),
			m("div[id=app"),
		])
	}
}

export { index }
	
