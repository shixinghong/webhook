package utils

import (
	"fmt"
	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
	"regexp"
)

const exp = `^(hub.myit.fun)(.*file-access)`

type patchSpec struct {
	Option string                `json:"op"`
	Path   string                `json:"path"`
	Value  appsv1.DeploymentSpec `json:"value"`
}

func AdmitDeployments(ar admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	klog.Info("admitting deployments")
	deploymentResource := metav1.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}

	if ar.Request.Resource != deploymentResource {
		err := fmt.Errorf("expect resource to be %s", deploymentResource)
		klog.Error(err)
		return ToAdmissionResponse(err)
	}
	raw := ar.Request.Object.Raw
	deployment := &appsv1.Deployment{}
	deserializer := Codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, deployment); err != nil {
		klog.Error(err)
		return ToAdmissionResponse(err)
	}

	klog.Info(deployment)

	// 定义准入规则
	reviewResponse := admissionv1.AdmissionResponse{}
	containers := deployment.Spec.Template.Spec.Containers
	for i, container := range containers {
		img := regexp.MustCompile(exp).ReplaceAllString(container.Image, "r.myit.fun:5000$2")
		containers[i].Image = img
	}

	// 测试注入docker in docker 容器
	//newDeployment := InjectDockerInDockerContainer(deployment)
	klog.Info(deployment.Spec.Template.Spec.Containers)

	reviewResponse.Allowed = true
	dpSpec := []patchSpec{
		{
			Option: "replace",
			Path:   "/spec",
			Value:  deployment.Spec,
		},
	}

	dpSpecJson, err := json.Marshal(dpSpec)
	if err != nil {
		klog.Error(err)
		return ToAdmissionResponse(err)
	}

	klog.Info(string(dpSpecJson))

	reviewResponse.Patch = dpSpecJson
	jsonPatchType := admissionv1.PatchTypeJSONPatch
	reviewResponse.PatchType = &jsonPatchType

	return &reviewResponse
}
