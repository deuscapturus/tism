
//component
var modal = {
	view: function(vnode) {
		return m("div", m("p", vnode.attrs.message), m("button", { onclick: function() {
			m.render(vnode.dom.parentNode, null)
		}}, "Close"))
	}
}

export { modal }
