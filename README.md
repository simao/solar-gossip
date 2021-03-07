# Solar gossip

This is a Proof Of Concept project to show how a gossip-protocol could help the [solar-protocol](https://github.com/alexnathanson/solar-protocol) be more resilient and scalable.

Solar-gossip uses a [gossip protocol](https://en.wikipedia.org/wiki/Gossip_protocol) to create a network of nodes that determine which node has more energy and should update a DNS service and serve public HTTP requests.

The gossip protocol is based on [SWIM: Scalable Weakly-consistent Infection-style Process Group Membership Protocol](https://research.cs.cornell.edu/projects/Quicksilver/public_pdfs/SWIM.pdf) and is implemented by [memberlist](https://github.com/hashicorp/memberlist/).

This repository is just a proof of concept, that are some edge cases that still need to be handled and it doesn't actually update any DNS.
