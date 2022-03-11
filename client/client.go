package main

import (
	"context"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"strings"
)

var str = `{
  "apiVersion": "admission.k8s.io/v1",
  "kind": "AdmissionReview",
  "request": {
    "uid": "705ab4f5-6393-11e8-b7cc-42010a800002",
    "kind": {"group":"apps","version":"v1","kind":"Deployment"},
    "resource": {"group":"apps","version":"v1","resource":"deployments"},
    "name": "test",
    "namespace": "default",
    "operation": "CREATE",
    "userInfo": {
      "username": "admin",
      "uid": "014fbff9a07c",
      "groups": ["system:authenticated","my-admin-group"],
      "extra": {
        "some-key":["some-value1", "some-value2"]
      }
    },
    "object": {"apiVersion": "apps/v1","kind": "Deployment","metadata": {"labels": {"app": "whoami"},"name": "whoami","namespace": "default"},"spec": {"replicas": 1,"selector": {"matchLabels": {"app": "whoami"}},"template": {"metadata": {"labels": {"app": "whoami"}},"spec": {"containers": [{"image": "hub.fastonetech.com/infra/file-access:v1.1.1.1","imagePullPolicy": "Always","name": "whoami"}],"restartPolicy": "Always"}}}},
    "dryRun": false
  }
}`

func main() {
	cfg := rest.Config{Host: "http://localhost:8080"}
	clientset, err := kubernetes.NewForConfig(&cfg)
	if err != nil {
		klog.Error(err)
		return
	}
	res := clientset.AdmissionregistrationV1().RESTClient().Post().Body(strings.NewReader(str)).Do(context.Background())
	raw, err := res.Raw()
	if err != nil {
		klog.Error(err)
		return
	}
	//klog.V(7).Infoln(raw)
	fmt.Println(string(raw))
}
