# rainbird-api
This is a restful API written in golang that talks to both the rainbird.com phone-api as well as the local API provided on the RainBird controller itself. It does this by utilizing jsonrpc 2.0 calls to match what the phone applications utilize. 

*This project has no affiliation with RainBird and is purely a reverse engineered effort. As a result it may break at any time as RainBird could potentially update their firmware and cloud based API to no longer work with this software.*

# Docker

Run a typical docker build from the project directory

`docker build -t rainbird-api .`

Then you can run it it with the port you want exposed. For example if you want port 8000 instead of the default 8080.

`docker run -p 8000:8080 rainbird-api`

## Warning

*You use this software at your own risk*