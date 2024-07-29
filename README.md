# One Time Password (OTP) App 

## Why another OTP app?
Authy (Twillio) no longer supports their desktop app and recently upon launch of the app the message "Device Removed" was received and there's no way to reconnect due to the message "The device does not meet the minimum integrity requirements" thus began the [GOTP](https://github.com/yelsaw/gotp) project for simple desktop use.

## What to expect?
This is a no frills OTP application which only supports TOTP.
In the future this project may become more end-user friendly, but for now it's bare-bones and has no support or docs other than whatever is provided herein, which right now isn't much :)

## Limitations
As mentioned, only supports [Time-based one-time password](https://en.wikipedia.org/wiki/Time-based_one-time_password) not [HMAC-based one-time password](https://en.wikipedia.org/wiki/HMAC-based_one-time_password)

## Getting started
Clone the repo, build or run the app.

Store and retrieve
```
go run <otp-string-key> <full-totp-url> 
```
Retrieve only
```
go run <otp-key>
```
