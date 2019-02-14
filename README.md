# Keeper
[![Go Report Card](https://goreportcard.com/badge/github.com/atakanyenel/keeper)](https://goreportcard.com/report/github.com/atakanyenel/keeper)

Keeper is a solution for cloud environments where IP address changes with every new deployment of the same VM. Keeper works as a DNS. Every time a a VM is instantiated , Keeper gets it's public IP address and sends it to a server, which has a static domain name. That way, any other device that wants to reach that server can just query it's current IP address from the server. 

Keeper is designed for minimal footprint on start-up , so it compiles to minimal executable possible by removing debug info. 
