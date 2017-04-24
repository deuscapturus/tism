// BEGIN WORKAROUND
// modal code repeated here as workaround for firefox bug https://bugzilla.mozilla.org/show_bug.cgi?id=1358882
//import { modal } from "../components/modal.js";
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
// END WORKAROUND

import { Token } from "./Token.js";

//model
var Keys = {
	available: [],
	error: false,
	getKey: function(id) {
		return m.request({
			method: "POST",
			url: "key/get",
			data: { "token": Token.current, "Id": id },
		})
		.then(function(result){
			console.log(result.pubkey)
			m.render(document.getElementById("tokenModal"), m(modal, {"message": result.pubkey}))
		})
		.catch(function(e) {
			m.render(document.getElementById("tokenModal"), m(modal, {"message": e.message}))
		})
	},
	deleteKey: function(id) {
		return m.request({
			method: "POST",
			url: "key/delete",
			data: { "token": Token.current, "Id": id },
		})
		.then(function(result){
			Keys.getList()
		})
		.catch(function(e) {
			m.render(document.getElementById("tokenModal"), m(modal, {"message": e.message}))
		})
	},
	getList: function() {
		return m.request({
			method: "POST",
			url: "key/list",
			data: { "token": Token.current },
		})
		.then(function(result){
			Keys.available = result
			Keys.error = false
		})
		.catch(function(e) {
			console.log(e.message)
			Keys.error = true
		})
	}
}

export { Keys }
