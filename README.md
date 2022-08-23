# cert-lister

```
# Use the version matching your os and cpu architecture
> ./cert-lister-****

KIND   NAME                                               KEY                          SERIAL NUMBER                           Age
SECRET kubedb/kubedb-kubedb-provisioner-license           key.txt                      7958275879838045348                     29d
SECRET kubedb/kubedb-kubedb-webhook-server-apiserver-cert tls.crt                      136749661798657898788018791024910026265 9y
SECRET kubedb/kubedb-kubedb-webhook-server-apiserver-cert ca.crt                       196830156477330672160882529336709569966 9y
SECRET kubedb/kubedb-kubedb-webhook-server-license        key.txt                      7958275879838045348                     29d
CFGMAP default/kube-root-ca.crt                           ca.crt                       0                                       9y
CFGMAP kube-node-lease/kube-root-ca.crt                   ca.crt                       0                                       9y
CFGMAP kube-public/kube-root-ca.crt                       ca.crt                       0                                       9y
CFGMAP kube-system/extension-apiserver-authentication     client-ca-file               0                                       9y
CFGMAP kube-system/extension-apiserver-authentication     requestheader-client-ca-file 0                                       9y
CFGMAP kube-system/kube-root-ca.crt                       ca.crt                       0                                       9y
CFGMAP kubedb/kube-root-ca.crt                            ca.crt                       0                                       9y
CFGMAP local-path-storage/kube-root-ca.crt                ca.crt                       0                                       9y
APISVC v1alpha1.mutators.autoscaling.kubedb.com           spec.caBundle                196830156477330672160882529336709569966 9y
APISVC v1alpha1.mutators.dashboard.kubedb.com             spec.caBundle                196830156477330672160882529336709569966 9y
APISVC v1alpha1.mutators.kubedb.com                       spec.caBundle                196830156477330672160882529336709569966 9y
APISVC v1alpha1.mutators.ops.kubedb.com                   spec.caBundle                196830156477330672160882529336709569966 9y
APISVC v1alpha1.mutators.schema.kubedb.com                spec.caBundle                196830156477330672160882529336709569966 9y
APISVC v1alpha1.validators.dashboard.kubedb.com           spec.caBundle                196830156477330672160882529336709569966 9y
APISVC v1alpha1.validators.kubedb.com                     spec.caBundle                196830156477330672160882529336709569966 9y
APISVC v1alpha1.validators.ops.kubedb.com                 spec.caBundle                196830156477330672160882529336709569966 9y
APISVC v1alpha1.validators.schema.kubedb.com              spec.caBundle                196830156477330672160882529336709569966 9y
```
