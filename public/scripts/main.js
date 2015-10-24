$(document).ready(function() {
    if (!window["WebSocket"]) {
        return;
    }

    var chat = $("#chat");
    var chatInput = $("#input");
    var conn = new WebSocket('ws://' + window.location.host + '/ws');

    // Textarea is editable only when socket is opened.
    conn.onopen = function(e) {
        chat.attr("disabled", false);
        chat.val("Connection received. Waiting for chats.")
    };
    conn.onclose = function(e) {
        chat.attr("disabled", true);
    };

    // Whenever we receive a message, update textarea.
    conn.onmessage = function(e) {
        if (e.data != chat.val()) {
            chat.val(e.data);
        }
    };

    var timeoutId = null;
    var typingTimeoutId = null;
    var isTyping = false;

    chatInput.on("keydown", function() {
        isTyping = true;
        window.clearTimeout(typingTimeoutId);
    });

    chatInput.on("keyup", function() {
        typingTimeoutId = window.setTimeout(function() {
            isTyping = false;
        }, 1000);

        window.clearTimeout(timeoutId);
        timeoutId = window.setTimeout(function() {
            if (isTyping) return;
            conn.send(chatInput.val());
        }, 1100);
    });
})
