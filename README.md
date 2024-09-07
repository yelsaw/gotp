# One-Time Password (OTP) App

## Why another OTP app?
Authy (Twillio) no longer supports their desktop app and recently upon launch of the app the message "Device Removed" was received and there's no way to reconnect due to the message "The device does not meet the minimum integrity requirements" thus began the [GOTP](https://github.com/yelsaw/gotp) project for simple desktop use.

## What to expect?
This is a no frills OTP application which only supports TOTP.
In the future this project may become more end-user friendly, but for now it's bare-bones and has no support or docs other than whatever is provided herein, which right now isn't much :)

## Limitations
As mentioned, only supports [Time-based one-time password](https://en.wikipedia.org/wiki/Time-based_one-time_password) not [HMAC-based one-time password](https://en.wikipedia.org/wiki/HMAC-based_one-time_password)

## Getting started

Navigate to the [Releases Page](https://github.com/yelsaw/gotp/releases) and download the code, or use a pre-compiled binary for your OS.
 - Linux (Tested on Debian)
 - Darwin (Tested on Intel and M3)
 - Windows (Untested as of first release)

*OR*

Clone the repo, build or run the app.

OTP code from URL
otpath provided is an example of a parsed QRCode from a service.
```
go run main.go "otpauth://totp/AppName:you@youremail.com?algorithm=SHA1&digits=6&issuer=AppName&period=30&secret=SECRET_STRING"
```

Returns:
```
Token: 123456
Expires in 30 seconds

Press q to quit

```
