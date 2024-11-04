# Prate - Simple CLI chat tool

**Prate** is a simple CLI based chat tool to use inside your office/private network 😄

# Internal Structure

Scanner :- This module scans entire subnet and keep list of _node_

Node :- Computet/ User who listen to the message. technically computer listening in configured port default port is **4387**

1. coroutine to check the users (scanner)
2. coroutine to send message to each users if needed
