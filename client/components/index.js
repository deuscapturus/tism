import { auth } from "./auth.js";
import { app } from "./app.js";
import { Token } from "../models/Token.js";

var index = {
	view: function() {
		return m("div[class=container]", [
			m("div[class=page-header", [
				m("h1", [
					"tISM ",
					m("small", "the Immutable Secrets Manager")
				])
			]),
			m(auth),
			m(app, {auth: Token.info, token: Token.current, err: Token.error}),
			m("div[id=app"),
		])
	}
}

export { index }
	
