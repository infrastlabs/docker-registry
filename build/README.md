**TODO**

- regisry-minio-clust


**conf**


- bash tls_sanKey.sh test.registry.ssl

```bash
echo 172.25.23.199 test.registry.ssl >> /etc/hosts

certDir=/etc/docker/certs.d/test.registry.ssl:8143
mkdir -p $certDir
cat > $certDir/cert.crt <<EOF
-----BEGIN CERTIFICATE-----
MIICwzCCAaugAwIBAgIUE6s/3VsYmpQu8eIB8FL+sXmCJbAwDQYJKoZIhvcNAQEL
jHrMC7kaV2PxUSJITQJchDd1RppUVorIOFiuQgRQCt0pMxZWObDseFbp646wsH0E
MWm6f86jczuhl3fvd3VsYBPpErhe34mDHBNqEZsu6NMa7k5U5/BO
-----END CERTIFICATE-----
EOF
```

**validate**

```bash
# Error response from daemon: Get https://test.registry.ssl:8143/v2/: x509: certificate is valid for deploy.xx.com.ssl, not test.registry.ssl
echo admin123 |docker login test.registry.ssl:8143 --username=admin --password-stdin

# pushImg
docker pull busybox
docker tag busybox test.registry.ssl:8143/busybox
docker push test.registry.ssl:8143/busybox

# abc/busybox
docker tag busybox test.registry.ssl:8143/abc/busybox:v2
docker push test.registry.ssl:8143/abc/busybox:v2

# host-21-67:/etc/docker/certs.d/test.registry.ssl:8143 # docker push test.registry.ssl:8143/busybox
The push refers to repository [test.registry.ssl:8143/busybox]
01fd6df81c8e: Pushed 
latest: digest: sha256:62ffc2ed7554e4c6d360bce40bbcf196573dd27c4ce080641a2c59867e732dee size: 527
```