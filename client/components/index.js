import { auth } from "./auth.js";
import { menu } from "./menu.js";
import { app } from "./app.js";
import { Token } from "../models/Token.js";

var index = {
	view: function() {
		return m("div", [
			m(auth),
			m(menu),
			m(app, {auth: Token.info})
		])
	}
}

export { index }
	
