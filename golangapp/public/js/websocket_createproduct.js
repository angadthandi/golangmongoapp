function connWs(){
    var conn;

    conn = new WebSocket("ws://" + document.location.host + "/ws");
    conn.onclose = function (evt) {
    };
    conn.onopen = function (evt) {
        console.log("websocket onopen evt:", evt);
    };
    conn.onmessage = function (evt) {
        console.log("websocket onmessage evt:", evt);
    };

    return conn;
}

var ws = connWs();

function createNewProduct() {
    var messageObj = {
        api: "CreateProduct",
        message: {
            product: {
                productname: document.getElementById('productname').value,
                productcode: document.getElementById('productcode').value
            }
        }
    };

    console.log(messageObj);
    var jsonMsg = JSON.stringify(messageObj);
    console.log(jsonMsg);
    ws.send(jsonMsg);
}