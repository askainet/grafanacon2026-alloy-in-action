# Side Missions: Expanding Your Alloy Utility Belt

Extra missions that show what else you can accomplish with Alloy.

It's recommended to comment out or remove your previous Alloy config before starting these exercises to reduce log noise.

> **Note:** If you completed the main missions, Mimir may have hit its series limit from the cardinality explosion in Mission 1. If you see "No data" when querying new metrics, recreate Mimir to clear the old data:
> ```bash
> docker compose down mimir && docker compose up -d mimir
> ```
> Wait about 30 seconds after recreating before querying again.

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

Open two terminals, both in the **project root directory** (`grafanacon2026-alloy-in-action/`).

**Terminal 1 - Alloy logs:** This streams the Alloy container's stdout so you can see what `loki.echo` prints. Run:
```bash
make alloy-logs
```
This is equivalent to `docker compose logs alloy -f`. Leave it running.

**Terminal 2 - Commands:** Reload the config, then write a test secret:
```bash
make alloy-reload
```

Then append a line containing a secret:
```bash
echo "gcon26_a1b2c3d4e5f60718" >> ./alloy/secrets.txt
```

Back in **Terminal 1** (the Alloy logs), you should see the secret replaced with the redaction string:
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
  target_match_label  = "TODO"
  metrics_match_label = "TODO"

  // Only copy the labels we care about (without this, all SD labels including __address__ would be copied)
  labels_to_copy = ["TODO", "TODO", "TODO", "TODO", "TODO"]

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
<summary>Hint 1: service registry URL</summary>

The service registry is at `http://mission-control:8080/api/sd/targets`. You can curl it from your terminal (using `localhost`) to see what it returns:
```bash
curl -s http://localhost:8080/api/sd/targets | jq .
```

</details>

<details>
<summary>Hint 2: discovery.http exports</summary>

`discovery.http` exposes a `.targets` export: a list of target maps containing `__address__` and all labels from the SD response. Reference it as `discovery.http.lgtm_registry.targets`.

</details>

<details>
<summary>Hint 2: how matching works</summary>

`prometheus.enrich` compares the value of `metrics_match_label` on each incoming metric against the value of `target_match_label` on each discovered target. When they match, the specified labels are copied onto the metric.

Here, the scraped metric has `backend="loki:3100"` and the SD target has `__address__="loki:3100"`. So `metrics_match_label = "backend"` and `target_match_label = "__address__"`.

Metrics that don't have a `backend` label (like `http_requests_total`) pass through unchanged.

</details>

<details>
<summary>Hint 3: available labels</summary>

Curl the service registry to see what labels each target has:
```bash
curl -s http://localhost:8080/api/sd/targets | jq .
```

The labels you see in the response (like `team`, `aws_region`, etc.) are what you can list in `labels_to_copy`.

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

## Running the OpenTelemetry Engine

