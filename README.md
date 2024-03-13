# LayerZero Takehome
This is a simple golang service to fetch the price of a cryptocurrency (default bitcoin) in fiat currency (default CNY). It fetches the price two days ago and now and makes these available on port 8080.

## Running Docker Image 🐳

1. Build Docker image:
```
make build
```
2. Run with default values (BTC & CNY)

```
make run
```

3. Alternatively, run with specify different currencies

```
docker run -p 8080:8080 -e MAIN_CURRENCY=<YOUR CRYPTOCURRENCY> -e VS_CURRENCY=<YOUR FIAT CURRENCY> currency-price-checker
```

## Building in K8S ☸️
These steps explain how to run this docker image on kind. For different providers, the steps would differ: the nginx controller would not be the same, and the docker image would be sourced from an image registry like GCR.


1. create kind cluster (from their [doc](https://kind.sigs.k8s.io/docs/user/ingress/#create-cluster))
```
cat <<EOF | kind create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
EOF
```

2. Load docker image inside the cluster:
```
make kind_load_image
```

3. Install nginx ingress controller:
```
make install_nginx
```

4. Wait for the nginx controller to be ready, then install the helm chart:
```
helm install currency-price-checker charts/ -f charts/values.yaml
```

The prices should now be available with:
```
curl http://pricechecker.com/
```
🥳🥳

*Requirement*:
You would need to modify your local host file (/etc/host) to resolve pricechecker.com locally
```
127.0.0.1 pricechecker.com
```

## Secrets 🤐
The coingecko API does not require any token on this endpoint, so I did not bake any secret into the pods. However, I left the blueprint to use [bitnami's sealed secrets](https://github.com/bitnami-labs/sealed-secrets) if we needed to add secrets.

Keeping track of Kubernetes secrets can be very messy. Sealed Secrets is an amazing solution to keep secrets encrypted inside version control.
To add sealedsecrets here we would need to install the SS controller and install the kubeseal cli for encryption.

## Potential Improvements 🏗️
- For this service to be fully scalable we would need to create a horizontal pod autoscaler (HPA) that monitors resource use in our deployment and spins new replicas to match our needs.
- If this service were to be exposed to the public internet, it would need to be put behind a load balancer and a bot protection service like cloudfront or akamai to limit rate and ban bad actors.