# Agent Management Sync

POC: Syncs a collection of namespaces following the filesystem structure:

Syncs a collection of namespaces following the filesystem structure:

```
<config_path>
├── <namespace_1_name>
│   ├── base.yaml
│   └── snips
│       ├── <snip_1_name>.yaml
│       └── <snip_n_name>.yaml
└── <namespace_n_name>
    └── ...


```

## Usage

```
CONFIG_PATH= AGENT_MANAGEMENT_HOST= AGENT_MANAGEMENT_USERNAME= AGENT_MANAGEMENT_PASSWORD= go run ./cmd/main.go
```
