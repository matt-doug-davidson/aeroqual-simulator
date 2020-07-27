# Aeroqual Simulator
## Environment Variables
|ENV Variable |  Required  |  Description  |
|:--:|:--:|:--:|
| SN | yes |Serial number of the simulated instrument |
| PORT | yes | TCP port used to connect to the simulated instrument |
| UN | yes | User name for logging into the simulated instrument |
| PW | yes |Password for logging into the simulated instrument |

## CLI
Configuration for this simulator comes from the command line in the form of environment variables. Here are a couple of examples:

Example with executable:
```bash
SN="AQY BB-585" PORT=6969 UN=allied PW="SibiuA\$rqua1" ./aeroqual-simulator
```
Example with go run:
```bash
SN="AQY BB-585" PORT=6969 UN=allied PW="SibiuA\$rqua1" go run aeroqual-simulator.go
```

## Docker

### Build

Example:
```bash
docker build -t mddofapex/aeroqual-simulator:v0.0.2 .
```

### Run

Example:
```bash
docker run --env SN="AQY BB-585" --env PORT=6969 --env UN=allied --env PW="SibiuA\$rqua1" -p 6969:6969 mddofapex/aeroqual-simulator:v0.0.1
```

### Push

```bash
docker push mddofapex/aeroqual-simulator:v0.0.2
```