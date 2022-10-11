package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"regexp"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	nsClient := clientset.CoreV1().Namespaces()
	nsList, err := nsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("I found %d namespaces \n", len(nsList.Items))

	const excludedNs string = "default|^kube-."

	rx, err := regexp.Compile(excludedNs)
	if err != nil {
		panic(err)
	}

	for _, b := range nsList.Items {
		currNs := b.ObjectMeta.Name

		fmt.Printf("Current namespace  %s\n ", currNs)

		found := rx.MatchString(currNs)
		if found {
			fmt.Println("Skipping...")
			continue
		}
		podsClient := clientset.CoreV1().Pods(currNs)

		podsList, err := podsClient.List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		podsNumber := len(podsList.Items)
                
                if podsNumber != 0 {
                    fmt.Println("Skipping...")
                    continue
                }
	
		svcClient := clientset.CoreV1().Services(currNs)

		svcList, err := svcClient.List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		
                svcNumber := len(svcList.Items)
		if svcNumber != 0 {
                    fmt.Println("Skipping...")
                    continue
                }

                fmt.Println(currNs, " namespace will be deleted...")

                deletePolicy := metav1.DeletePropagationForeground
                if err := nsClient.Delete(context.TODO(), currNs, metav1.DeleteOptions{
                        PropagationPolicy: &deletePolicy,
                }); err != nil {
                        panic(err)
                }
                fmt.Println(currNs, " namespace has beed deleted successfully")
	}

}
