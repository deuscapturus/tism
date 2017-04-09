//component
var menu = {
	view: function() {
		return m("ul", [
			m("li", m("a[href=/encrypt]", {oncreate: m.route.link}, "Encrypt")),
			m("li", m("a[href=/decrypt]", {oncreate: m.route.link}, "Decrypt")),
			m("li", m("a[href=/tokens]", {oncreate: m.route.link}, "Manage Tokens")),
			m("li", m("a[href=/keys]", {oncreate: m.route.link}, "Manage Keys"))
		])
	}
}

export { menu }
