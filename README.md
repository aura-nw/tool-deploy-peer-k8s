# Tool deploy peer cosmos
## Target
Deploy cosmos network using [blogd](https://github.com/phanhoc/blog/tree/dev) repository to k8s.        
Current version support:
- [ ] Create persistence peer
- [ ] Create seed connect to persistence peer
- [x] Create peer connect to seed

## Prerequisite

Go version 1.16 or above [Install go](https://www.google.com/url?sa=t&rct=j&q=&esrc=s&source=web&cd=&cad=rja&uact=8&ved=2ahUKEwiQorLn9_vzAhUPOSsKHccqDlgQFnoECAIQAQ&url=https%3A%2F%2Fgolang.org%2Fdoc%2Finstall&usg=AOvVaw2iPQ4PF_CePbLDqF11JL33)  
Kubectl in current home directory (has .kubeconfig directory)
## Installation
```
cd $repo
go mod download
go build
```

## Usage

Select option: 
```
cd $repo
go run main.go 
    ```
    Select Option: 
    Create persistence peer node
    Create seed node
  ▸ Create peer node
    Exit
    ```    
```
Enter information necessary:
```
    # default value in dev-qa
    Create peer node
    Peer name: bn
    Namespace peer: akc-testnet
    ChainId: akc-test
    Seed node address: node-cosmos-seed
    Namespace seed: akc-testnet
    Pod name seed: node-cosmos-seed-0
    Home path in seed: /data-node/seed
```
Peer name: peer name (moniker will be `peer-{peer name}`)  
Namespace peer: Namespace in k8 where pod peer will be created  
ChainId: Chain Id in seed  
Seed node address: seed peer address domain name (svc name in k8)  
Namespace seed: Namespace where pod seed existed  
Pod name seed: pod name seed in k8  
Home path in seed: Home path network where has directory /config/genesis.json (in order to copy genesis from seed to another peer)

## Script

### Create account
Run command in pod peer:
```
    kubectl exec -it peer-$peername-0 -- bash
    blogd keys add $account
```
### Call faucet to get token
Run command in any pod:
```
    kubectl exec -it peer-$peername-0 -- bash
    curl --location --request GET 'http://node-cosmos-faucet/?address=$address'
    # for i in {1..5}; do curl --location --request GET 'http://node-cosmos-faucet/?address=$address';sleep 2; done
    blogd query bank balances $address
```
address is wallet address created in previous step
### Upgrade peer to validator

```
    # Send token to one validator to active new validator
    blogd tx bank send $address $addressValidator 1stake --chain-id akc-test

    blogd tx staking create-validator \
    --amount=10000000stake \
    --pubkey=$(blogd tendermint show-validator) \
    --moniker="peer-hb" \
    --commission-rate="1" \
    --commission-max-rate="1" \
    --commission-max-change-rate="1" \
    --min-self-delegation="1" \
    --from=peer-hb \
    --chain-id=akc-test \
    --gas=auto \
    --gas-adjustment=1.4 \
    --keyring-backend=test

    blogd query staking validators
    blogd query tendermint-validator-set
```
