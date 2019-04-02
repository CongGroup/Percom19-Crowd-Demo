# Introduction

This demo is for paper **Aggregating Crowd Wisdom via Blockchain: A Private, Correct, and Robust Realization**. 

A data consumer can use this app to get an sum aggregation of data providers. The smart contract is deployed at Kovan testnet.
You can see from [here](https://kovan.etherscan.io/address/0x2a152cad883162011dcdf696a4861edc81619e14)

## Every part of this project

### back-end

Provide core function for dealing with front-end and blockchain. 

### front-end

Subpath /admin: An admin page providing friendly interface for data consumer and service provider  to manage crowdsourcing task.
Rootpath /: Provide friendly interface for data providers to submit their data and get their reward.

### client

A simulation of data consumer. You can set the number of data consumers to exploring this demo.

## Usage (Only for running in cityu server)

1. Clone this project.

``` bash
git clone https://github.com/CongGroup/Percome19-Crowd-Demo.git
```

2. Pull image

```bash
docker-compose pull
```

3. Run
```bash
docker-compose up
```
