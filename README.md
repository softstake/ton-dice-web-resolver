# ton-dice-web-resolver
Resolving bets and send them to the dice server to save.

Resolve bets:
 - ton-api grpc 

Save bets:
 - ton-dice-web-server grpc


## build 
```
export GITHUB_TOKEN = 'token'
docker build --build-arg GITHUB_TOKEN="$GITHUB_TOKEN" -t dice-resolver .
```

## run (develop)
```docker run --name dice-resolver --network dice-network -e CONTRACT_ADDR=kf_P3yXtu1ab1CBEl6rG2jb_LtCPIElEoqCZTiffBOte0uLL -e STORAGE_HOST=dice-server -e STORAGE_PORT=5300 -e TON_API_HOST=ton-api -e TON_API_PORT=5400 -d dice-resolver```
 
## ENV VARS
    * CONTRACT_ADDR - Dice contract address in the TON network, required variable, no default value.
    * KEY_FILE_BASE - Dice contract owner private key file name without extension, default value is 'owner'.
    * STORAGE_HOST - Host of the 'ton-dice-web-server' service, requred variable, no default value.
    * STORAGE_PORT - Port of the 'ton-dice-web-server' service, default value is '5300'.
    * TON_API_HOST - Host of the 'ton-api' service, requred variable, no default value.
    * TON_API_PORT - Port of the 'ton-api' service, default value is '5400'.
