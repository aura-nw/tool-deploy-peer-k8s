package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"proj1/model"
	"proj1/util"
	"strconv"
	"text/template"
	"time"

	"k8s.io/client-go/kubernetes"
)

func createPersistencePeer(config model.Config, clientset *kubernetes.Clientset) {
	peerName := util.GetInputFromKeyboard("Peer name", "")
	namespace := util.GetInputFromKeyboard("Namespace", config.Namespace)
	chainId := util.GetInputFromKeyboard("ChainId", config.ChainId)

	td := model.ConfigNodePersistencePeerTemplate{Namespace: namespace, PeerName: peerName, ChainId: chainId,
		ExternalAddress: fmt.Sprintf("%s-%s.%s", "peer-persistence", peerName, namespace)}

	t, err := template.ParseFiles(filepath.Join(config.PathTemplate, config.TemplateNodePersistencePeer))
	if err != nil {
		panic(err)
	}
	filePath := filepath.Join(config.PathOutput, strconv.FormatInt(time.Now().Unix(), 10))
	os.Mkdir(filePath, 0777)
	f, err := os.Create(filepath.Join(filePath, config.TemplateNodePersistencePeer))
	if err != nil {
		panic(err)
	}
	err = t.Execute(f, td)
	if err != nil {
		panic(err)
	} else if !config.RunLocal {
		util.ApplyFileYamlAndWatchPod(clientset, filePath, "seed-0", config.Namespace, true, true)
	}
}

func createNodeSeed(config model.Config, clientset *kubernetes.Clientset) {
	td := model.ConfigNodeSeedTemplate{Namespace: config.Namespace}
	t, err := template.ParseFiles(filepath.Join(config.PathTemplate, config.TemplateNodeSeed))
	if err != nil {
		panic(err)
	}
	filePath := filepath.Join(config.PathOutput, strconv.FormatInt(time.Now().Unix(), 10))
	os.Mkdir(filePath, 0777)
	f, err := os.Create(filepath.Join(filePath, config.TemplateNodeSeed))
	if err != nil {
		panic(err)
	}
	err = t.Execute(f, td)
	if err != nil {
		panic(err)
	} else if !config.RunLocal {
		util.ApplyFileYamlAndWatchPod(clientset, filePath, "seed-0", config.Namespace, true, true)
		util.CopyFileFromPod("seed-0", config.Namespace, "/root/.blog/config/genesis.json", "genesis.json")
		util.CreateConfigMapFromFile(clientset, config.Namespace, "seed-genesis", "genesis.json")
	}
}

func createNodePeer(config model.Config, clientset *kubernetes.Clientset) {
	peerName := util.GetInputFromKeyboard("Peer name", "")
	namespace := util.GetInputFromKeyboard("Namespace peer", config.Namespace)
	chainId := util.GetInputFromKeyboard("ChainId", config.ChainId)
	seedNodeAddress := util.GetInputFromKeyboard("Seed node address", config.SeedNodeAddress)
	seedP2PPort := util.GetInputFromKeyboard("Seed node P2P port", config.SeedSvcPortP2P)

	namespaceSeed := util.GetInputFromKeyboard("Namespace seed", config.Namespace)
	//select pod seed to connect
	podSeedName := util.GetInputFromSelect("Select pod seed", util.GetListNamePodInNamespace(clientset, namespaceSeed))

	seedBinaryPath := util.GetInputFromKeyboard("Binary path in seed", config.SeedBinaryPath)

	seedHomePath := util.GetInputFromKeyboard("Home path in seed", config.SeedHomePath)
	//get seed id from seed pod
	seedNodeId, err := util.RunCommandInPod(podSeedName, namespace, fmt.Sprintf("%s tendermint show-node-id --home %s", seedBinaryPath, seedHomePath))

	if err != nil {
		panic(err)
	}
	configInTemplate := model.ConfigNodePeerTemplate{Namespace: namespace, PeerName: peerName, ChainId: chainId, SeedNodeAddress: seedNodeAddress,
		SeedP2PPort: seedP2PPort, SeedNodeID: seedNodeId}

	templateSvc, err := template.ParseFiles(filepath.Join(config.PathTemplate, config.TemplateNodePeerSVC))
	if err != nil {
		panic(err)
	}
	filePath := filepath.Join(config.PathOutput, strconv.FormatInt(time.Now().Unix(), 10))
	os.Mkdir(filePath, 0777)
	file, err := os.Create(filepath.Join(filePath, config.TemplateNodePeerSVC))
	if err != nil {
		panic(err)
	}
	err = templateSvc.Execute(file, configInTemplate)
	if err != nil {
		panic(err)
	} else if !config.RunLocal {
		//create svc
		util.ApplyFileYamlAndWatchPod(clientset, filePath, "peer-0", config.Namespace, false, false)

		//get IP from SVC
		configInTemplate.ClusterIPService = fmt.Sprintf("%s:%s", util.GetClusterIPService(clientset, "peer-"+peerName, namespace), seedP2PPort)

		templateSts, err := template.ParseFiles(filepath.Join(config.PathTemplate, config.TemplateNodePeer))
		if err != nil {
			panic(err)
		}

		//generate yaml with ip svc in config map
		file, err := os.Create(filepath.Join(filePath, config.TemplateNodePeer))
		if err != nil {
			panic(err)
		}

		templateSts.Execute(file, configInTemplate)

		util.ApplyFileYamlAndWatchPod(clientset, filePath, "peer-"+peerName+"-0", config.Namespace, true, false)

		//copy genesis from seed to local
		util.CopyFileFromPod(podSeedName, config.Namespace,
			filepath.Join(seedHomePath, "config/genesis.json"), "genesis.json")

		//copy genesis from local to pvc in peer
		util.CopyFileToPod("peer-"+peerName+"-0", config.Namespace,
			"genesis.json", filepath.Join(config.PeerPvcPath, "genesis-seed.json"))
	}
}

func main() {

	// Load config
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	// Load kubeconfig from current user bastion
	var clientset *kubernetes.Clientset

	if !config.RunLocal {
		clientset, _ = util.LoadConfigKubernetes()
	}

	for {
		result := util.GetInputFromSelect("Select Option", []string{"Create persistence peer node (maintaining)", "Create seed node (maintaining)", "Create peer node", "Exit"})

		switch result {
		case "Create persistence peer node (maintaining)":
			createPersistencePeer(config, clientset)
		case "Create seed node (maintaining)":
			createNodeSeed(config, clientset)
		case "Create peer node":
			createNodePeer(config, clientset)
		case "Exit":
			return
		default:
			break

		}
	}

	// // Load kubeconfig from current user bastion
	// clientset, err := util.LoadConfigKubernetes()
	// if err != nil {
	// 	fmt.Printf("Prompt failed %v\n", err)
	// }
	// pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	// fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
}
