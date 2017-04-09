import { menu } from "./menu.js";
import { encrypt } from "./encrypt.js";
import { decrypt } from "./decrypt.js";

//component
var app = {
	view: function(vnode) {
		if (Object.keys(vnode.attrs.auth).length !== 0) {
			return m("div", [
				m(menu),
				m(eval(m.route.param("task")), {auth: vnode.attrs.auth, token: vnode.attrs.token})
			])
		} else {
			return m("h4", vnode.attrs.err)
		}
	}
}

export { app }
