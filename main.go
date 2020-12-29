package main

import (
	"context"
	"encoding/json"
	"github.com/k8s-autoops/autoops"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"log"
	"net/http"
	"os"
	"regexp"
)

const (
	denyMessageNoAutoCreatedIngressService = "禁止使用 [+工作负载] 按钮， 这会让 Rancher 使用非标准 Kubernetes 行为。请先在 [工作负载] 编辑页面声明 集群IP 类型的端口规则，然后在本页面使用 [+服务] 按钮，并使用下拉框选择端口"
)

var (
	regexpAutoCreatedIngressService = regexp.MustCompile(`^ingress-[a-fA-F0-9]{32}$`)
)

func exit(err *error) {
	if *err != nil {
		log.Println("exited with error:", (*err).Error())
		os.Exit(1)
	} else {
		log.Println("exited")
	}
}

func main() {
	var err error
	defer exit(&err)

	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	s := &http.Server{
		Addr: ":443",
		Handler: autoops.NewMutatingAdmissionHTTPHandler(
			func(ctx context.Context, request *admissionv1.AdmissionRequest, patches *[]map[string]interface{}) (deny string, err error) {
				var buf []byte
				if buf, err = request.Object.MarshalJSON(); err != nil {
					return
				}
				if request.Resource.Resource == "services" {
					// disable ingress-xxxxxxxxxxxx auto created service
					var obj corev1.Service
					if err = json.Unmarshal(buf, &obj); err != nil {
						return
					}
					if regexpAutoCreatedIngressService.MatchString(obj.Name) {
						deny = denyMessageNoAutoCreatedIngressService
						log.Printf("Create denied for Service: %s", obj.Name)
						return
					}
				}
				if request.Resource.Resource == "ingresses" {
					// disable usage or creation for ingress-xxxxx
					var obj extensionsv1beta1.Ingress
					if err = json.Unmarshal(buf, &obj); err != nil {
						return
					}
					for _, rule := range obj.Spec.Rules {
						if rule.HTTP == nil {
							continue
						}
						for _, path := range rule.HTTP.Paths {
							if regexpAutoCreatedIngressService.MatchString(path.Backend.ServiceName) {
								deny = denyMessageNoAutoCreatedIngressService
								log.Printf("Create/Update denied for Ingress: %s", obj.Name)
								return
							}
						}
					}
				}
				return
			},
		),
	}

	if err = autoops.RunAdmissionServer(s); err != nil {
		return
	}
}
