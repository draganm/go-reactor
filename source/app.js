require('es5-shim');
var React = require('react');
var ReactDOM = require('react-dom');
var Bootstrap = require('react-bootstrap');
var _ = require("lodash");
var request = require('superagent');

var Dropzone = require('react-dropzone');

var EventEmitter = require('events');

var modelChangeEmitter = new EventEmitter();

var connectingModel = {
  title: "connecting ...",
  model: {
    el: "bs.PageHeader",
    ch: [{te: "Connecting to server ..."}]
  }
};

var connectedModel = {
  title: "connected!",
  model: {
    el: "bs.PageHeader",
    ch: [{ te: "Connected!" }]
  }
};

var disconnectedModel = {
  title: "disconnected!",
  model: {
    el: "bs.PageHeader",
    ch: [{ te: "Disconnected!"}]
  }
};

var url = window.location.href.replace(/#.*$/g,"")

url = url.replace(/^http/,"ws")+"ws";

var ws = null;

var connectSocket = function() {
  ws = new WebSocket(url);

  var oldHashChange = window.onhashchange;

  ws.onopen = function(){
    modelChangeEmitter.emit('modelChange', connectedModel);
    window.onhashchange = function(event) {
      ws.send(JSON.stringify({id: "main_window", type: "popstate", value: window.location.hash}));
    };
    ws.send(JSON.stringify({id: "main_window", type: "popstate", value: window.location.hash}));

  }

  ws.onclose = function(){
    window.onhashchange = oldHashChange;
    ws=null
    setTimeout(connectSocket,1000)
  	modelChangeEmitter.emit('modelChange', disconnectedModel);
  }

  ws.onmessage = function(evt){
    var pageModel = JSON.parse(evt.data);
    modelChangeEmitter.emit('modelChange', pageModel);
  }
}

connectSocket();

var Wrapper = React.createClass({
  render: function() {

    var model = this.props.model


    var elementName = model.el
    var parts = elementName.split(".")

    var elementClass = elementName
    var attrs = _.clone(model.at || {})

    if (parts[0] === "bs" && parts.length == 2) {
      elementClass = Bootstrap[parts[1]]
    }

    if (parts[0] === "bs" && parts.length == 3) {
      elementClass = Bootstrap[parts[1]][parts[2]]
    }

    if (parts[0] == "Dropzone") {
      elementClass=Dropzone
      attrs.onDrop=function(files){
        var req = request.post(attrs["target"]);
        files.forEach((file)=> {
            req.attach(file.name, file);
        });
        req.end();
      }
    }

    var children = (model.ch || []).map(function(child) {
      if (child.te || child.te == "") {
        return child.te;
      } else {
        return React.createElement(Wrapper, {model: child})
      }
    });


    var eventHandlers = _.map(model.ev || [], function(eventData){

      var preventDefault = eventData.pd

      var stopPropagation = eventData.sp

      var eventName = eventData.name

      var extraValues = eventData.xv || []

      var onName = "on" + eventName.charAt(0).toUpperCase() + eventName.slice(1);
      attrs[onName] = function(evt) {
        if (preventDefault) {
          evt.preventDefault();
        }
        if (stopPropagation) {
          evt.stopPropagation();
          evt.nativeEvent.stopImmediatePropagation();
        }

        if (evt) {

          var xv = _.reduce(extraValues, function(x, n){
            x[n] = evt[n]
            return x
          }, {})

          var files = (evt.target.files || []);

          if (files.length > 0) {

            var reader = new FileReader();
            reader.onload = function(theFile) {
              if (ws) {
                ws.send(JSON.stringify({id: model.id, type: eventName, value: files[0].name, data: reader.result, xv: xv}))
              }
            };

            reader.readAsText(files[0]);

          } else {
            if (ws) {
              ws.send(JSON.stringify({id: model.id, type: eventName, value: evt.target.value, xv: xv}))
            }
          }

        } else {
          if (ws) {
            ws.send(JSON.stringify({id: model.id, type: eventName}))
          }
        }

      }
    });

    var params = [elementClass, attrs].concat(children);

    return React.createElement.apply(this, params)
  }
});

var Page = React.createClass({
  getInitialState: function() {
    return connectingModel;
  },
  componentDidMount: function() {
    var that = this;
    modelChangeEmitter.on('modelChange', function(update) {
      if (update.eval) {
        eval(update.eval)
      }
      if (update.title) {
        document.title = update.title
      }
      if (update.location) {
        window.location = update.location
      }
      if (update.model) {
        that.setState({model: update.model});
      }
    });
  },
  render: function() {
    return <Wrapper model={this.state.model}/>
  }
});



ReactDOM.render(<Page/>, document.getElementById('react-application'));
