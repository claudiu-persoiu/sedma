var communicationHandler = function (url, callback) {
    var socket, socketOpened = false;

    var connect = function () {
        socket = new WebSocket(url);

        socket.onopen = function () {
            socketOpened = true;
        };

        socket.onmessage = function (event) {
            callback(JSON.parse(event.data));
        };

        socket.onclose = function () {
            socketOpened = false;
            throw "Disconnected";
        };
    };

    return {
        send: function (action, data) {
            if (!data) {
                data = "";
            }

            if (socketOpened) {
                return socket.send(JSON.stringify({"action": action, "data": data}));
            }
            return false;
        },
        isConnected: function () {
            return socketOpened;
        },
        connect: function () {
            connect();
        }
    };
};