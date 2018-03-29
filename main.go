package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/golang/glog"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/feynmanliang/epoxy/controller"
	"github.com/feynmanliang/epoxy/server"
)

func createConnection() (*kubernetes.Clientset, error) {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		glog.Fatal(err.Error())
	}

	// create the clientset
	return kubernetes.NewForConfig(config)
}

func main() {
	// establish kubernetes connection
	clientset, err := createConnection()
	if err != nil {
		glog.Fatal(err.Error())
	}

	// start the controller
	stopController := make(chan struct{})
	defer close(stopController)
	go controller.NewController(clientset).Run(1, stopController)

	// Start the server
	stopServer := make(chan struct{})
	defer close(stopServer)
	go server.Run(stopServer)

	// Wait forever
	select {}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
