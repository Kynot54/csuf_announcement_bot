# csuf_announcement_bot
Discord Bot Written in Golang to Provide Regularly Scheduled Updates on Campus News from CSUF into our ACM@CSUF Chapter's Discord Channel

## Prerequisites

---

1. Installation of [Go](https://go.dev/dl/) and [SQLite3](https://www.sqlite.org/download.html) on Host Machine
2. Create a [Discord Bot](https://discord.com/developers/applications) on the Discord Developer Portal

## Compilation and Run Steps

---

1. Git Clone this Repository into a Folder of your choice
2. Run the following script that is stored in __./bin/init.sh__
    sh ./bin/init.sh

3. To Compile and Run the Project
    go build
    ./csuf_announcement_bot

### Notes

---
This Bot can be re-used for other organizations that utilize RSS Feeds to Retrieve News, Updates, and More...
