# Smart Home V2

This repository is an updated and improved version of [Smart Home V1](https://github.com/S4ndf1re/SmartHome).
While improvements can be found in the frontend, scipting / plugin system and overall architecture,
the overalls platform security was degraded, especially regarding the previously encrypted mqtt message communication.
This was a time constraint, as this remained a hobby project.

Non the less, the following features are implemented:

- Scripting using JavaScript, each script exposing
  - Direct Communication With MQTT Brokers
  - Web Component Creation with live updates using WebSockets
- Esp8266 scripts for reading Mifare1k Chips
- Esp8266 scripts for actuating a doorlock mechanism
- Simple User-Management

Similar to V1, this version was in use for roughly two years and provided improved stability over V1.
