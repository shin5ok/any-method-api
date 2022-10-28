# Set up application
## Prepare environtment values
```
export PROJECT=<your project>
export APPNAME=myapp
export TAG=0.01
```
## Build container and Push it to GCR
Run it in top dir.
```
bash ./build.sh
```

## Normal deploy
```
envsubst < manifests.yaml | kubectl apply -f -
```
Just wait for deploying them completely.

## Just test
```
INGRESS_IP=$(kubectl get ing $APPNAME -o json | jq .status.loadBalancer.ingress[].ip -r)

curl http://$INGRESS_IP/foo/bar
```

## Re-deploy pod with specified latency rate
Add purposeful randomized latency from 500ms to 1000ms before sending responses.
```
# Add purposeful latency to aboud 50% responses
RAND_DIV=2 MODE=sleep envsubst < manifests.yaml | kubectl apply -f -

# Add purposeful latency to all responses
RAND_DIV=1 MODE=sleep envsubst < manifests.yaml | kubectl apply -f -
```
If you specified MODE=error it would return 503 responses following to the RAND_DIV rate.

## Revert to normal application without latency or 503 error
```
RAND_DIV= envsubst < manifests.yaml | kubectl apply -f -
```


# Set up loading client

## Prepare environtment values, just in case.
```
export PROJECT=<your project>
export TAG=0.01
```

## Build container and Push it to GCR
```
cd loading-client/
bash ./build.sh
```

## Run it to load the specified target
Make sure if your current dir is 'loading-client'.
```
IP=$INGRESS_IP envsubst < manifests.yaml | kubectl apply -f -
```

# Configure PodMonitoring resource to collect metrics
change dir to top dir,
```
envsubst < podmonitoring.yaml | kubectl apply -f -
```

