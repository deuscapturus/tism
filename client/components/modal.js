
//component
var modal = {
	view: function(vnode) {
		return m("div[id=insideModal]", m("p", vnode.attrs.message), m("button", { onclick: function() {
			m.render(document.getElementById("insideModal"), null)
		}}, "Close"))
	}
}

export { modal }
