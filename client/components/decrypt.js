import { Decrypt } from "../models/Decrypt.js";

//component
var decrypt = {
	view: function(vnode) {
		return m("div", [
			m("h3", "DECRYPT"),
			m("label[for=keys]", "Encryption Keys"),
			m("div[id=io]", [
				m("label[for=input]", "Input"),
				m("textarea[id=input]", {
					oninput: m.withAttr("value", function(value) {
						Decrypt.decrypt(vnode.attrs.token, value)
					})
				}),
				m("label[for=output]", "Output"),
				m("textarea[id=output]", {
					"value": Decrypt.output,
					"class": Decrypt.error ? "error" :  null
				})
			])
	
		])
	}
}

export { decrypt }
