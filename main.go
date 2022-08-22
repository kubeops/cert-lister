package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	cmscheme "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned/scheme"
	"gomodules.xyz/cert"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2/klogr"
	aggapi "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	aggscheme "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

func NewUncachedClient() (client.Client, error) {
	ctrl.SetLogger(klogr.New())
	cfg := ctrl.GetConfigOrDie()
	cfg.QPS = 100
	cfg.Burst = 100

	mapper, err := apiutil.NewDynamicRESTMapper(cfg)
	if err != nil {
		return nil, err
	}

	scheme := runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		return nil, err
	}
	// apiservices
	if err := aggscheme.AddToScheme(scheme); err != nil {
		return nil, err
	}
	// cert-manager
	if err := cmscheme.AddToScheme(scheme); err != nil {
		return nil, err
	}

	return client.New(cfg, client.Options{
		Scheme: scheme,
		Mapper: mapper,
		//Opts: client.WarningHandlerOptions{
		//	SuppressWarnings:   false,
		//	AllowDuplicateLogs: false,
		//},
	})
}

func main() {
	kc, err := NewUncachedClient()
	if err != nil {
		panic(err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)

	err = ListSecrets(kc, w)
	if err != nil {
		panic(err)
	}
	err = ListAPIServices(kc, w)
	if err != nil {
		panic(err)
	}

	w.Flush()
}

func ListSecrets(kc client.Client, w io.Writer) error {
	var list core.SecretList
	err := kc.List(context.TODO(), &list)
	if err != nil {
		return err
	}
	for _, item := range list.Items {
		if item.Type == core.SecretTypeTLS {
			if v, ok := item.Data[core.TLSCertKey]; ok {
				if crts, err := cert.ParseCertsPEM(v); err == nil {
					for _, crt := range crts {
						fmt.Fprintf(w, "SECRET\t%s/%s\t%v\n", item.GetNamespace(), item.GetName(), crt.SerialNumber)
					}
				}
			}
		}
	}
	return nil
}

func ListAPIServices(kc client.Client, w io.Writer) error {
	var list aggapi.APIServiceList
	err := kc.List(context.TODO(), &list)
	if err != nil {
		return err
	}
	for _, item := range list.Items {
		if len(item.Spec.CABundle) > 0 {
			if crts, err := cert.ParseCertsPEM(item.Spec.CABundle); err == nil {
				for _, crt := range crts {
					fmt.Fprintf(w, "APISVC\t%s\t%v\n", item.GetName(), crt.SerialNumber)
				}
			}
		}
	}
	return nil
}
