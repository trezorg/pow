# "Word of Wisdom" TCP-server with protection from DDOS based on Proof of Work

[![Actions Status]( https://github.com/trezorg/pow/actions/workflows/go.yml/badge.svg)](https://github.com/trezorg/pow/actions)

## 1. Description
This project is a solution for some interview question on Golang.

## 2. Getting started

### 2.1 Start server and client by docker-compose:
```
make start
```

### 2.2 Launch tests:
```
make test
```

### 2.3 Check logs:
```
make logs
```

### 2.4 Stop server and client by docker-compose:
```
make stop
```

### 2.5 Build:
```
make build
```

## 3. Proof of Work
Idea of Proof of Work for DDOS protection is that client, which wants to get some resource from server, 
should firstly solve some challenge from server. 
This challenge should require more computational work on client side and verification of challenge's solution - much less on the server side.

### 4 Selection of an algorithm
There is some different algorithms of Proof Work. 
I compared next three algorithms as more understandable and having most extensive documentation:
+ [Merkle tree](https://en.wikipedia.org/wiki/Merkle_tree)
+ [Hashcash](https://en.wikipedia.org/wiki/Hashcash)
+ [Guided tour puzzle](https://en.wikipedia.org/wiki/Guided_tour_puzzle_protocol)

After comparison, I chose Hashcash. Other algorithms have next disadvantages:
+ In Merkle tree server should do too much work to validate client's solution. For tree consists of 4 leaves and 3 depth server will spend 3 hash calculations.
+ In guided tour puzzle client should regularly request server about next parts of guide, that complicates logic of protocol.

Hashcash, instead has next advantages:
+ simplicity of implementation
+ lots of documentation and articles with description
+ simplicity of validation on server side
+ possibility to dynamically manage complexity for client by changing required leading zeros count

Of course Hashcash also has disadvantages like:

1. Compute time depends on power of client's machine. 
For example, very weak clients possibly could not solve challenge, or too powerful computers could implement DDOS-attackls.
But complexity of challenge could be dynamically solved by changing of required zeros could from server.
2. Pre-computing challenges in advance before DDOS-attack. 
Some clients could parse protocol and compute many challenges to apply all of it in one moment.
It could be solved by sending a seed value to client and store this seed on the server. 
