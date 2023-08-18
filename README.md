# gorial
Serial port communication and debugger in your terminal

# Usage

You need at least Go version 1.21.0 installed.

Cloning the repo:

``` bash
git clone https://github.com/dio-av/gorial.git
```
Go to the gorial folder and build using the command:

``` bash
go build cmd/gorialApp.go
```

## Running the application

Pass the flags setting the baud rate (115200 default) and the COM Port name:

- Windows
``` bash
gorialApp.exe -baud=115200 -port="COM3"
```
- Linux
``` bash
./gorialApp -baud=115200 -port="/dev/ttyUSB1"
```
