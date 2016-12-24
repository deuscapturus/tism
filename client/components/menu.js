//component
var menu = {};

//view
menu.view = function(ctrl) {
	return m("ul[id=menu]", [
		m("li", m("a[href=/encrypt]", {config: m.route}, "Encrypt")),
		m("li", m("a[href=/decrypt]", {config: m.route}, "Decrypt")),
		m("li", m("a[href=/tokens]", {config: m.route}, "Manage Tokens")),
		m("li", m("a[href=/keys]", {config: m.route}, "Manage Keys"))
	
	])
};

//initialize
//m.mount(document.getElementById("menu"), menu);
