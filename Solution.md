# be-code-challenge

## Overview

The entrypoint to the code is in cmd/challenge/challenge.go. The service has the following HTTP handlers:

- _/hourly_ - use GET method. Returns the amount of fees being paid for Ethereum transactions per hour in the following JSON format:
```
[
  {
    "t": 1603114500,
    "v": 123.45
  },
  ...
]
```

The /hourly endpoint doesn't have query parameters for filtering by time range, that's a trade-off I made. I also didn't optimize performance of
a service or a database query. According to the requirement to make the service production ready I added _/readiness_ 
and _/liveness_ endpoints to make the service run in a Kubernetes cluster. 

- _/readiness_ - check if the database is ready and returns a 500 status if it's not.
- _/liveness_ - return simple status info if the service is alive.

I moved the _docker-compose.yml_ file to ./infra directory and created a Makefile with _make image, up, test and clean_ commands.
So you will have to run _make up_ instead of _docker-compose up_ to run the service. The installation and running steps described below.

## Installation

1. Clone this repository in the current directory:

   ```
   git clone https://github.com/illyasch/be-code-challenge
   ```

2. Build a Docker image:

   ```bash
   make image
   ```

3. Start the local development environment (uses Docker):

   ```
   make up
   ```

   At this point you should have the challenge service running. To confirm the state of the running Docker container, run

   ```
   $ docker ps
   ```

## How to

### Run unit tests

from the docker container

```
make test
```

### Run manual tests

   ```
   $ curl http://localhost:8080/hourly
[{"t":1599436800,"v":17781937815.707344},{"t":1599440400,"v":25796173158.88589},{"t":1599444000,"v":34821055861.44104},{"t":1599447600,"v":29814493424.40487},{"t":1599451200,"v":27821774201.37403},{"t":1599454800,"v":25575311595.65763},{"t":1599458400,"v":33138772595.681362},{"t":1599462000,"v":35671504748.405235},{"t":1599465600,"v":29861742077.137997},{"t":1599469200,"v":31870528305.547283},{"t":1599472800,"v":29492777739.82145},{"t":1599476400,"v":28476895183.216103},{"t":1599480000,"v":31458479835.443005},{"t":1599483600,"v":36114881483.45387},{"t":1599487200,"v":39990571952.80128},{"t":1599490800,"v":32351461366.072742},{"t":1599494400,"v":35702769826.03566},{"t":1599498000,"v":28225350833.576153},{"t":1599501600,"v":23619974534.4661},{"t":1599505200,"v":20792555662.410378},{"t":1599508800,"v":20324756156.83894},{"t":1599512400,"v":18641503209.089615},{"t":1599516000,"v":16778240397.619787},{"t":1599519600,"v":17399949324.595974}]
   ```