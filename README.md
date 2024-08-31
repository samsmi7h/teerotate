# Tee Rotate
Rotate your log files, within your Go app. No need for external crons. Easy to test & monitor.

## Roadmap
* Add post-rotation hooks: e.g. upload to S3

# How to Use
```go
logger := teerotate.NewRotatingFileLogger("/path/to/dir/", time.Hour)
logger.Print("my first log at: %s", time.Now())
```

## Graceful shutdown
**Important** -- To avoid losing logs when you shutdown, make sure you wait for `logger.Close()` to complete.
Otherwise you risk not flushing the last logs to your file.

Example for handling interrupts:

```go
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

l := teerotate.NewRotatingFileLogger("/tmp", time.Hour)

go func() {
	<-sigChan
	fmt.Println("got signal")
	l.Close()
	fmt.Println("finished")
}()
```

## Log files
While they are being written log files have the form: `2024-08-31T15:50:16.log.live`.

After they have been rotated the log files have the `.live` prefix removed: `2024-08-31T15:50:16.log`.

