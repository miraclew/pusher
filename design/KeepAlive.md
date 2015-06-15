# KeepAlive

## Client side ping

 1. After connection established, start a timer to send interval ping message
 2. Create a counter to record missed pong, increase the counter when send a ping 
 3. Reset the counter to 0 if receive a pong
 4. If the counter > MAX_MISSED_PONG, close connection and reconnect


Code
`

    const PING_MSG = 'p';
    const PONG_MSG = 'q';
    const MAX_MISSED_PONG = 3;

    var missed_pong = 0;

    const PING_INTERVAL = 5000; // 5 seconds
    var ping_timer = null;

    const RECONNECT_INTERVAL = 3000; // 3 seconds
    var reconnect_timer = null;
    
    function on_connected(conn) {
        // setup timer
        ping_timer = setInterval(PING_INTERVAL, function() {
            if(missed_pong > MAX_MISSED_PONG) {
                conn.close();
                schedule_reconnect();
            }
            
            conn.send(PING_MSG);
            missed_pong++;
        });    
        
        if(reconnect_timer != null) {
            clearInterval(reconnect_timer);
        }
    }
    
    function on_message_received(conn, $msg) {
        if($msg == PONG_MSG) {
            missed_pong = 0;            
        } else {
            // process normal payload
        }
    }
    
    function on_closed() {
        if(ping_timer != null) {
            clearInterval(ping_timer)
        }
    }
    
    function schedule_reconnect() {
        // reconnect
        reconnect_timer = setInterval(RECONNECT_INTERVAL, function() {
            connect();    
        });
    }
    
    function connect() {
        // connect to ws server
    }
`


