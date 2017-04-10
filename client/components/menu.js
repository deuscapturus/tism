//component
var menu = {
	view: function() {
		return m("ul[class=nav nav-tabs]", [
			m("li", {class: (m.route.param("task") == "decrypt") ? "active" : null}, m("a[href=/decrypt]", {oncreate: m.route.link}, "Decrypt")),
			m("li", {class: (m.route.param("task") == "encrypt") ? "active" : null}, m("a[href=/encrypt]", {oncreate: m.route.link}, "Encrypt"))
		])
	}
}

export { menu }
