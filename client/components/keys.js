import { keysList } from "../components/keysList.js";
import { createKey } from "../components/createKey.js";

//component
var keys = {
	view: function(vnode) {
		return m("div", [
			m("div[class=row]", [
				m("div[class=col-xs-12 col-md-3]", m(createKey)),
				m("div[class=col-xs-12 col-md-9]", m(keysList))
			])
		])
	}
}

export { keys }
