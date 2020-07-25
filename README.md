# Aeroqual Simulator

|ENV Variable |  Required  |  Description  |
|:--:|:--:|:--:|
| SN | yes |Serial number of the simulated instrument |
| PORT | yes | TCP port used to connect to the simulated instrument |
| UN | yes | User name for logging into the simulated instrument |
| PW | yes |Password for logging into the simulated instrument |

Example:
```bash
SN="AQY BB-585" PORT=6969 UN=allied PW="SibiuA\$rqua1" ./aeroqual-simulator

SN="AQY BB-585" PORT=6969 UN=allied PW="SibiuA\$rqua1" go run aeroqual-simulator.go
```