import { Decrypt } from "../models/Decrypt.js";

//component
var decrypt = {
	view: function(vnode) {
		return m("div", [
			m("div[id=io][class=form-group]", [
				m("label[for=input][class=control-label]", "Input"),
				m("textarea[id=input][class=form-control][rows=4][style=word-break: break-all;]", {
					oninput: m.withAttr("value", function(value) {
						Decrypt.decrypt(vnode.attrs.token, value)
					})
				}),
				m("label[for=output][class=control-label]", "Output"),
				m("textarea[id=output][class=form-control][rows=2][readonly]", {
					"value": Decrypt.output,
					"class": Decrypt.error ? "error" :  null
				})
			])
	
		])
	}
}

export { decrypt }
