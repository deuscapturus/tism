import { Token } from "./Token.js";

//model
var Encrypt = {
  error: false,
  input: "",
  output: "",
  selectedKey: "",
  encrypt: function() {
    return m
      .request({
        method: "POST",
        url: "encrypt",
        data: {
          token: Token.current,
          decsecret: Encrypt.input,
          id: Encrypt.selectedKey
        },
        deserialize: function(value) {
          return value;
        }
      })
      .then(function(result) {
        Encrypt.output = result;
        Encrypt.error = false;
      })
      .catch(function(e) {
        Encrypt.output = e.message;
        Encrypt.error = true;
      });
  }
};

export { Encrypt };
