# Tool deploy peer cosmos
## Target
Deploy cosmos network using [blogd](https://github.com/phanhoc/blog/tree/dev) repository to k8s.        
Current version support:
- [ ] Create persistence peer
- [ ] Create seed connect to persistence peer
- [x] Create peer connec to seed

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
