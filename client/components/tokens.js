import { Keys } from "../models/Keys.js";
import { CreateToken } from "../models/CreateToken.js";

//component
var tokens = {
  oninit: function(vnode) {
    Keys.getList();
  },
  view: function(vnode) {
    return m("form[class=form-inline]", [
      m(
        "div[class=form-group]",
        m(
          "button[type=button][class=btn btn-default]",
          {
            onclick: function() {
              CreateToken.create();
            }
          },
          m(
            "span[class=glyphicon glyphicon-plus][aria-hidden=true]",
            " Create Token"
          )
        ),
        m(
          "div[class=checkbox]",
          m(
            "label",
            m("input[class=checkbox][type=checkbox]", {
              onclick: m.withAttr("checked", function(value) {
                CreateToken.toggleAdmin(value);
              })
            }),
            "Make Admin"
          )
        )
      ),
      m("table[class=table table-striped table-hover]", [
        m("thead", [
          m("tr", [m("td", "ID"), m("td", "Name"), m("td", "Created")])
        ]),
        m(
          "tbody",
          Keys.available.map(function(key) {
            return m(
              "tr",
              {
                onclick: function() {
                  CreateToken.toggle(key.Id);
                  console.log(CreateToken.selectedKeys);
                },
                class: CreateToken.selectedKeys.some(x => x === key.Id)
                  ? "table-active"
                  : null
              },
              [
                m("th[scope=row]", key.Id),
                m("td", key.Name),
                m("td", key.CreationTime)
              ]
            );
          })
        )
      ])
    ]);
  }
};

export { tokens };
