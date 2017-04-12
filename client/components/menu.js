import { Token } from "../models/Token.js"

//component
var menu = {
	view: function() {
		return m("div", m("ul[class=nav nav-tabs]", [
			m("li", {class: (m.route.param("task") == "decrypt") ? "active" : null}, m("a[href=/decrypt]", {oncreate: m.route.link}, "Decrypt")),
			m("li", {class: (m.route.param("task") == "encrypt") ? "active" : null}, m("a[href=/encrypt]", {oncreate: m.route.link}, "Encrypt")),
			m("li", {class: (m.route.param("task") == "keys") ? "active" : null}, m("a[href=/keys]", {oncreate: m.route.link}, "Keys")),
	//		(Token.info.admin == 1) ? m("li[class=dropdown]",
	//			{class: (m.route.param("task") == "tokens"|| m.route.param("task") == "keys") ? "active" : null}, [
	//			m("a[href=#][class=dropdown-toggle][data-toggle=dropdown][role=button][aria-haspopup=true][aria-expanded=false]",
	//				"Admin ",
	//				m("span[class=caret]"),
	//			),
	//			m("ul[class=navbar-nav navbar-right]", m("ul[class=dropdown-menu]", [
	//				m("li", m("a[href=/tokens]", {oncreate: m.route.link}, "Tokens")),
	//				m("li", m("a[href=/keys]", {oncreate: m.route.link}, "Keys"))
	//			]))
	//		]) : null
		]))
	}
}

export { menu }
