# jellyfin systeminfo proxy

This proxy is designed to fix chromecast network issue with Jellyfin.

## Installation

Deploy jellyfin-proxy helm chart.

Override system info path for an ingress in jellyfin helm chart.

```
# ingress values.yaml
...
paths:
  - path: /System/Info/Public
    pathType: Exact
    service:
      name: jellyfin-proxy
      port: 8080
  - path: /
    pathType: "Prefix"
    service:
      name: jellyfin
      port: 8096
...
```


