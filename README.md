# Agent Management Sync

POC: Syncs a collection of namespaces following the filesystem structure:

Syncs a collection of namespaces following the filesystem structure:

```
.
└── cfg
    └── <namespace_name>
        ├── base.yaml
        └── snips
            ├── <snip_1_name>.yaml
            ├── <snip_2_name>.yaml
            └── <snip_n_name>.yaml
```

## Usage

```
AGENT_MANAGEMENT_HOST= AGENT_MANAGEMENT_USERNAME= AGENT_MANAGEMENT_PASSWORD= go run ./cmd/main.go
```
