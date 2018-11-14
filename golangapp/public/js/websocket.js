window.onload = function () {
    var conn;
    var msg = document.getElementById("msg");
    var log = document.getElementById("log");
    document.getElementById("getProductsBtn").onclick = function () {
        console.log("getProductsBtn");
        var messageObj = {
            api: "GetProducts",
            message: {
                data: null
            }
        };

        console.log(messageObj);
        var jsonMsg = JSON.stringify(messageObj);
        console.log(jsonMsg);
        conn.send(jsonMsg);
    };
    function appendLog(item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }
    document.getElementById("form").onsubmit = function () {
        if (!conn) {
            return false;
        }
        if (!msg.value) {
            return false;
        }

        var msgObj = {
            // api: "GetProducts",
            api: "Default",
            message: {
                data: msg.value
            }
        };
        // conn.send(msg.value);
        console.log(msgObj);
        var jsonMsg = JSON.stringify(msgObj);
        console.log(jsonMsg);
        conn.send(jsonMsg);
        msg.value = "";
        return false;
    };
    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws");
        conn.onclose = function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);
        };
        conn.onopen = function (evt) {
            console.log("websocket onopen evt:", evt);
        };
        conn.onmessage = function (evt) {
            console.log("websocket onmessage evt:", evt);
            var messages = evt.data.split('\n');
            for (var i = 0; i < messages.length; i++) {
                var item = document.createElement("div");
                item.innerText = messages[i];
                appendLog(item);
            }
        };
    } else {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendLog(item);
    }
};