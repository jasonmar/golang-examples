# Example 7 - High Throughput JSON Ingestion

This example go application implements a simple REST API which accepts a JSON payload with multiple entries.

The design is based on a blog by [marcio.io](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/), who reports handling over 1M requests per minute on 4 c4.large AWS EC2 instances.

Decoding a massive volume of JSON and asynchronously performing tasks is a common use case for Go. This example application provides a fully working demonstration of how to accomplish this.


# Concurrency

### goroutines
At initialization, a Dispatcher is created and a goroutine is created for creating workers.
Each worker's Work() method is called in a goroutine, beginning the task loop.
Pointer copy from a global task channel to an individual worker channel occurs in a goroutine.
A goroutine is created for each HTTP Request. This is where JSON decoding occurs.

### channels
The Dispatcher makes a channel on which workers can send their own channels.
At the start of the task loop, workers send their task channel then block until a pointer to a task is received.
Each worker also has a quit channel which causes the task loop to return if a bool is received.
There is a single global task channel which is read only by the Dispatcher.
The global task channel is buffered, with size set by environment variable read at startup.
Each HTTP request sends a Task to the global task channel after successfully decoding a JSON payload.
The Dispatch loop blocks until a task is received from the global task channel then spawns a goroutine which blocks until a worker task channel is ready to receive.

### Scalability

HTTP requests are completed as soon as all tasks are sent to the global task channel, with task completion handled asynchronously. This reduces latency of the HTTP requests which should in theory enable higher throughput.
Most of the time workers will consume tasks fast enough to prevent the request handler from blocking until the global task channel is ready to receive.
If the buffer capacity of the global task channel is exhausted, HTTP request handlers will block and clients will experience increased latency. This is intentional and is meant to prevent new requests from impacting in-progress tasks.
Request latency can monitored and used to trigger autoscaling. As new instances come on line workers should be able to consume fast enough to prevent the buffer from being filled and latency should return to normal.



# Use of Pointers

### Can you spot the error in the code below?

```
for i, p := range d.Data {
    queue <- &Task{Time: t, Parcel: &p}
}
```

### What's wrong:

`for i, p := range d.Data` allocates memory for a single instance p and overwrites it each iteration.
Thus if you take a pointer to p on separate iteration, all those pointers will point to the same address.

```
for i, p := range d.Data {
    queue <- &Task{Time: t, Parcel: &p} // All tasks will end up with a pointer to the last Parcel
}
```

### Correct implementation:

```
for i := range d.Data {
    p := &d.Data[i] // get a pointer to an indexed element
    queue <- &Task{Time: t, Parcel: p} // each task now references a distinct Parcel
}
```


# Testing

### Test Query
```
curl -H "Content-Type: application/json" -X POST \
-d '{"Data":[
    {"id":"1","time":"2017-01-01T00:00:01Z","data":"abc"},
    {"id":"2","time":"2017-01-01T00:00:02Z","data":"def"},
    {"id":"3","time":"2017-01-01T00:00:03Z","data":"ghi"},
    {"id":"4","time":"2017-01-01T00:00:04Z","data":"jkl"},
    {"id":"5","time":"2017-01-01T00:00:05Z","data":"mno"},
    {"id":"6","time":"2017-01-01T00:00:06Z","data":"pqr"},
    {"id":"7","time":"2017-01-01T00:00:07Z","data":"stu"},
    {"id":"8","time":"2017-01-01T00:00:08Z","data":"vwx"},
    {"id":"9","time":"2017-01-01T00:00:09Z","data":"yz0"}
]}' http://127.0.0.1:8080
```

### Test Results

```
2017/01/03 22:35:46 Handler: Received Parcel 1
2017/01/03 22:35:46 Handler: Received Parcel 2
2017/01/03 22:35:46 Handler: Received Parcel 3
2017/01/03 22:35:46 Handler: Received Parcel 4
2017/01/03 22:35:46 Worker 1: Received Task with Parcel 3
2017/01/03 22:35:46 Handler: Received Parcel 5
2017/01/03 22:35:46 Worker 2: Received Task with Parcel 4
2017/01/03 22:35:46 Worker 0: Received Task with Parcel 1
2017/01/03 22:35:46 Worker 3: Received Task with Parcel 2
2017/01/03 22:35:46 Handler: Received Parcel 6
2017/01/03 22:35:46 Handler: Received Parcel 7
2017/01/03 22:35:46 Worker 5: Received Task with Parcel 6
2017/01/03 22:35:46 Worker 4: Received Task with Parcel 5
2017/01/03 22:35:46 Handler: Received Parcel 8
2017/01/03 22:35:46 Handler: Received Parcel 9
2017/01/03 22:35:46 Worker 6: Received Task with Parcel 8
2017/01/03 22:35:46 Worker 7: Received Task with Parcel 7
2017/01/03 22:35:46 Task: Completed for Parcel 2 (22.0749ms)
2017/01/03 22:35:46 Worker 3: Received Task with Parcel 9
2017/01/03 22:35:46 Task: Completed for Parcel 5 (31.6356ms)
2017/01/03 22:35:46 Task: Completed for Parcel 8 (39.6447ms)
2017/01/03 22:35:46 Task: Completed for Parcel 9 (40.645ms)
2017/01/03 22:35:46 Task: Completed for Parcel 3 (43.6554ms)
2017/01/03 22:35:46 Task: Completed for Parcel 6 (44.6584ms)
2017/01/03 22:35:46 Task: Completed for Parcel 4 (51.6774ms)
2017/01/03 22:35:46 Task: Completed for Parcel 7 (53.6962ms)
2017/01/03 22:35:46 Task: Completed for Parcel 1 (60.6983ms)
```