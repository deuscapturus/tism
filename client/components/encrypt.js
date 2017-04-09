import { Encrypt } from "../models/Encrypt.js";
import { Keys } from "../models/Keys.js";

//component
var encrypt = {
	oninit: function(vnode) {
		Keys.getList(vnode.attrs.token)
	},
	view: function(vnode) {
		return m("div", [
            m("h3", "ENCRYPT"),
			Keys.available.map( function(key) {
				return (key.Id != "ALL") ? m("label[class=keys]", {"for": key.Id}, [key.Name, 
					m("input[type=radio][name=key][class=keys]", {
						"value": key.Id,
						"id": key.Id,
						oninput: m.withAttr("value", function(value) {
							Encrypt.selectedKey = value
						})
					})
				]) : null
			}),
			m("div[id=io]", [
				m("label[for=input]", "Input"),
				m("textarea[id=input]", {
					oninput: m.withAttr("value", function(value) {
						Encrypt.encrypt(vnode.attrs.token, value, Encrypt.selectedKey)
					})
				}),
				m("label[for=output]", "Output"),
				m("textarea[id=output]", {
					"value": Encrypt.output,
					"class": Encrypt.error ? "error" :  null
				})
			])
		])
	}
}

export { encrypt }
