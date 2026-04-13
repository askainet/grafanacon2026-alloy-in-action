# Side Missions: Expanding Your Alloy Utility Belt

Extra missions that show what else you can accomplish with Alloy.

It's recommended to comment out or remove your previous Alloy config before starting these exercises to reduce log noise.

## Debugging components

Alloy includes components specifically for debugging your telemetry pipelines:

- [loki.echo](https://grafana.com/docs/alloy/latest/reference/components/loki/loki.echo/)
- [prometheus.echo](https://grafana.com/docs/alloy/latest/reference/components/prometheus/prometheus.echo/)
- [otelcol.processor.debug](https://grafana.com/docs/alloy/latest/reference/components/otelcol/otelcol.exporter.debug/)

Wire your pipelines to `forward_to` these components and Alloy will print the telemetry to `stdout`.

## Secret filtering

Alloy's [`loki.secretfilter`](https://grafana.com/docs/alloy/latest/reference/components/loki/loki.secretfilter/) component scans incoming log lines and redacts sensitive data like API keys. It uses a [GitLeaks](https://github.com/gitleaks/gitleaks/blob/master/config/gitleaks.toml) config to define which patterns are treated as secrets.

> [!CAUTION]
> `loki.secretfilter` may not catch all secrets and should be used as part of a defense-in-depth strategy. You should still avoid logging secrets.

> [!WARNING]
> `loki.secretfilter` runs regex matching on every log line, which can become a bottleneck in high-throughput pipelines. Keep this in mind when planning production deployments.

`loki.secretfilter` ships with a default gitleaks config, but it can change between Alloy versions. It's best to provide your own so you have full control over what gets matched. You can add as many `[[rules]]` entries as you need.

This repo includes a custom config at `alloy/gitleaks.toml` with a rule that matches `gcon26_` demo tokens:

```toml
[[rules]]
id = "grafanacon-token"
description = "GrafanaCon 2026 Demo Token"
regex = '''gcon26_[a-f0-9]{16,}'''
keywords = ["gcon26_"]
```

### Mission Objectives

- [ ] Configure `loki.source.file` to read from the secrets test file
- [ ] Wire `loki.secretfilter` to filter log lines using the custom gitleaks config
- [ ] Forward filtered output to `loki.echo` to see the results
- [ ] Verify that secrets are redacted in Alloy's stdout

### Build the Pipeline

**Pipeline:**

```
loki.source.file → loki.secretfilter → loki.echo
```

Add this to your `config.alloy` file and fill in the TODOs:

```alloy
// Step 1: Read a file that may contain secrets
loki.source.file "secrets" {
  targets = [{"__path__" = "/etc/alloy/secrets.txt"}]
  file_match {
    enabled = true
  }
  forward_to = [TODO]  // Forward to the secret filter's receiver
}

// Step 2: Filter secrets from incoming log lines
loki.secretfilter "you_shall_not_pass" {
  gitleaks_config = TODO           // Path to your gitleaks config inside the container
  redact_with     = "/\\(ಠ_ಠ) |  < YOU SHALL NOT PASS! >"
  forward_to      = [TODO]         // Forward to loki.echo
}

// Step 3: Print filtered output to stdout
loki.echo "secrets" { }
```

<details>
<summary>Hint 1: file paths</summary>

The Docker setup mounts the `alloy/` directory to `/etc/alloy/` inside the container. So the secrets test file lives at `/etc/alloy/secrets.txt` and the gitleaks config is at `/etc/alloy/gitleaks.toml`.

</details>

<details>
<summary>Hint 2: wiring components together</summary>

Each Alloy component that accepts input exposes a `.receiver` export. To forward from one component to the next, reference `<component_type>.<label>.receiver`.

For example, to forward to `loki.secretfilter "you_shall_not_pass"`, use `loki.secretfilter.you_shall_not_pass.receiver`.

</details>

<details>
<summary>Full solution</summary>

```alloy
loki.source.file "secrets" {
  targets = [{"__path__" = "/etc/alloy/secrets.txt"}]
  file_match {
    enabled = true
  }
  forward_to = [loki.secretfilter.you_shall_not_pass.receiver]
}

loki.secretfilter "you_shall_not_pass" {
  gitleaks_config = "/etc/alloy/gitleaks.toml"
  redact_with     = "/\\(ಠ_ಠ) |  < YOU SHALL NOT PASS! >"
  forward_to      = [loki.echo.secrets.receiver]
}

loki.echo "secrets" { }
```

</details>

### Verify Your Work

Open two terminals. Run `make alloy-logs` in one to watch Alloy's output.

Reload the config:
```bash
make alloy-reload
```

In the other terminal, append a line containing a secret:
```bash
echo "gcon26_a1b2c3d4e5f60718" >> ./alloy/secrets.txt
```

In the Alloy logs terminal, you should see the secret replaced with the redaction string:
```
level=info component_path=/ component_id=loki.echo.secrets receiver=loki.echo.secrets entry="/\\(ಠ_ಠ) |  < YOU SHALL NOT PASS! >" ...
```

Lines that don't match any rule pass through unchanged.

## Metric enrichment with service discovery

Alloy's [`prometheus.enrich`](https://grafana.com/docs/alloy/latest/reference/components/prometheus/prometheus.enrich/) component adds labels from service discovery targets onto incoming metrics. It matches a label on each metric against a label on the discovered targets and copies the specified labels from the matched target onto the metric.

This is useful when your application metrics reference backend services but don't carry infrastructure context — like which team owns the service, what AWS region it runs in, or what instance type it's on. Instead of baking that into every application, you maintain a central service registry and let Alloy enrich metrics at collection time.

```
  Scraped Metric                                             SD Target
 ┌────────────────────────┐          compare                ┌───────────────────────────┐
 │ request_duration{      │ "instance" vs "__address"       │ __address__: "api:8080" ◀── match!
 │   instance="api:8080"  │───────────────────────────────▶ │ team: "payments"          │
 │ } 0.5                  │                                 │ region: "us-east-1"       │
 └────────────────────────┘                                 ├───────────────────────────┤
                                                            │ ...other targets...       │
                                                            └───────────────────────────┘
                                                              │
                                                              ▼
                                                    ┌──────────────────────────────────┐
                                                    │ request_duration{                │
                                                    │   instance="api:8080",           │
                                                    │   team="payments",        ◀── copied from SD target labels
                                                    │   region="us-east-1"      ◀── copied from SD target labels
                                                    │ } 0.5                            │
                                                    └──────────────────────────────────┘
```

`prometheus.enrich` compares the metric's `instance` label (`metrics_match_label`) against each SD target's `__address__` (`target_match_label`). When they match, the specified labels are copied onto the metric. Metrics without a matching label pass through unchanged.

### Scenario

Mission-control monitors its backend dependencies (Grafana, Loki, Tempo, Mimir) with periodic health checks. These produce two metrics:

- `backend_health_checks_total{backend="loki:3100"}` — how many checks have been sent
- `backend_up{backend="loki:3100"}` — whether the backend is reachable (1=up, 0=down)

The `backend` label tells you *what* was checked, but not *who owns it* or *where it runs*. HQ maintains a service registry at `mission-control:8080/api/sd/targets` that returns [Prometheus HTTP SD](https://prometheus.io/docs/prometheus/latest/http_sd/)-formatted targets with AWS infrastructure metadata — account IDs, regions, availability zones, VPC IDs, instance types, and team ownership.

**Your task:** Use [`discovery.http`](https://grafana.com/docs/alloy/latest/reference/components/discovery/discovery.http/) to fetch the service registry and [`prometheus.enrich`](https://grafana.com/docs/alloy/latest/reference/components/prometheus/prometheus.enrich/) to match the `backend` label on scraped metrics against the `__address__` label on the SD targets, copying infrastructure metadata onto each metric before sending to Mimir.

**Before enrichment:**
```
backend_up{backend="loki:3100"} 1
```

**After enrichment:**
```
backend_up{backend="loki:3100", team="sigint", environment="classified", aws_region="us-east-1", availability_zone="us-east-1b", instance_type="r6g.2xlarge"} 1
```

### Mission Objectives

- [ ] Fetch service discovery targets from the mission-control registry using `discovery.http`
- [ ] Scrape health-check metrics from `mission-control:8080`
- [ ] Enrich metrics by matching the `backend` label against SD target `__address__` labels
- [ ] Copy infrastructure labels (`team`, `environment`, `aws_region`, `availability_zone`, `instance_type`) onto matched metrics
- [ ] Send enriched metrics to Mimir and verify them in Grafana

### Build the Pipeline

**Pipeline:**

```
                      discovery.http
                            |
                            v
prometheus.scrape → prometheus.enrich → prometheus.remote_write
```

- Scrape metrics from `mission-control:8080` (the backend health check metrics are exposed at `/metrics`)
- Fetch service discovery targets from `http://mission-control:8080/api/sd/targets`
- Enrich metrics by matching the metric's `backend` label against the SD target's `__address__` label
- Copy `team`, `environment`, `aws_region`, `availability_zone`, and `instance_type` onto each matched metric
- Send enriched metrics to Mimir at `http://mimir:9009/api/v1/push`

### Starter Code

Add this to your `config.alloy` file and fill in the TODOs:

```alloy
/*
  Bonus: Metric Enrichment via Service Discovery
  Pipeline: prometheus.scrape -> prometheus.enrich -> prometheus.remote_write
  With:     discovery.http feeding targets into prometheus.enrich
*/

// Step 1: Fetch infrastructure metadata from the service registry
discovery.http "lgtm_registry" {
  url              = "TODO"  // Where is the service registry?
  refresh_interval = "30s"
}

// Step 2: Scrape metrics from mission-control (includes backend_up, backend_health_checks_total)
prometheus.scrape "mission_control" {
  scrape_interval = "10s"
  targets         = [{"__address__" = "mission-control:8080"}]
  forward_to      = [TODO]  // Forward to the enrichment component's receiver
}

// Step 3: Enrich metrics with infrastructure metadata from the service registry
prometheus.enrich "infra_metadata" {
  targets = TODO  // Use the discovery component's .targets export

  // Match the "backend" label on metrics against "__address__" on SD targets
  target_match_label  = TODO
  metrics_match_label = TODO

  // Only copy the labels we care about (without this, all SD labels including __address__ would be copied)
  labels_to_copy = [TODO]

  forward_to = [TODO]  // Forward to remote_write
}

// Step 4: Send enriched metrics to Mimir
prometheus.remote_write "docker_mimir" {
  endpoint {
    url = "http://mimir:9009/api/v1/push"
  }
}
```

<details>
<summary>Hint 1: discovery.http exports</summary>

`discovery.http` exposes a `.targets` export: a list of target maps containing `__address__` and all labels from the SD response. Reference it as `discovery.http.lgtm_registry.targets`.

</details>

<details>
<summary>Hint 2: how matching works</summary>

`prometheus.enrich` compares the value of `metrics_match_label` on each incoming metric against the value of `target_match_label` on each discovered target. When they match, the specified labels are copied onto the metric.

Here, the scraped metric has `backend="loki:3100"` and the SD target has `__address__="loki:3100"`. So `metrics_match_label = "backend"` and `target_match_label = "__address__"`.

Metrics that don't have a `backend` label (like `http_requests_total`) pass through unchanged.

</details>

<details>
<summary>Full solution</summary>

```alloy
/*
  Bonus: Metric Enrichment via Service Discovery
  Pipeline: prometheus.scrape -> prometheus.enrich -> prometheus.remote_write
  With:     discovery.http feeding targets into prometheus.enrich
*/

// Fetch infrastructure metadata from the service registry
discovery.http "lgtm_registry" {
  url              = "http://mission-control:8080/api/sd/targets"
  refresh_interval = "30s"
}

// Scrape metrics from mission-control (includes backend_up, backend_health_checks_total)
prometheus.scrape "mission_control" {
  scrape_interval = "10s"
  targets         = [{"__address__" = "mission-control:8080"}]
  forward_to      = [prometheus.enrich.infra_metadata.receiver]
}

// Enrich metrics with infrastructure metadata from the service registry
prometheus.enrich "infra_metadata" {
  targets = discovery.http.lgtm_registry.targets

  target_match_label  = "__address__"
  metrics_match_label = "backend"

  labels_to_copy = ["team", "environment", "aws_region", "availability_zone", "instance_type"]

  forward_to = [prometheus.remote_write.docker_mimir.receiver]
}

// Send enriched metrics to Mimir
prometheus.remote_write "docker_mimir" {
  endpoint {
    url = "http://mimir:9009/api/v1/push"
  }
}
```

</details>

### Verify Your Work

Reload the config:
```bash
make alloy-reload
```

Wait about 30 seconds for scrapes to run, then open Explore in Grafana, select Mimir as the data source, and run:
```
backend_up
```

You should see the `backend_up` metric now carrying `team`, `environment`, `aws_region`, `availability_zone`, and `instance_type` labels — infrastructure metadata that didn't exist on the original `/metrics` endpoint. You can now build dashboards that group backend health by team or availability zone, without the application needing to know anything about where it runs.