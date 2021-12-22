package main

import (
	"GOThinkUtils/thinkutils"
	"GOThinkUtils/thinkutils/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"html/template"
	"net/http"
	"strings"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func onConnect(pConn *websocket.Conn) {
	log.Info("New Connect")
}

func onClose(pConn *websocket.Conn) {
	log.Info("Conn closed")
}

func onMessage(pConn *websocket.Conn, msg []byte) {
	log.Info("[%p] recv: %s", pConn, thinkutils.StringUtils.BytesToString(msg))
	err := pConn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Info("write:", err.Error())
	}
}

func homeHandler(c *gin.Context) {
	homeTemplate.Execute(c.Writer, "ws://"+c.Request.Host+"/echo")
}

func ipHandler(c *gin.Context) {
	szHost := strings.Split(c.Request.Host, ":")[0]
	c.JSON(http.StatusOK, thinkutils.AjaxResultSuccessWithData(szHost))
}

func main() {
	pHandler := &thinkutils.WSHandler{
		OnConnect: onConnect,
		OnMsg:     onMessage,
		OnClose:   onClose,
	}

	eng := gin.Default()
	// 路由组1 ，处理GET请求
	eng.GET("/echo", pHandler.Handler)
	eng.GET("/", homeHandler)
	eng.GET("/ip", ipHandler)

	eng.Run(":8080")
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
