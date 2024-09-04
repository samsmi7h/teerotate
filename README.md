# Tee Rotate
Rotate your log files, within your Go app. No need for external crons. Easy to test & monitor.

## Roadmap
* Rotate condition includes min size
* Handle situation where underlying file has been (re)moved

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


## Post-Rotation Hooks
Once you've rotated out a log file, you probably want to do something with it.

Rather than needing to setup a cron job to do this work, you can provide a hook to do it.

All hooks run in the background.

**Important**: Hooks should be added BEFORE any logs are sent. Otherwise, behaviour is undefined.

```go
logger.WithPostRotationHook(func() {....})
```
