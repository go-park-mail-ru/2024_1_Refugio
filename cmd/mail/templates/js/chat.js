$(function(){

    var socket = null;
    var msgBox = $("#chatbox textarea");
    var messages = $("#messages");

    $("#chatbox").submit(function(){

        if (!msgBox.val()) return false;
        if (!socket) {
            alert("Error: There is no socket connection.");
            return false;
        }

        socket.send(msgBox.val());
        msgBox.val("");
        return false;

    });

    if (!window["WebSocket"]) {
        alert("Error: Your browser does not support web sockets.")
    } else {
        const href = window.location.href;
        const hrefArr = href.split('/');
        const mail = hrefArr[hrefArr.length - 1];
        //let variable = "ivan@mailhub.su";
        //`web/something/${variable}`
        socket = new WebSocket(`ws://localhost:8080/api/v1/auth/web/websocket_connection/${mail}`);
        //socket = new WebSocket("ws://localhost:8080/api/v1/auth/web/websocket_connection");
        socket.onclose = function() {
            alert("Connection has been closed.");
        }
        socket.onmessage = function(e) {
            messages.append($("<li>").text(e.data));
        }
    }

});