Everything you've built so far uses Alloy's default engine, which reads configuration written in [Alloy syntax](https://grafana.com/docs/alloy/latest/get-started/configuration-syntax/) (the `config.alloy` file you've been editing). As of [Alloy v1.14.0](https://grafana.com/blog/native-opentelemetry-inside-alloy-now-you-can-get-the-best-of-both-worlds/), Alloy also ships with an [experimental OpenTelemetry engine](https://grafana.com/docs/alloy/latest/introduction/otel_alloy/) that lets you run Alloy as a fully compatible OTel Collector, configured with standard upstream collector YAML, while retaining access to Alloy's features and integrations.

> [!CAUTION]
> The OpenTelemetry engine is **experimental** and, per the Alloy docs, [is subject to frequent breaking changes and may be removed with no equivalent replacement](https://grafana.com/docs/alloy/latest/introduction/otel_alloy/). Make sure you understand the risks before using an experimental feature in your environments.

### Default engine vs. OpenTelemetry engine

| | Default engine | OpenTelemetry engine |
|---|---|---|
| Configuration | Alloy's native syntax | Standard upstream OTel Collector YAML |
| Component focus | Alloy's native components, including first-class support for Prometheus-native workflows (scraping, service discovery, remote_write) | All [OpenTelemetry Collector core](https://github.com/open-telemetry/opentelemetry-collector) components plus a curated selection from [contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib), enabling native OpenTelemetry pipelines end-to-end |
| Stability | Stable default, backward compatibility guarantees | Experimental, operational parity still being developed |
| Can coexist? | Yes, can run alongside the OTel engine via the [Alloy Engine extension](https://grafana.com/docs/alloy/latest/set-up/otel_engine/) | Same |

Adopting the OTel engine is optional and fully backwards compatible. Existing Alloy configurations keep working unchanged unless you opt in.

> [!NOTE]
> Default ports and logging formats may differ slightly between the two engines, so expect some operational differences when switching.

> [!TIP]
> If you need a component that isn't in Alloy's bundled set, you can build a custom Alloy binary with your own component selection using the [OpenTelemetry Collector Builder (OCB)](https://github.com/open-telemetry/opentelemetry-collector/tree/main/cmd/builder). Official docs on doing this with Alloy are coming soon.

### Mission Objectives

- [ ] Translate the Foundation 1 traces pipeline from Alloy syntax into an upstream OTel Collector YAML config
- [ ] Swap the Alloy container over to the OTel engine and verify traces still reach Tempo

### Build the Pipeline

**Pipeline (same shape as Foundation 1):**

```
otlp receiver -> batch processor -> otlphttp exporter
```

The building blocks are the same ones you already know, just expressed in upstream OTel Collector YAML. Mechanical differences to keep in mind:

- Pipelines are declared under `service.pipelines`, not wired through `output` / `forward_to`
- Each component reference uses a type and an optional instance name separated by `/` (for example, `otlphttp/tempo`)
- It's YAML, so indentation is significant

### Starter Code

Try translating the Alloy config from the Traces foundation mission into OTel Collector YAML. There's starter code below if you'd like, or try doing it from scratch if you're familiar with Collector configs ([configuration reference](https://opentelemetry.io/docs/collector/configuration/)).

Create a new file at [alloy/otel-config.yaml](alloy/otel-config.yaml) and fill in the TODOs:

```yaml
# Foundation 1 traces pipeline, expressed as an OTel Collector config

receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

processors:
  batch:
    send_batch_size: 100
    send_batch_max_size: 200
    timeout: 250ms

exporters:
  otlp_http/tempo:
    endpoint: TODO  # Send to http://tempo:4318
    tls:
      insecure: true
      insecure_skip_verify: true

service:
  pipelines:
    traces:
      receivers: [TODO]   # Which receiver feeds this pipeline?
      processors: [TODO]  # Which processor should traces flow through?
      exporters: [TODO]   # Which exporter ships traces to Tempo?
```

<details>
<summary>Hint: referencing components in the pipeline</summary>

Each entry in `service.pipelines.*.{receivers,processors,exporters}` is the component's identifier from above. Use the bare type (`otlp`, `batch`) for components without an instance name, or `type/name` (like `otlphttp/tempo`) for ones that have one.

</details>

<details>
<summary>Full solution</summary>

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

processors:
  batch:
    send_batch_size: 100
    send_batch_max_size: 200
    timeout: 250ms

exporters:
  otlp_http/tempo:
    endpoint: http://tempo:4318
    tls:
      insecure: true
      insecure_skip_verify: true

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlphttp/tempo]
```

</details>

### Run the OTel engine

To see traces flow through the OTel engine, swap the Alloy container's command. Both engines bind the same OTLP ports (`4317` and `4318`), so they can't both handle the same ports simultaneously without extra configuration.

1. In [docker-compose.yml](docker-compose.yml), replace the `alloy` service's `command:` block with:

   ```yaml
   command:
     - otel
     - run
     - --config=file:/etc/alloy/otel-config.yaml
   ```

2. Recreate the container:

   ```bash
   docker compose up -d alloy
   ```

3. Open the [Mission Control Overview dashboard](http://localhost:3000/d/mission-control-overview/mission-control-overview) and confirm the Traces panel has recent traces. Mission Control is still sending OTLP traces to `alloy:4318`, and the OTel engine is now receiving and exporting them to Tempo using the YAML config you just wrote.

When you're done experimenting, you can revert the `command:` block in `docker-compose.yml` to the original `run /etc/alloy/config.alloy` form and `docker compose up -d alloy` again to return to the default engine.