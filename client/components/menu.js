//component
var menu = {};

//view
menu.view = function(ctrl) {
	return m("ul[id=menu]", [
		m("li", m("a[href=#encrypt]", "Encrypt")),
		m("li", m("a[href=#decrypt]", "Decrypt")),
		m("li", m("a[href=#tokens]", "Manage Tokens")),
		m("li", m("a[href=#keys]", "Manage Keys"))
	
	])
};

//initialize
m.mount(document.getElementById("menu"), menu);
