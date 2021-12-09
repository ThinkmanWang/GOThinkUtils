package main

import (
	"GOThinkUtils/thinkutils"
	"GOThinkUtils/thinkutils/logger"
	"github.com/gorilla/websocket"
	"html/template"
	"net/http"
)

var (
	log      *logger.LocalLogger = logger.DefaultLogger()
	upgrader                     = websocket.Upgrader{}
)

type OnConnectCallback func(pConn *websocket.Conn)
type OnCloseCallback func(pConn *websocket.Conn)
type OnMsgCallback func(pConn *websocket.Conn, msg []byte)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Info("upgrade:", err)
		return
	}

	defer func() {
		onClose(c)
		c.Close()
	}()

	go onConnect(c)

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}

		switch mt {
		case websocket.BinaryMessage, websocket.TextMessage:
			go onMessage(c, message)
		default:
			continue
		}

	}
}

func onConnect(pConn *websocket.Conn) {
	log.Info("New Connect")
}

func onClose(pConn *websocket.Conn) {
	log.Info("Conn closed")
}

func onMessage(pConn *websocket.Conn, msg []byte) {
	log.Info("recv: %s", thinkutils.StringUtils.BytesToString(msg))
	err := pConn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Info("write:", err.Error())
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
	http.HandleFunc("/echo", wsHandler)
	http.HandleFunc("/", homeHandler)
	http.ListenAndServe("127.0.0.1:8080", nil)
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))
