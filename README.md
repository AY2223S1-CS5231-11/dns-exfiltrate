# DNS Data Exfiltration

## Setting Up DNS Records

In order for DNS exfiltration to work, we need to set up our own DNS server that is the authoritative nameserver for some domain.

1. Create an A record pointing to the IP address of the DNS server to exfiltrate data to.
1. Create an NS record pointing to the previously created A record.

   This NS record will be used as the domain for all of our DNS queries.
