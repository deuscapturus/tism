import { CreateKey } from "../models/CreateKey.js";

//component
var createKey = {
  view: function(vnode) {
    return m(
      "form",
      {
        onsubmit: function() {
          CreateKey.create();
        }
      },
      [
        m("div[class=form-group", [
          m("label[class=control-label]", "Name"),
          m(
            "div",
            m("input[class=form-control][type=text]", {
              oninput: m.withAttr("value", function(value) {
                CreateKey.name = value;
              })
            })
          )
        ]),
        m("div[class=form-group", [
          m("label[class=control-label]", "Comment"),
          m(
            "div",
            m("input[class=form-control][type=text]", {
              oninput: m.withAttr("value", function(value) {
                CreateKey.comment = value;
              })
            })
          )
        ]),
        m("div[class=form-group", [
          m("label[class=control-label]", "Email"),
          m(
            "div",
            m("input[class=form-control][type=email]", {
              oninput: m.withAttr("value", function(value) {
                CreateKey.email = value;
              })
            })
          )
        ]),
        m("button[type=submit][class=btn btn-default]", "Create New Key")
      ]
    );
  }
};

export { createKey };
