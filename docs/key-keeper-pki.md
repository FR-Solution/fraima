# SSL

| Key usage extension | Description  |
|:--------------------|:-------------| 
| Digital signature   |	Use when the public key is used with a digital signature mechanism to support security services other than non-repudiation, certificate  or CRL signing. A digital signature is often used for entity authentication and data origin authentication with integrity.|
| Non-repudiation     | Use when the public key is used to verify digital signatures used to provide a non-repudiation service. Non-repudiation protects against the signing entity falsely denying some action (excluding certificate or CRL signing). |
| Key encipherment    | Use when a certificate will be used with a protocol that encrypts keys. An example is S/MIME enveloping, where a fast (symmetric) key is encrypted with the public key from the certificate. SSL protocol also performs key encipherment. |
| Data encipherment   | Use when the public key is used for encrypting user data, other than cryptographic keys. |
| Key agreement       | Use when the sender and receiver of the public key need to derive the key without using encryption. This key can then can be used to encrypt messages between the sender and receiver. Key agreement is typically used with Diffie-Hellman ciphers. |
| Certificate signing | Use when the subject public key is used to verify a signature on certificates. This extension can be used only in CA certificates. |
| CRL signing         | Use when the subject public key is to verify a signature on revocation information, such as a CRL. |
| Encipher only       | Use only when key agreement is also enabled. This enables the public key to be used only for enciphering data while performing key agreement. |
| Decipher only       | Use only when key agreement is also enabled. This enables the public key to be used only for deciphering data while performing key agreement. |

### CA
| CA  | Type  |Root CA| Description | 
|:---|:-----|:-----|:-----------|
| kubernetes | intermediate | root | Основной CA кластера |
| etcd | intermediate | root | Основной CA ETCD кластера|
| fron-proxy | intermediate | root | CA для [Aggregation Layer](https://kubernetes.io/docs/tasks/extend-kubernetes/configure-aggregation-layer/)|
| root | root-ca | - | Root CA, Корневой сертификат инфраструктуры |


### ETCD

| Component  | ARG  | Description | Usages | CN |  CA |  OU | 
|:-----------|:-----| :------------|:----|:----|--|:--|
| etcd | --trusted-ca-file      | Trusted certificate authority |  ca| - | etcd | - | - |
|      | --cert-file            | Certificate used for SSL/TLS connections to etcd. When this option is set, advertise-client-urls can use the HTTPS schema.т | DigitalSignature,KeyEncipherment,ServerAuth | system:etcd-server | etcd | system:etcd |
|      | --key-file             | Key for the certificate. Must be unencrypted.| - | - | etcd | - | 
|      | --peer-trusted-ca-file | Trusted certificate authority. | ca |  - | etcd | - | 
|      | --peer-cert-file       | Certificate used for SSL/TLS connections between peers. This will be used both for listening on the peer address as well as sending requests to other peers.| DigitalSignature,KeyEncipherment,ServerAuth,ClientAuth | system:etcd-peer | etcd | system:etcd | 
|      | --peer-key-file        | Key for the certificate. Must be unencrypted.| - | - |  etcd |- |
|      | --client-cert-auth=true        | If an etcd server is launched with the option --client-cert-auth=true, the field of Common Name (CN) in the client’s TLS cert will be used as an etcd user. In this case, the common name authenticates the user and the client does not need a password. | bool | - | - | - |

#### Базовая настройка политик:
```bash
etcdctl user add system:etcd-peer --interactive=false # If an etcd server is launched with the option --client-cert-auth=true, the field of Common Name (CN) in the client’s TLS cert will be used as an etcd user. In this case, the common name authenticates the user and the client does not need a password.
etcdctl role add system:etcd
etcdctl role grant-permission system:etcd read /
etcdctl role grant-permission system:etcd write /
etcdctl user grant-role system:etcd-peer system:etcd
```

### KBE-APISERVER

| Component  | ARG  | Description | Type | CN |  CA |  OU | 
|:-----------|:-----| :------------|:----|:----|:--|:--|
| kube-apiserver | --client-ca-file     | If set, any request presenting a client certificate signed by one of the authorities in the client-ca-file is authenticated with an identity corresponding to the CommonName of the client certificate. | CA | - | kubernetes | - |
| | --tls-cert-file                     | File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated after server cert). If HTTPS serving is enabled, and --tls-cert-file and --tls-private-key-file are not provided, a self-signed certificate and key are generated for the public address and saved to the directory specified by --cert-dir. | DigitalSignature,KeyEncipherment,ServerAuth | system:kube-apiserver-server         | kubernetes | - |
| | --tls-private-key-file              | File containing the default x509 private key matching --tls-cert-file. | - | - | - |
| | --etcd-cafile                      | SSL Certificate Authority file used to secure etcd communication. | CA | - |etcd | - |
| | --etcd-certfile                     | SSL certification file used to secure etcd communication. | DigitalSignature,KeyEncipherment,ClientAuth | system:kube-apiserver-etcd-client    | etcd | - |
| | --etcd-keyfile                      | SSL key file used to secure etcd communication. | - | - | - | 
| | --kubelet-client-certificate        | Path to a client cert file for TLS. | DigitalSignature,KeyEncipherment,ClientAuth | system:kube-apiserver-kubelet-client | kubernetes | system:masters |
| | --kubelet-client-key                | Path to a client key file for TLS. | - | - | - | 
| | --proxy-client-cert-file            | Client certificate used to prove the identity of the aggregator or kube-apiserver when it must call out during a request. This includes proxying requests to a user api-server and calling out to webhook admission plugins. It is expected that this cert includes a signature from the CA in the --requestheader-client-ca-file flag. That CA is published in the 'extension-apiserver-authentication' configmap in the kube-system namespace. Components receiving calls from kube-aggregator should use that CA to perform their half of the mutual TLS verification. | DigitalSignature,KeyEncipherment,ClientAuth | system:kube-apiserver-front-proxy-client | front-proxy | system:masters |
| | --proxy-client-key-file             | Private key for the client certificate used to prove the identity of the aggregator or kube-apiserver when it must call out during a request. This includes proxying requests to a user api-server and calling out to webhook admission plugins. | - | - | - |
| | --requestheader-client-ca-file      | Root certificate bundle to use to verify client certificates on incoming requests before trusting usernames in headers specified by --requestheader-username-headers. WARNING: generally do not depend on authorization being already done for incoming requests. | CA | - | front-proxy | - |
| | --service-account-key-file          | File containing PEM-encoded x509 RSA or ECDSA private or public keys, used to verify ServiceAccount tokens. The specified file can contain multiple keys, and the flag can be specified multiple times with different files. If unspecified, --tls-private-key-file is used. Must be specified when --service-account-signing-key-file is provided | DigitalSignature,KeyEncipherment | system:kubernetes-sa | kubernetes | - |
| | --service-account-signing-key-file  | Path to the file that contains the current private key of the service account token issuer. The issuer will sign issued ID tokens with this private key. | key | - | - | 
| kubeconfig | ca |- | - | - | - | - |
|            | client |- | - | - | - | - |
|            | key |- | - | - | - | - |

