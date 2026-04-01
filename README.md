# Grafana Alloy in Action

<img width="1059" height="592" alt="image" src="https://github.com/user-attachments/assets/7971d341-29de-4bed-88d7-c419c653ef90" />

<img width="2509" height="1406" alt="image" src="https://github.com/user-attachments/assets/e3a18293-1ffc-4405-8dbe-cd8179227344" />

# Resources 
- [Grafana Alloy documentation](https://grafana.com/docs/alloy/latest/)
  - [Alloy configuration blocks](https://grafana.com/docs/alloy/latest/reference/config-blocks/)
  - [Alloy components](https://grafana.com/docs/alloy/latest/reference/components/)
  - [Collect and forward data with Grafana Alloy](https://grafana.com/docs/alloy/latest/collect/)
  - [Grafana Alloy Tutorials](https://grafana.com/docs/alloy/latest/tutorials/)
    
# What is Alloy? When does it make sense to use it? 

<img width="2557" height="1413" alt="image" src="https://github.com/user-attachments/assets/75f0f830-636e-4c6e-a535-c8895acfd777" />
<img width="2499" height="1399" alt="image" src="https://github.com/user-attachments/assets/98f4e0f9-0623-4733-8bc9-517cf2ace5e9" />
<img width="2504" height="1408" alt="image" src="https://github.com/user-attachments/assets/1ec1c6ee-3126-4df3-8782-0a2f2fd39d23" />
<img width="2506" height="1408" alt="image" src="https://github.com/user-attachments/assets/d38d3741-d8ed-4d0f-95da-b4cb867f3b25" />
<img width="2511" height="1409" alt="image" src="https://github.com/user-attachments/assets/a1c430ac-17db-4325-bc4a-fa554d5b85f4" />
<img width="2508" height="1404" alt="image" src="https://github.com/user-attachments/assets/dbc34673-7dfc-4031-b800-75cff73bcf23" />

# Alloy configuration language 101

### Think of Alloy as our secret weapon that collects, transforms, and exports our telemetry data. 

![Alt Text](https://media0.giphy.com/media/v1.Y2lkPTc5MGI3NjExeTkyOTc1eW1idjM4MW1sbXNtbmZjNzBrcHIyZWo4bXg2dGcwdWNyciZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/Mb9dQnfZXSBYMhU2Nv/giphy.gif)

To instruct Alloy on how we want that done, we must write these instructions in a language (`Alloy syntax`) that Alloy understands. 

<img width="2508" height="1409" alt="image" src="https://github.com/user-attachments/assets/427f1f24-0869-4d6a-87d5-5e91d5b0c9a9" />
<img width="2508" height="1407" alt="image" src="https://github.com/user-attachments/assets/ad4a092f-e914-48f5-aa78-a49b575eb912" />
<img width="2506" height="1409" alt="image" src="https://github.com/user-attachments/assets/9da4ee1f-6d41-4a52-a2d3-352702a30aa5" />
<img width="2508" height="1409" alt="image" src="https://github.com/user-attachments/assets/a9de2552-8654-44ac-ad9d-c6721884603b" />
<img width="2507" height="1408" alt="image" src="https://github.com/user-attachments/assets/3ea214f5-abeb-4798-b9ff-14d8f5c31d34" />
<img width="2509" height="1408" alt="image" src="https://github.com/user-attachments/assets/05b433bb-e80a-4205-9c1b-ade73650f2f0" />
<img width="2510" height="1410" alt="image" src="https://github.com/user-attachments/assets/b1ab4f30-736f-403b-b5a1-eca171d24f02" />
<img width="2510" height="1409" alt="image" src="https://github.com/user-attachments/assets/491e8022-904b-4c90-aac4-c83d0a75d756" />
<img width="2510" height="1404" alt="image" src="https://github.com/user-attachments/assets/aa299e50-4f2e-44d9-adc0-95aabd4ea9ac" />
<img width="2507" height="1409" alt="image" src="https://github.com/user-attachments/assets/1eb9b0c2-2b18-4ea5-a7f0-3035e6112e16" />
<img width="2504" height="1407" alt="image" src="https://github.com/user-attachments/assets/d344bd85-7d5b-4020-b063-554dcca8188f" />
<img width="2511" height="1409" alt="image" src="https://github.com/user-attachments/assets/1ad9451c-40f0-4b4d-bd0e-e4adcd70bb9d" />

The `usage` section gives you an example of how this particular component could be configured. 

<img width="2510" height="1409" alt="image" src="https://github.com/user-attachments/assets/346cdf21-594d-4606-ab55-01dee9470e9a" />

The `arguments` and `blocks` sections list what you could do with the data. Pay close attention to the name, type, description, default, and required columns so Alloy could understand what you want it to do! 

<img width="2508" height="1411" alt="image" src="https://github.com/user-attachments/assets/6a6102f6-82e9-4e31-bff6-faac9a5fcf6b" />

<img width="2509" height="1411" alt="image" src="https://github.com/user-attachments/assets/77f29151-513f-4258-8ad5-6c370ad9265a" />

Focusing on these 3 things will point us in the right direction as we configure our pipeline. 
<img width="2506" height="1411" alt="image" src="https://github.com/user-attachments/assets/eeda8160-d3f6-4b9f-8ce7-52205c460f7e" />

# Learning environment setup

<img width="2506" height="1411" alt="image" src="https://github.com/user-attachments/assets/c188f4a8-4c43-4104-880a-6bf6de1f9ddf" />

Before getting started, make sure you:
- install [Docker Desktop](https://www.docker.com/products/docker-desktop/) and [Docker Compose](https://docs.docker.com/compose/install/)
  
- clone the repo for the learning environment :
```
git clone <UPDATE WITH THE OFFICIAL REPO HERE>
```

To start the environment, run the following command from within the project's root directory:
```
make start     
```

To stop the environment, run the following command from within the project's root directory:
```
make stop
```
Open the project using a text editor of your choice.
- Expand the `alloy/` folder and open the `config.alloy` file.
- We will use this file to build pipelines for the training exercises and missions.

# Foundations 

## Traces

<img width="2501" height="1402" alt="image" src="https://github.com/user-attachments/assets/acf8e2cb-0820-4e3e-a85d-3358a4f1fa0d" />

### Objectives

- Use [`otelcol.receiver.otlp`](https://grafana.com/docs/alloy/latest/reference/components/otelcol/otelcol.receiver.otlp/) to receive traces from the application
- Use [`otelcol.processor.batch`](https://grafana.com/docs/alloy/latest/reference/components/otelcol/otelcol.processor.batch/) to batch traces for efficient export
- Use [`otelcol.exporter.otlphttp`](https://grafana.com/docs/alloy/latest/reference/components/otelcol/otelcol.exporter.otlphttp/) to send traces to Tempo

### Scenario

Field agents are conducting operations around the world, and `mission-control` is tracking every request that flows through the system. 
These operations generate **traces**, detailed records of each request's journey through the application.

Right now, those traces are being generated but going nowhere. Without a pipeline to collect them, you're flying blind.

**Your task:** 
Build a trace pipeline so you can see what's happening inside mission-control.

Once it's configured, you'll be able to view traces in Grafana dashboard and drill into individual operations.

### Build the Pipeline

**Pipeline:** 

`otelcol.receiver.otlp` → `otelcol.processor.batch` → `otelcol.exporter.otlphttp`

- Receive OTLP traces on `0.0.0.0:4317` (gRPC) and `0.0.0.0:4318` (HTTP)
- Batch traces before sending (improves efficiency)
- Export to Tempo at `http://tempo:4318`

### Starter Code
Copy this into your `config.alloy` file and fill in the TODOs:
```alloy
/*
  Section 1: Traces Pipeline
  Pipeline: otelcol.receiver.otlp -> otelcol.processor.batch -> otelcol.exporter.otlphttp
*/

// Step 1: Receive OTLP traces from mission-control

//The receiver listens for incoming traces. Mission-control sends traces using OTLP, so you need both gRPC and HTTP endpoints which are 4317 annd 4318 by default.
otelcol.receiver.otlp "default" {
  grpc {
    endpoint = "0.0.0.0:4317"
  }

  http {
    endpoint = "0.0.0.0:4318"
  }
// OpenTelemetry components use `output` to send data and `.input` to receive it:
  output {
    traces = [TODO]  // Forward to the batch processor's input
  }
}

// Step 2: Batch traces before exporting

otelcol.processor.batch "default" {
  send_batch_size     = 100
  send_batch_max_size = 200
  timeout             = "250ms"

  output {
    traces = [TODO]  // Forward to the Tempo exporter's input
  }
}


// Step 3: Export traces to Tempo
otelcol.exporter.otlphttp "docker_tempo" {
  client {
    endpoint = "TODO"  // Send to http://tempo:4318
    tls {
      insecure             = true
      insecure_skip_verify = true
    }
  }
}
```
<details>
<summary>Full solution</summary>
  
```alloy
/*
  Section 1: Traces Pipeline
  Pipeline: otelcol.receiver.otlp -> otelcol.processor.batch -> otelcol.exporter.otlphttp
*/
// Receive OTLP traces from mission-control
otelcol.receiver.otlp "default" {
  grpc {
    endpoint = "0.0.0.0:4317"
  }
  http {
    endpoint = "0.0.0.0:4318"
  }
  output {
    traces = [otelcol.processor.batch.default.input]
  }
}
// Batch traces before exporting
otelcol.processor.batch "default" {
  output {
    traces = [otelcol.exporter.otlphttp.docker_tempo.input]
  }
}
// Export traces to Tempo
otelcol.exporter.otlphttp "docker_tempo" {
  client {
    endpoint = "http://tempo:4318"
    tls {
      insecure             = true
      insecure_skip_verify = true
    }
  }
}
```
</details>


### Reloading the Config

Whenever you make changes to the config file, you need to reload Alloy:

```bash
make alloy-reload
```

If the config is valid, you'll see:

```
config reloaded
```

### Verify Your Work

1. Open **Grafana** at http://localhost:3000
2. Go go the dasbhaord the **Mission Control Overview** dashboard. The Traces panel at the bottom should now show data.
<img width="2502" height="1407" alt="image" src="https://github.com/user-attachments/assets/000ff34f-3869-4744-a604-211f84587f90" />
<img width="2512" height="1411" alt="image" src="https://github.com/user-attachments/assets/33dc061c-b87f-4c77-b4d5-4947509d5f65" />
3. When you click on one the traces, you can see the full breakdown of what happened and how long it took
<img width="2511" height="1414" alt="image" src="https://github.com/user-attachments/assets/21b3c77c-c918-4e77-bb7f-b7969da3485a" />

## Logs

### Objectives

<img width="2503" height="1409" alt="image" src="https://github.com/user-attachments/assets/d5f91a67-44db-4b4a-8bcb-b0c459cb8e85" />

- Use [`local.file_match`](https://grafana.com/docs/alloy/latest/reference/components/local/local.file_match/) to discover log files
- Use [`loki.source.file`](https://grafana.com/docs/alloy/latest/reference/components/loki/loki.source.file/) to tail and read log files
- Use [`loki.process`](https://grafana.com/docs/alloy/latest/reference/components/loki/loki.process/) to parse JSON and extract labels
- Use [`loki.write`](https://grafana.com/docs/alloy/latest/reference/components/loki/loki.write/) to send logs to Loki

### Scenario

Mission-control generates operational logs: JSON-formatted records of everything happening in the system. Agent check-ins, status updates, internal events. It's all there.

These logs are written to files inside the container, but no one's reading them. Without a pipeline to collect and ship them, you're missing critical intel.

**Your task:** 
Build a log pipeline to tail mission-control's log files, parse the JSON, and send them to Loki.

Once it's configured, you'll see logs flowing into Grafana where you can search and filter by component.

### Build the Pipeline

**Pipeline:** 

`local.file_match` → `loki.source.file` → `loki.process` → `loki.write`

- Match log files at `/var/log/alloy/*.log`
- Read log files and forward log lines
- Parse JSON and extract the `component` field as a label
- Write to Loki at `http://loki:3100/loki/api/v1/push`

### Starter Code

Add this to your `config.alloy` file (below your traces pipeline) and fill in the TODOs:

```alloy
/*
  Section 2: Logs Pipeline
  Pipeline: local.file_match -> loki.source.file -> loki.process -> loki.write
*/

// Step 1: Find log files to read

local.file_match "mission_control" {
  path_targets = [{"__path__" = "/var/log/alloy/*.log"}]
}

// Step 2: Read the log files

loki.source.file "mission_control" {
  targets    = TODO  // hint: component_type.label.targets
  forward_to = [TODO]  // Forward to loki.process receiver
}

/*
Step 3: Extract "component" from JSON and promote to label
Original log: {"level":"INFO","component":"agents","msg":"Request received"}
Labels are indexed — makes filtering fast in Grafana
*/

loki.process "mission_control_logs" {
//`stage.json` parses the log body
  stage.json {
    expressions = {
      component = "component",    // Extract the "component" field from the JSON log body
    }
  }

// Add the extracted field as a Loki label
//`stage.labels` adds extracted fields to Loki labels
  stage.labels {
    values = {
      component = "",  // Empty string means "use the extracted value with the same name."
    }
  }

  forward_to = [TODO]  // Forward to the loki.write receiver
}

// Step 4: Send logs to Loki

loki.write "docker_loki" {
  endpoint {
    url = "http://loki:3100/loki/api/v1/push"
  }
}
```

<details>
<summary>Full solution</summary>

```alloy
/*
  Section 2: Logs Pipeline
  Pipeline: local.file_match -> loki.source.file -> loki.process -> loki.write
*/

// Find log files to read
local.file_match "mission_control" {
  path_targets = [{"__path__" = "/var/log/alloy/*.log"}]
}

// Read the log files
loki.source.file "mission_control" {
  targets    = local.file_match.mission_control.targets
  forward_to = [loki.process.mission_control_logs.receiver]
}

// Parse JSON logs and extract labels
loki.process "mission_control_logs" {
  stage.json {
    expressions = {
      component = "component",
    }
  }

  stage.labels {
    values = {
      component = "",
    }
  }

  forward_to = [loki.write.docker_loki.receiver]
}

// Send logs to Loki
loki.write "docker_loki" {
  endpoint {
    url = "http://loki:3100/loki/api/v1/push"
  }
}
```

</details>

### Reloading the Config

Whenever you make changes to the config file, you need to reload Alloy:

```bash
make alloy-reload
```

If the config is valid, you'll see:

```
config reloaded
```

### Verify Your Work

Check the **Mission Control Overview** dashboard. The logs panel should now show data.

<img width="2507" height="1407" alt="image" src="https://github.com/user-attachments/assets/3f1eb1ba-987b-4418-822c-b8210c408bc1" />

Next, open Explore from the Grafana sidebar. 

Explore is Grafana's query playground. It lets you run ad-hoc queries without building a dashboard.

Select Loki as the data source and try this query:
{component="agents"}

You should see logs filtered to just field agent activity. 
<img width="2509" height="1405" alt="image" src="https://github.com/user-attachments/assets/ce5b5337-78da-4eb1-b17e-0859435b046b" />
Check Common labels at the top. It shows component=agents. That's the label we extracted from the JSON.

Without that label, you'd have to search the entire log body to filter by component. With it, you can filter instantly.

## Metrics

<img width="2510" height="1409" alt="image" src="https://github.com/user-attachments/assets/26bd91f6-50af-41cf-9868-1c493b91edeb" />

### Objectives

- Use [`prometheus.scrape`](https://grafana.com/docs/alloy/latest/reference/components/prometheus/prometheus.scrape/) to scrape metrics from the application
- Use [`prometheus.relabel`](https://grafana.com/docs/alloy/latest/reference/components/prometheus/prometheus.relabel/) to normalize label names for consistency
- Use [`prometheus.remote_write`](https://grafana.com/docs/alloy/latest/reference/components/prometheus/prometheus.remote_write/) to export metrics to Mimir

### Scenario

You're the operations team responsible for monitoring field agents deployed around the world. Your agents check in regularly, and you need to track who's active and where they're located.

Those check-ins are captured as metrics. The `mission-control` app exposes them at `/metrics`, but right now nothing is collecting them.

**Your task:** 
Build a metrics pipeline to scrape mission-control, normalize the labels, and send them to Mimir.

Once it's configured, you'll see a world map showing agent locations, grouped by country.

### The Problem

Two metrics track your agents: `active_agents` tells you who's currently online, and `agent_comms_total` counts their check-ins. Together, they paint a picture of field activity.

But these metrics label the same data differently. `active_agents` uses `agent_id` and `country_code`, while `agent_comms_total` uses `id` and `region`. This makes it harder to correlate metrics later. You'll use `prometheus.relabel` to normalize them.

### Build the Pipeline

**Pipeline:** 

`prometheus.scrape` → `prometheus.relabel` → `prometheus.remote_write`

- Scrape metrics from `mission-control:8080`
- Rename `id` → `agent_id` and `region` → `country_code`
- Send to Mimir at `http://mimir:9009/api/v1/push`

### Starter Code

Add this to your `config.alloy` file (below your logs pipeline) and fill in the TODOs:

```alloy
/*
  Section 3: Metrics Pipeline
  Pipeline: prometheus.scrape -> prometheus.relabel -> prometheus.remote_write
*/

// Step 1: Scrape metrics from mission-control
prometheus.scrape "mission_control" {
  scrape_interval = "5s"
  scrape_timeout  = "4s"

  targets    = [{"__address__" = "mission-control:8080"}]
  forward_to = [TODO]  // Forward to prometheus.relabel.normalize_agent_labels receiver
}

// Step 2: Normalize label names

/*
  Normalize labels so all agent metrics use consistent naming

  Before:
    active_agents{agent_id="ALPHA-007", country_code="US"}
    agent_comms_total{id="ALPHA-007", region="US"}

  After:
    active_agents{agent_id="ALPHA-007", country_code="US"}
    agent_comms_total{agent_id="ALPHA-007", country_code="US"}

  Rename: id → agent_id, region → country_code
*/

prometheus.relabel "normalize_agent_labels" {
  rule {
    action        = "replace"
    source_labels = ["id"]
    regex         = "(.+)"
    target_label  = "agent_id"
  }

  rule {
    action = "labeldrop"
    regex  = "^id$"
  }

 // We’ll do the same with the “region” label to rename it to “country_code” and drop the old label

  rule {
    action        = "replace"
    source_labels = ["TODO"]
    regex         = "(.+)"
    target_label  = "TODO"
  }

  rule {
    action = "labeldrop"
    regex  = "^TODO$"
  }
  forward_to = [TODO]  // Forward to prometheus.remote_write receiver
}

// Step 3: Send metrics to Mimir
prometheus.remote_write "docker_mimir" {
  endpoint {
    url = "http://mimir:9009/api/v1/push"
  }
}

```

// Step 2: Normalize label names

<details>
<summary>Full solution</summary>

```alloy
/*
  Section 3: Metrics Pipeline
  Pipeline: prometheus.scrape -> prometheus.relabel -> prometheus.remote_write
*/

// Scrape metrics from mission-control
prometheus.scrape "mission_control" {
  scrape_interval = "5s"
  scrape_timeout  = "4s"

  targets    = [{"__address__" = "mission-control:8080"}]
  forward_to = [prometheus.relabel.normalize_agent_labels.receiver]
}

// Step 2: Normalize label names

/*
  Normalize labels so all agent metrics use consistent naming

  Before:
    active_agents{agent_id="ALPHA-007", country_code="US"}
    agent_comms_total{id="ALPHA-007", region="US"}

  After:
    active_agents{agent_id="ALPHA-007", country_code="US"}
    agent_comms_total{agent_id="ALPHA-007", country_code="US"}

  Rename: id → agent_id, region → country_code
*/

prometheus.relabel "normalize_agent_labels" {
  rule {
    action        = "replace"
    source_labels = ["id"]
    regex         = "(.+)"
    target_label  = "agent_id"
  }

  rule {
    action = "labeldrop"
    regex  = "^id$"
  }

  rule {
    action        = "replace"
    source_labels = ["region"]
    regex         = "(.+)"
    target_label  = "country_code"
  }

  rule {
    action = "labeldrop"
    regex  = "^region$"
  }
  forward_to = [prometheus.remote_write.docker_mimir.receiver]  // Forward to prometheus.remote_write receiver
}

// Step 3: Send metrics to Mimir
prometheus.remote_write "docker_mimir" {
  endpoint {
    url = "http://mimir:9009/api/v1/push"
  }
}
```

</details>

### Reloading the Config

Whenever you make changes to the config file, you need to reload Alloy:

```bash
make alloy-reload
```

If the config is valid, you'll see:

```
config reloaded
```

### Verify Your Work

Check the **Mission Control Overview** dashboard and view the Active Agents panel. 
<img width="2508" height="1409" alt="image" src="https://github.com/user-attachments/assets/6088be81-95ea-427a-8f0d-2d2ff41b7660" />

To verify your relabel rules are working, go to Explore, select Mimir as a data source and run this query:
```
agent_comms_total{agent_id=~".+"}
```
You should see that agent_comms_total now has agent_id and country_code labels (purple box). That confirms that our relabeling rules worked. 
<img width="2506" height="1410" alt="image" src="https://github.com/user-attachments/assets/5fd78140-ef2a-4763-9058-67063bf62205" />


### How to use the Alloy UI to debug pipelines

Alloy UI is a useful tool that helps you visualize how Alloy is configured and what it is doing so you are able to debug efficiently. 

Navigate to `localhost:12347` to see the list of components (orange box) that alloy is currently configured with.
Click on the blue ‘view’ button on the right side (red arrow).
<img width="2506" height="1402" alt="image" src="https://github.com/user-attachments/assets/3d2c591c-2519-4aed-8f8b-7072b7a0dd91" />

This page shows us the health of the component, the arguments it's using, and its current exports (green box). 

This page also gives us quick access to the component’s documentation (orange arrow) and a Live Debugging view (yellow arrow). 
<img width="2501" height="1407" alt="image" src="https://github.com/user-attachments/assets/9ce8d733-7f14-4686-b128-248f2a0faa37" />

When we click on the Live Debugging view, we will be able to see a real-time stream of telemetry flowing through a component.
<img width="2509" height="1411" alt="image" src="https://github.com/user-attachments/assets/ae97f9a3-574a-4d15-a67c-5667e2dbca40" />

Navigate to the ‘Graph’ tab to access the graph of components and how they are connected.
<img width="2506" height="1409" alt="image" src="https://github.com/user-attachments/assets/e0960d65-efcc-43b3-9681-2cdd83a8ad62" />

The number (pink box) shown on the dotted lines shows the rate of transfer between components. The window at the top (orange box) configures the interval over which alloy should calculate the per-second rate, so a window of ‘10’ means that alloy should look over the last 10 seconds to compute the rate.

The color of the dotted line signifies what type of data are being transferred between components. See the color key (purple box) for clarification. 

**To debug the piplines using the Alloy UI**
- Ensure that no component is reported as unhealthy.
- Ensure that the arguments and exports for misbehaving components appear correct.
- Ensure that the live debugging data meets your expectations.

# BREAK

# Missions 

Training's over. Each mission throws a new crisis at the pipelines you've built. Activate with `make missionN`, reset all with `make reset`.

## Mission I: Rogue Dimension

```bash
make mission1
```
<img width="2513" height="1411" alt="image" src="https://github.com/user-attachments/assets/418058d0-ff9c-41e6-b823-a3fd25b624de" />

An adversary discovered that our server records the full request path as a metric label. They're now flooding us with requests to thousands of random URLs — paths like `/api/a3f8c2e1` that don't map to any real endpoint. Every unique path creates a new time series in Mimir, and cardinality is climbing fast.

**Your orders:** Using what you learned in Foundation III, fetch the list of legitimate paths with `remote.http` from `http://mission-control:8080/api/metrics/allowed-paths`, then use `prometheus.relabel` with a `keep` action on the `path` label to filter out the garbage. Only real routes should make it through to Mimir.

<img width="2504" height="1407" alt="image" src="https://github.com/user-attachments/assets/6e3da07a-2b2f-47d8-b673-0bd6a0636860" />

### Starter Code

Add these components to your `config.alloy` and fill in the TODOs. You'll also need to update `normalize_agent_labels`'s `forward_to` to route through the new relabel component instead of going directly to `remote_write`.

```alloy
// Step 1: Fetch the allowlist of legitimate paths
remote.http "allowed_paths_regex" {
  url = "http://mission-control:8080/api/metrics/allowed-paths"
}

// Step 2: Filter — only keep metrics where "path" is a known legitimate route
prometheus.relabel "mission1" {
  rule {
    action        = "keep"
    source_labels = ["path"]
    regex         = TODO
  }

  forward_to = [prometheus.remote_write.docker_mimir.receiver]
}

// Step 3: Update normalize_agent_labels forward_to:
// Before: forward_to = [prometheus.remote_write.docker_mimir.receiver]
// After:  forward_to = [prometheus.relabel.mission1.receiver]
```

<details>
<summary>Hint 1: parsing the response</summary>

The response body is JSON. Alloy has a standard library with functions for encoding and decoding — check the [encoding functions](https://grafana.com/docs/alloy/latest/reference/stdlib/encoding/).

</details>

<details>
<summary>Hint 2: accessing the response body</summary>

`remote.http` exposes the fetched response body as `.content`. You can use it directly in expressions elsewhere in your config.

</details>

<details>
<summary>Hint 3: putting it together</summary>

`encoding.from_json()` parses a JSON string into an object. Fields on the resulting object are accessible with dot notation — the same way you reference component exports.

</details>

### Verify Your Work

```bash
make mission1-verify
```

## Mission II: Operation Cold Storage

```bash
make mission2
```
<img width="2509" height="1410" alt="image" src="https://github.com/user-attachments/assets/001ea468-26d0-46b2-bca2-4d94fd96c680" />

After the last incident, the higher-ups want us collecting the DEBUG logs we were previously dropping. Turns out those include request logs that could have helped us track down the attacker. But pumping everything into Loki would blow the budget. The directive: archive all logs to S3 for auditing, but only send `INFO`/`WARN`/`ERROR` to Loki for fast queries.

**Your orders:** The skills you picked up in Foundation II will come in handy here. Split your log pipeline into two parallel paths:
1. **All logs** -> S3 via `otelcol.receiver.loki` -> `otelcol.processor.batch` -> `otelcol.exporter.awss3`
2. **Non-DEBUG only** -> Loki via a second `loki.process` with `stage.drop`

<img width="2511" height="1411" alt="image" src="https://github.com/user-attachments/assets/ce6e21b3-4c1e-41db-b9ed-a1450fef230a" />

### Starter Code

Extend your existing log pipeline. Update `loki.process.mission_control_logs`'s `forward_to` to fan out to two destinations, then add the components below and fill in the TODOs.

```alloy
// Step 1: Update loki.process.mission_control_logs forward_to:
// forward_to = [
//   TODO,  // all logs -> S3 path
//   TODO,  // non-DEBUG -> Loki path
// ]

// Path 1: Bridge Loki logs to OTLP format for S3 export
otelcol.receiver.loki "all_logs" {
  output {
    logs = [TODO]
  }
}

otelcol.processor.batch "s3_logs" {
  timeout             = "10s"
  send_batch_size     = 100
  send_batch_max_size = 200

  output {
    logs = [TODO]
  }
}

otelcol.exporter.awss3 "audit_logs" {
  s3_uploader {
    s3_bucket = TODO
    s3_prefix = "logs"

    endpoint            = "http://localstack:4566"
    disable_ssl         = true
    s3_force_path_style = true
  }

  marshaler {
    type = "body"
  }

  sending_queue {
    batch {
      flush_timeout = "500ms"
      min_size      = 100
      sizer         = "items"
    }
  }
}

// Path 2: Filter logs before sending to Loki
loki.process "filter_debug" {
  stage.json {
    expressions = {
      level = TODO,  // which JSON field contains the log level?
    }
  }

  // TODO: add a stage to drop the logs you don't want forwarded to Loki

  forward_to = [loki.write.docker_loki.receiver]
}
```

<details>
<summary>Hint 1: wiring the fan-out</summary>

Loki and OTLP components use the same `.receiver` pattern you've already seen with Prometheus. `otelcol.receiver.loki` exposes a `.receiver` export that accepts Loki log entries.

</details>

<details>
<summary>Hint 2: filtering logs</summary>

Check the [`loki.process` docs](https://grafana.com/docs/alloy/latest/reference/components/loki/loki.process/) for stages that can filter log lines based on an extracted field value.

</details>

<details>
<summary>Hint 3: stage details</summary>

The stage takes `source` (the extracted field to compare) and `value` (the string to match against).

</details>

### Verify Your Work

```bash
make mission2-verify
```

## Mission III: Selective Surveillance

```bash
make mission3
```
<img width="2506" height="1410" alt="image" src="https://github.com/user-attachments/assets/7e06835d-2fe1-4ae2-87af-bb9623913bd5" />

We need to keep our network traffic to a minimum to maintain a low profile. Right now we're sending every single trace to Tempo, and that kind of volume is going to attract attention.

**Your orders:** Go back to the pipeline you built in Foundation I and add head sampling to cut the volume down. Insert `otelcol.processor.probabilistic_sampler` between the OTLP receiver and batch processor. Set `sampling_percentage = 5.0` to keep only 5% of traces. While you're at it, think about what we're trading away here. What intelligence might slip through the cracks?
<img width="2502" height="1405" alt="image" src="https://github.com/user-attachments/assets/ba811483-6c22-4aec-b5cf-b4ab7700b85a" />

### Starter Code

Add this component to your `config.alloy` and fill in the TODOs. You'll also need to update `otelcol.receiver.otlp.default`'s output to route through it.

```alloy
// Update otelcol.receiver.otlp "default" output:
// Before: traces = [otelcol.processor.batch.default.input]
// After:  traces = [TODO]

otelcol.processor.probabilistic_sampler "mission3" {
  sampling_percentage = TODO

  output {
    traces = [TODO]
  }
}
```

<details>
<summary>Hint: wiring</summary>

OTLP components use `.input` to receive data, just like you wired the batch processor in the foundations.

</details>

### Verify Your Work

```bash
make mission3-verify
```

## Mission IV: Leave No (Error) Trace

```bash
make mission4
```
<img width="2510" height="1410" alt="image" src="https://github.com/user-attachments/assets/54db1bff-6bab-4403-b7ae-52e7ea1b3078" />

A field agent has sent us an encrypted dead drop and is transmitting the decryption key in fragments through span attributes on error traces. One piece every 15 seconds, five pieces total. Head sampling would toss most of them before we ever get the chance to reassemble the key.

**Your orders:** Swap out the head sampler for `otelcol.processor.tail_sampling`. Configure two policies:
1. `status_code` policy: keep all error traces (every fragment counts)
2. `probabilistic` policy: sample 5% of everything else

Once you've recovered all the fragments, piece together the token and crack open the dead drop.

```bash
make access-token                 # Check token fragment recovery (wait ~75s)
make deaddrop KEY="your-token"    # Unlock the dead drop with the assembled token
```
<img width="2508" height="1410" alt="image" src="https://github.com/user-attachments/assets/bbd15e2c-1971-44b2-b5cf-59a00a3160c0" />

### Starter Code

Replace `otelcol.processor.probabilistic_sampler.mission3` with this component and fill in the TODOs. Update `otelcol.receiver.otlp.default`'s output to route through it instead.

```alloy
// Update otelcol.receiver.otlp "default" output:
// Before: traces = [otelcol.processor.probabilistic_sampler.mission3.input]
// After:  traces = [TODO]

otelcol.processor.tail_sampling "mission4" {
  decision_wait = "10s"

  policy {
    name = "keep_all_errors"
    type = TODO

    // TODO: add the inner block for this policy type
  }

  policy {
    name = "sample_normal_traffic"
    type = TODO

    // TODO: add the inner block for this policy type
  }

  output {
    traces = [otelcol.processor.batch.default.input]
  }
}
```

<details>
<summary>Hint 1: policy types</summary>

Check the [`otelcol.processor.tail_sampling` docs](https://grafana.com/docs/alloy/latest/reference/components/otelcol/otelcol.processor.tail_sampling/) for the available policy types. You need one that matches on span status, and one that samples by percentage.

</details>

<details>
<summary>Hint 2: inner block structure</summary>

Each policy type has a corresponding inner block with the same name as the type, containing that type's arguments. For example:

```alloy
policy {
  name = "example"
  type = "latency"

  latency {
    threshold_ms = 5000
  }
}
```

</details>

<details>
<summary>Hint 3: error status codes</summary>

The [`status_code` block](https://grafana.com/docs/alloy/latest/reference/components/otelcol/otelcol.processor.tail_sampling/#status_code) lists the valid values for `status_codes`.

</details>

### Verify Your Work

```bash
make mission4-verify
```

## All Commands

```bash
# Environment
make start / stop / restart / clean

# Monitoring
make logs / alloy-logs / metrics / status

# Missions
make mission1 / mission2 / mission3 / mission4 / reset
make mission1-verify / mission2-verify / mission3-verify / mission4-verify

# Mission 4 endgame
make access-token
make deaddrop KEY="..."
```

## Troubleshooting

- **Alloy not receiving traces?** Check `docker compose logs alloy` for connection errors
- **Port conflicts?** Check ports 8080, 3000, 3100, 3200, 4317, 4318, 9009, 12347
