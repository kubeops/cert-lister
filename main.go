/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	cmscheme "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned/scheme"
	"gomodules.xyz/cert"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/duration"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
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

	hc, err := rest.HTTPClientFor(cfg)
	if err != nil {
		return nil, err
	}
	mapper, err := apiutil.NewDynamicRESTMapper(cfg, hc)
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
	fmt.Fprintln(w, "KIND\tNAME\tKEY\tSERIAL NUMBER\tCN\tAge")

	err = ListSecrets(kc, w)
	if err != nil {
		panic(err)
	}
	err = ListConfigMaps(kc, w)
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
		for k, v := range item.Data {
			if crts, err := cert.ParseCertsPEM(v); err == nil {
				for _, crt := range crts {
					fmt.Fprintf(w, "SECRET\t%s/%s\t%s\t%v\t%s\t%s\n", item.GetNamespace(), item.GetName(), k, crt.SerialNumber, crt.Subject.CommonName, ConvertToHumanReadableDateType(crt.NotAfter))
				}
			}
		}
	}
	return nil
}

func ListConfigMaps(kc client.Client, w io.Writer) error {
	var list core.ConfigMapList
	err := kc.List(context.TODO(), &list)
	if err != nil {
		return err
	}
	for _, item := range list.Items {
		for k, v := range item.Data {
			if crts, err := cert.ParseCertsPEM([]byte(v)); err == nil {
				for _, crt := range crts {
					fmt.Fprintf(w, "CFGMAP\t%s/%s\t%s\t%v\t%s\t%s\n", item.GetNamespace(), item.GetName(), k, crt.SerialNumber, crt.Subject.CommonName, ConvertToHumanReadableDateType(crt.NotAfter))
				}
			}
		}
		for k, v := range item.BinaryData {
			if crts, err := cert.ParseCertsPEM(v); err == nil {
				for _, crt := range crts {
					fmt.Fprintf(w, "CFGMAP\t%s/%s\t%s\t%v\t%s\t%s\n", item.GetNamespace(), item.GetName(), k, crt.SerialNumber, crt.Subject.CommonName, ConvertToHumanReadableDateType(crt.NotAfter))
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
					fmt.Fprintf(w, "APISVC\t%s\t%s\t%v\t%s\t%s\n", item.GetName(), "spec.caBundle", crt.SerialNumber, crt.Subject.CommonName, ConvertToHumanReadableDateType(crt.NotAfter))
				}
			}
		}
	}
	return nil
}

// ConvertToHumanReadableDateType returns the elapsed time since timestamp in
// human-readable approximation.
// ref: https://github.com/kubernetes/apimachinery/blob/v0.21.1/pkg/api/meta/table/table.go#L63-L70
// But works for timestamp before or after now.
func ConvertToHumanReadableDateType(timestamp time.Time) string {
	if timestamp.IsZero() {
		return "<unknown>"
	}
	var d time.Duration
	now := time.Now()
	var sign string
	if now.After(timestamp) {
		d = now.Sub(timestamp)
		sign = "- "
	} else {
		d = timestamp.Sub(now)
	}
	return sign + duration.HumanDuration(d)
}
