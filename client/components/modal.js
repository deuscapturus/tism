
//component
var modal = {
	oncreate: function(vnode) {
		vnode.dom.parentNode.style.display = "block"
	},
	onremove: function(vnode) {
		vnode.dom.parentNode.style.display = "none"
	},
	view: function(vnode) {
		return m("div[class=modal-dialog modal-lg]", m("div[class=modal-content]", [
			m("div[class=modal-header]", m("h5[class=modal-title]", "New Token" )), 
			m("div[class=modal-body]", [
				m("p[style=word-wrap: break-word;]", vnode.attrs.message),
				m("button[class=close]", { onclick: function() {
					m.render(vnode.dom.parentNode, null)
				}}, "Close")]
			)]))
	}
}

export { modal }
