$(document).ready(function() {
    if (!window["WebSocket"]) {
        return;
    }

    var chat = $("#chat"),
        status = $("#status"),
        input = $("#input"),
        text = $("#text"),
        conn = new WebSocket('ws://' + window.location.host + '/ws');

    // Textarea is editable only when socket is opened.
    conn.onopen = function(e) {
        status.text("Connected. All systems Go.");
        status[0].className = "alert alert-success";
    };
    conn.onclose = function(e) {
        status.text("Connection closed.")
        status[0].className = "alert alert-danger";
    };

    // Whenever we receive a message, update the chat window.
    conn.onmessage = function(e) {
        if (e.data != chat.text()) {
            chat.text(chat.text() + "\n" + e.data);
        }
    };

    input.on('submit', function(e) {
        e.preventDefault();
        conn.send(text.val());
        text.val("");
    })
})
