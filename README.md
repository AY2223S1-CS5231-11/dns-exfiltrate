# DNS Data Exfiltration

## Setting Up DNS Records

In order for DNS exfiltration to work, we need to set up our own DNS server that is the authoritative nameserver for some domain.

1. Create an A record pointing to the IP address of the DNS server to exfiltrate data to.
1. Create an NS record pointing to the previously created A record.

   This NS record will be used as the domain for all of our DNS queries.

## Design Decisions

### Identifying Targets

Targets cannot be identified via the IP address seen by the server.
This is because our DNS server will not be the local DNS server that the target talk to.
As such, regardless of the DNS query type (iterative or recursive), our DNS server never talks directly to the target.

In addition, some local DNS servers have a range of IP addresses which they use to communicate to other DNS servers.
Given the connectionless nature of UDP (the transport protocol used by DNS), every DNS request that we receive from the same target can come from different IP addresses.
