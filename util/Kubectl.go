package util

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ApplyFileYamlAndWatchPod(clientset *kubernetes.Clientset, filePath string, podName string, namespace string, isRunning bool, isReady bool) error {
	fmt.Println("File path create node seed: ", filePath)
	output, err := exec.Command("kubectl", "apply", "-f", filePath).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", output)

	// Wait for pod to be ready
	if isRunning || isReady {
		waitingPodReady(clientset, podName, namespace, isRunning, isReady)
	}
	return nil
}

func waitingPodReady(clientset *kubernetes.Clientset, podName string, namespace string, isRunning bool, isReady bool) {
	spinner := CreateSpinner()
	spinner.Start()
	spinner.Message("waiting for pod to be ready")
	for {
		pod, _ := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
		if (isRunning && pod.Status.Phase == "Running") && ((isReady && pod.Status.ContainerStatuses[0].Ready == true) || !isReady) {
			spinner.Message("pod is ready")
			break
		}
		time.Sleep(2 * time.Second)
	}
	spinner.Stop()
}
func CopyFileFromPod(podName string, namespace string, src string, dest string) error {
	fmt.Println("Copy file from pod: ", podName)
	_, err := exec.Command("kubectl", "cp", fmt.Sprintf("%s/%s:/%s", namespace, podName, src), dest).Output()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func CopyFileToPod(podName string, namespace string, src string, dest string) error {
	fmt.Println("Copy file to pod: ", podName)
	_, err := exec.Command("kubectl", "cp", src, fmt.Sprintf("%s/%s:/%s", namespace, podName, dest)).Output()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func GetListNamePodInNamespace(clientset *kubernetes.Clientset, namespace string) []string {
	var listNamePod []string
	pods, _ := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	for _, pod := range pods.Items {
		listNamePod = append(listNamePod, pod.Name)
	}
	return listNamePod
}
func CreateConfigMapFromFile(clientset *kubernetes.Clientset, namespace string, nameConfigMap string, filePath string) error {
	fmt.Println("Create config map: ", nameConfigMap)
	_, err := exec.Command("kubectl", "create", "configmap", nameConfigMap, "--from-file", filePath, "-n", namespace).Output()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func RunCommandInPod(podName string, namespace string, command string) (string, error) {
	fmt.Println("Run command in pod: ", podName)
	output, err := exec.Command("kubectl", "exec", "-it", podName, "-n", namespace, "--", "bash", "-c", command).Output()
	if err != nil {
		log.Fatal(err)
	}
	outputString := string(output)
	// remove '\n' at the end of string
	outputString = outputString[:len(outputString)-1]
	return outputString, err
}

func GetClusterIPService(clientset *kubernetes.Clientset, nameService string, namespace string) string {
	fmt.Println("Get IP service: ", nameService)
	service, _ := clientset.CoreV1().Services(namespace).Get(context.TODO(), nameService, metav1.GetOptions{})
	return service.Spec.ClusterIP
}
