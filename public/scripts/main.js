$(document).ready(function() {
    if (!window["WebSocket"]) {
        return;
    }

    var chat = $("#chat");
    var status = $("#status")
    var input = $("#input");
    var conn = new WebSocket('ws://' + window.location.host + '/ws');

    // Textarea is editable only when socket is opened.
    conn.onopen = function(e) {
        status.text("Connected. All systems Go.");
        status[0].className = "alert alert-success";
    };
    conn.onclose = function(e) {
        status.text("Connection closed.")
        status[0].className = "alert alert-danger";
    };

    // Whenever we receive a message, update textarea.
    conn.onmessage = function(e) {
        console.log("Receiving message");
        if (e.data != chat.text()) {
            // *** ALSO
            // Append random name or chat # assigned by backend
            chat.text(chat.text() + "\n" + e.data);
        }
    };

    var timeoutId = null;
    var typingTimeoutId = null;
    var isTyping = false;

    input.on("keydown", function() {
        isTyping = true;
        window.clearTimeout(typingTimeoutId);
    });

    input.on("keyup", function() {
        typingTimeoutId = window.setTimeout(function() {
            isTyping = false;
        }, 1000);

        window.clearTimeout(timeoutId);
        timeoutId = window.setTimeout(function() {
            if (isTyping) return;
            conn.send(input.text());
        }, 1100);
    });
})
