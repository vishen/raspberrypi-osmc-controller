# OSMC Remote Controller (OSMC / Remote PIS-0445) on Raspberry PI on Linux

A Go example library for using OSMC remote control on Raspberry PI on Linux.


## Building for Raspberry PI

	GOOS=linux GOARCH=arm GOARM=7 go build -o ./build/armv7-linux/osmc-controller

## Remote Controller Functions

```
- HOME
- ARROW_UP
- ARROW_LEFT
- ARROW_RIGHT
- ARROW_DOWN
- OK
- INFO
- RETURN
- PLAY_PAUSE
- STOP
- PREV
- NEXT
- LIST
```
