import { Token } from "./Token.js";
import { modal } from "../components/modal.js";

//model
var CreateToken = {
    selectedKeys: [],
	toggle: function(id) {
		// Remove id from selectedKeys, otherwise add it
		if (CreateToken.selectedKeys.some(x => x === id)) {
			CreateToken.selectedKeys = CreateToken.selectedKeys.filter(x => x !== id)
		} else {
			CreateToken.selectedKeys.push(id)
		}
	},
	err: false,
	create: function() {
		return m.request({
			method: "POST",
			url: "token/new",
			data: { "token": Token.current, "keys": CreateToken.selectedKeys },
			deserialize: function(value) {return value}
		})
		.then(function(result){
			m.render(document.getElementById("tokenModal"), m(modal, {"message": result}))
			CreateToken.err = false
		})
		.catch(function(e) {
			console.log(e.message)
			CreateToken.err = true
		})
	}
}

export { CreateToken }