### KUBE-CONTROLLER-MANAGER
| Component  | ARG  | Description | Type | CN |  CA |  OU | 
|:-----------|:-----| :------------|:----|:----|:--|:--|
| kube-controller-manager | --root-ca-file | If set, this root certificate authority will be included in service account's token secret. This must be a valid PEM-encoded CA bundle. | CA | - | kubernetes | - |
| | --client-ca-file                       | If set, any request presenting a client certificate signed by one of the authorities in the client-ca-file is authenticated with an identity corresponding to the CommonName of the client certificate. | CA | - | kubernetes | - |
| | --requestheader-client-ca-file         | Root certificate bundle to use to verify client certificates on incoming requests before trusting usernames in headers specified by --requestheader-username-headers. WARNING: generally do not depend on authorization being already done for incoming requests. | CA | - | front-proxy | - |
| | --tls-cert-file                        | File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated after server cert). If HTTPS serving is enabled, and --tls-cert-file and --tls-private-key-file are not provided, a self-signed certificate and key are generated for the public address and saved to the directory specified by --cert-dir. | DigitalSignature,KeyEncipherment,ServerAuth | system:kube-controller-manager | kubernetes | system:kube-controller-manager |
| | --tls-private-key-file                 | File containing the default x509 private key matching --tls-cert-file. | - | - | - | - |
| | --cluster-signing-cert-file            | Filename containing a PEM-encoded X509 CA certificate used to issue cluster-scoped certificates. If specified, no more specific --cluster-signing-* flag may be specified. | DigitalSignature,KeyEncipherment,CertSign | system:kubernetes-sign | kubernetes | - |
| | --cluster-signing-key-file             | Filename containing a PEM-encoded RSA or ECDSA private key used to sign cluster-scoped certificates. If specified, no more specific --cluster-signing-* flag may be specified. | - | - | - | - |
| | --service-account-private-key-file     | Filename containing a PEM-encoded private RSA or ECDSA key used to sign service account tokens. | DigitalSignature,KeyEncipherment | system:kubernetes-sa | kubernetes | - |
| kubeconfig | ca |- | - | - | - | - |
|            | client |- | - | - | - | - |
|            | key |- | - | - | - | - |

### KUBE-SCHEDULER
| Component  | ARG  | Description | Type | CN |  CA |  OU | 
|:-----------|:-----| :------------|:----|:----|:--|:--|
| kube-scheduler | --client-ca-file  | If set, any request presenting a client certificate signed by one of the authorities in the client-ca-file is authenticated with an identity corresponding to the CommonName of the client certificate. | CA | - | kubernetes | - |
|                | --requestheader-client-ca-file | Root certificate bundle to use to verify client certificates on incoming requests before trusting usernames in headers specified by --requestheader-username-headers. WARNING: generally do not depend on authorization being already done for incoming requests. | CA | - | front-proxy | - |
|                | --tls-cert-file | File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated after server cert). If HTTPS serving is enabled, and --tls-cert-file and --tls-private-key-file are not provided, a self-signed certificate and key are generated for the public address and saved to the directory specified by --cert-dir. | DigitalSignature,KeyEncipherment,ServerAuth | system:scheduler-server | kubernetes | system:scheduler-server |
|                | --tls-private-key-file | File containing the default x509 private key matching --tls-cert-file. | - | - | - | - |
| kubeconfig | ca |  | - | - | - | - |
|            | client |- | - | - | - | - |
|            | key |- | - | - | - | - |
