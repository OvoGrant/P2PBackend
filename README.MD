# Project Title

COMP4911 P2P Backend

## Requirements

- Any machine that can be accessed publicly

## Installation

1. Clone the git repository
2. build the project with go build -o p2p_backend
3. run the compiled executable

   ```sh
   go build . -o p2p_backend 
   ```

## Usage

This program acts as an index server for clients running the corresponding client program. Currently 
clients are hard-corded to point to p2pcomp4911.xyz, but this can be changed by modifying a single line of code in the client code