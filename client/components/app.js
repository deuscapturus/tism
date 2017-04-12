import { menu } from "./menu.js";
import { encrypt } from "./encrypt.js";
import { decrypt } from "./decrypt.js";
import { keys } from "./keys.js";
import { tokens } from "./tokens.js";

//component
var app = {
	view: function(vnode) {
		if (Object.keys(vnode.attrs.auth).length !== 0) {
			return m("div", [
				m(menu),
				m(eval(m.route.param("task")), {auth: vnode.attrs.auth, token: vnode.attrs.token})
			])
		} else {
			if (vnode.attrs.err) { return m("span[class=alert alert-danger center-block]", vnode.attrs.err) } else { return null }
		}
	}
}

export { app }
