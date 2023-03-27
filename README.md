# whereami-go
Golang port of https://github.com/GoogleCloudPlatform/kubernetes-engine-samples/tree/main/whereami

what `whereami-go` doesn't do that the original `whereami` does:
- gRPC (i.e. `whereami-go` is HTTP only)
- OTEL tracing exports (*but* does do trace header propagation so if you're doing tracing via Istio / ASM, that will work because the headers are there)
- doesn't pretty-print the JSON (just run it through `jq` if you want nicer formatting from CLI)