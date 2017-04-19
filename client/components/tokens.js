import { Keys } from "../models/Keys.js";
import { CreateToken } from "../models/CreateToken.js";

//component
var tokens = {
	oninit: function(vnode) {
		Keys.getList()
	},
	view: function(vnode) {
		return m("div", [
			m("button[type=button][class=btn btn-default]", { onclick: function() { CreateToken.create() }}, m("span[class=glyphicon glyphicon-plus][aria-hidden=true]", " Create Token")),
			m("table[class=table table-striped table-hover]", [
				m("thead", [
					m("tr", [
						m("td", "ID"),
						m("td", "Name"),
						m("td", "Created")
					])
				]),
				m("tbody", 
					Keys.available.map( function(key) {
						return m("tr", {onclick: function() {
							CreateToken.toggle(key.Id)
							console.log(CreateToken.selectedKeys)
						}, class: (CreateToken.selectedKeys.some(x => x === key.Id)) ? "table-active" : null
						},[
							m("th[scope=row]", key.Id),
							m("td", key.Name),
							m("td", key.CreationTime)
						])
					})
				)
			])
		])
	}
}

export { tokens }
