//component
var menu = {
	view: function() {
		return m("ul", [
			m("li", m("a[href=/decrypt]", {oncreate: m.route.link}, "Decrypt")),
			m("li", m("a[href=/encrypt]", {oncreate: m.route.link}, "Encrypt"))
		])
	}
}

export { menu }
