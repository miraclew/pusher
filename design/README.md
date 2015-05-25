#design


## User Message Queue Design

            L(head)     R(tail)
receiver <- [ 1, 2, 3, 4, 5 ] <- sender

## Send message

1. Push to queue(rpush)
2. Client offline, return (iOS push if need)
3. Otherwise, trigger process queue

## Accept connection

1. Auth client
2. Trigger process queue

## Handle ack
1. Remove message from queue (lrem)
2. Check queue length, return if it's 0
3. Schedule next trigger process queue, cancel exist schedule

## Process queue:

1. Get first 10 elements of queue, from left
2. Write messages to connection
