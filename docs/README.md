# Mittens Documentation# Mittens Documentation# Mittens Documentation# Kubetap Documentation



Mittens is a kubectl plugin for intercepting and inspecting traffic to Kubernetes Services using mitmproxy. It deploys an mitmproxy sidecar container next to your target Service, allowing you to see and modify traffic in real-time without port forwarding.



## Key ConceptsMittens is a kubectl plugin for intercepting and inspecting traffic to Kubernetes Services using mitmproxy. It deploys an mitmproxy sidecar container next to your target Service, allowing you to see and modify traffic in real-time without port forwarding.



### What is Mittens?



Mittens (from German "Mitten" = "in the middle") is a lightweight kubectl plugin that makes it easy to tap into Kubernetes Service traffic for debugging, testing, and security analysis. Unlike traditional approaches that require manual port forwarding or external proxies, mittens:## Key ConceptsMittens is a kubectl plugin for intercepting and inspecting traffic to Kubernetes Services using mitmproxy. It deploys an mitmproxy sidecar container next to your target Service, allowing you to see and modify traffic in real-time without port forwarding.<p align="center">



- **Deploys automatically** - Sets up mitmproxy as a sidecar in seconds

- **Captures all traffic** - Intercepts traffic from any source to your Service

- **Interactive inspection** - Browse and modify requests/responses in mitmproxy's TUI### What is Mittens?  <img src='img/kubetap.png' class='smallimg' height='600'/>

- **Auto-cleanup** - Removes proxy setup when you exit

- **K9s integration** - One-click tapping from the k9s dashboard



## Use CasesMittens (from German "Mitten" = "in the middle") is a lightweight kubectl plugin that makes it easy to tap into Kubernetes Service traffic for debugging, testing, and security analysis. Unlike traditional approaches that require manual port forwarding or external proxies, mittens:## Key Concepts</p>



### Security Testing and Vulnerability Assessment



When assessing web applications, security testers need visibility into all traffic destined for a target Service—not just requests they generate locally.- **Deploys automatically** - Sets up mitmproxy as a sidecar in seconds



Traditional approach with BurpSuite or mitmproxy:- **Captures all traffic** - Intercepts traffic from any source to your Service

- Limited to traffic from the tester's machine

- Cannot see traffic from internal services or scheduled jobs- **Interactive inspection** - Browse and modify requests/responses in mitmproxy's TUI### What is Mittens?Kubetap is a [![open source](img/GitHub-Mark-32px.png)][kubetapGH]

- Requires complex port forwarding setup

- **Auto-cleanup** - Removes proxy setup when you exit

**Mittens approach:**

- Intercepts all inbound traffic to the Service, regardless of source- **K9s integration** - One-click tapping from the k9s dashboard[open source](https://github.com/Lappihuan/kubetap) CNI-agnostic

- See requests from other microservices, cron jobs, webhooks, etc.

- One command to start tapping: `kubectl mittens on <service>`

- Full request/response modification capabilities

## Use CasesMittens (from German "Mitten" = "in the middle") is a lightweight kubectl plugin that makes it easy to tap into Kubernetes Service traffic for debugging, testing, and security analysis. Unlike traditional approaches that require manual port forwarding or external proxies, mittens:project by [![Soluble](img/soluble-logo-very-small.png)][soluble]

### Developer Debugging



When debugging production or staging environments, developers often need visibility into incoming traffic patterns and data structures that are difficult to replicate locally.

### Security Testing and Vulnerability Assessmentthat automates the process of proxying Kubernetes Services.

Traditional approach:

- Deploy debugger locally and somehow connect to cluster

- Add logging statements and redeploy

- Use tools like Telepresence (heavyweight and not always practical)When assessing web applications, security testers need visibility into all traffic destined for a target Service—not just requests they generate locally.- **Deploys automatically** - Sets up mitmproxy as a sidecar in seconds



**Mittens approach:**

- Deploy mitmproxy sidecar in seconds

- Inspect live traffic in interactive TUITraditional approach with BurpSuite or mitmproxy:- **Captures all traffic** - Intercepts traffic from any source to your Service---

- See exact request/response format, headers, timing

- Modify responses for testing edge cases- Limited to traffic from the tester's machine

- No code changes or redeployment needed

- Cannot see traffic from internal services or scheduled jobs- **Interactive inspection** - Browse and modify requests/responses in mitmproxy's TUI

### Microservices Integration Testing

- Requires complex port forwarding setup

When testing microservices that communicate with each other, you need to:

- Monitor what downstream services are sending- **Auto-cleanup** - Removes proxy setup when you exit<div class="video">

- Test error conditions and edge cases

- Validate request formats and timing**Mittens approach:**



Mittens gives you real-time visibility into this cross-service traffic without complex setup.- Intercepts all inbound traffic to the Service, regardless of source- **K9s integration** - One-click tapping from the k9s dashboard  <iframe src="https://www.youtube.com/embed/hBroFtlxvkM" frameborder="0" allowfullscreen></iframe>



## Getting Started- See requests from other microservices, cron jobs, webhooks, etc.



- **[Quick Start Guide](getting_started/quick-start.md)** - Get up and running in 5 minutes- One command to start tapping: `kubectl mittens on <service>`</div>

- **[Installation](getting_started/installation.md)** - Install mittens (recommended: krew)

- **[Usage Guide](getting_started/usage.md)** - Complete command reference and workflows- Full request/response modification capabilities

- **[K9s Integration](getting_started/k9s-integration.md)** - Use mittens directly from k9s with Ctrl+M

## Use Cases

## How It Works

### Developer Debugging

1. **Enable**: Run `kubectl mittens on <service>` in your target Kubernetes cluster

2. **Deploy**: Mittens deploys an mitmproxy sidecar container next to your target pod[VIDEO][kubetapDemo]: Kubetap introduction by

3. **Intercept**: All traffic destined for the Service flows through mitmproxy

4. **Inspect**: Open an interactive mitmproxy TUI session using tmux + kubectl execWhen debugging production or staging environments, developers often need visibility into incoming traffic patterns and data structures that are difficult to replicate locally.

5. **Modify**: View, filter, and modify requests/responses in real-time

6. **Cleanup**: Press `q` to exit, mittens automatically removes the sidecar and cleans up### Security Testing and Vulnerability Assessment[Matt Hamilton][erinerGH]



## RequirementsTraditional approach:



- **Kubernetes**: 1.19 or later- Deploy debugger locally and somehow connect to cluster

- **kubectl**: 1.19 or later with plugin support

- **Container Runtime**: Docker, containerd, or similar- Add logging statements and redeploy

- **Network Access**: Can connect to mitmproxy listening port (default: 7777)

- Use tools like Telepresence (heavyweight and not always practical)When assessing web applications, security testers need visibility into all traffic destined for a target Service—not just requests they generate locally.## Use Cases

## Key Features



1. **Zero Configuration** - Works out of the box with sensible defaults

2. **Interactive Debugging** - Full mitmproxy TUI for inspecting traffic**Mittens approach:**

3. **Request Modification** - Change headers, bodies, response codes on the fly

4. **Automatic Cleanup** - All resources removed when you exit- Deploy mitmproxy sidecar in seconds

5. **K9s Integration** - One-key tapping from k9s dashboard (Ctrl+M)

6. **Multi-Protocol** - HTTP, HTTPS, WebSockets, gRPC (mitmproxy features)- Inspect live traffic in interactive TUITraditional approach with BurpSuite or mitmproxy:### Security Testing

7. **Persistent Sessions** - Stay connected even if your terminal disconnects (tmux)

8. **Pod and Service Support** - Tap services or specific pods- See exact request/response format, headers, timing

9. **Custom Images** - Use your own mitmproxy image or configuration

10. **No Code Changes** - Inspect services without redeployment- Modify responses for testing edge cases- Limited to traffic from the tester's machine



## Architecture- No code changes or redeployment needed



Unlike traditional approaches that use port forwarding:- Cannot see traffic from internal services or scheduled jobsWhen assessing web applications, it is common to use BurpSuite, MITMproxy, Zap,



```### Microservices Integration Testing

Traditional Approach:

Service Pod           Tunnel              Local Machine- Requires complex port forwarding setupor other intercepting proxy to capture and modify HTTP requests on a

┌──────────────┐      ========            ┌──────────────┐

│ App          │◄──────port-fwd─────────►│ Proxy Tool   │When testing microservices that communicate with each other, you need to:

│ Container    │                         │ (BurpSuite)  │

└──────────────┘                         └──────────────┘- Monitor what downstream services are sendingsecurity tester’s machine. These requests are intercepted and modified on the



Mittens Approach (All in Cluster):- Test error conditions and edge cases

┌──────────────────────────────────────┐

│  Kubernetes Cluster                  │- Validate request formats and timing**Mittens approach:**tester’s local machine prior to being sent to the remote server.

│  ┌──────────────────────────────────┐│

│  │  Service Pod                     ││

│  │  ┌────────────────────────────┐ ││

│  │  │ App Container              │ ││Mittens gives you real-time visibility into this cross-service traffic without complex setup.- Intercepts all inbound traffic to the Service, regardless of source

│  │  └────────────────────────────┘ ││

│  │  ┌────────────────────────────┐ ││

│  │  │ mitmproxy Sidecar          │ ││

│  │  └─────────────┬──────────────┘ ││## Getting Started- See requests from other microservices, cron jobs, webhooks, etc.<img src='img/traditional-webapp-testing.png' class='img'/>

│  └───────────────┼────────────────┘│

│                  │ tmux session    │

│    ┌─────────────▼──────────────┐  │

│    │ Your Terminal              │  │- **[Quick Start Guide](getting_started/quick-start.md)** - Get up and running in 5 minutes- One command to start tapping: `kubectl mittens on <service>`

│    │ (mitmproxy TUI)            │  │

│    └────────────────────────────┘  │- **[Installation](getting_started/installation.md)** - Install mittens (recommended: krew)

└──────────────────────────────────────┘

```- **[Usage Guide](getting_started/usage.md)** - Complete command reference and workflows- Full request/response modification capabilitiesWhile this paradigm allowed testers to capture and modify all traffic that the



**Benefits:**- **[K9s Integration](getting_started/k9s-integration.md)** - Use mittens directly from k9s with Ctrl+M

- All traffic is transparently captured (no client changes needed)

- Interactive session stays local to the pod (no network tunneling overhead)testers themselves create, testers can not see traffic destined for the target

- Automatic cleanup when session ends

- Works across cluster boundaries (pod can be on any node)## How It Works



## Quick Example### Developer Debuggingserver that originates from other Services. The lack of visibility of intranet



```sh1. **Enable**: Run `kubectl mittens on <service>` in your target Kubernetes cluster

# 1. Install mittens

kubectl krew install mittens2. **Deploy**: Mittens deploys an mitmproxy sidecar container next to your target podtraffic that does not originate from a tester’s machine can hamper a tester’s



# 2. List available services3. **Intercept**: All traffic destined for the Service flows through mitmproxy

kubectl mittens list

4. **Inspect**: Open an interactive mitmproxy TUI session using tmux + kubectl execWhen debugging production or staging environments, developers often need visibility into incoming traffic patterns and data structures that are difficult to replicate locally.ability to competently review complex systems and environments.

# 3. Start tapping a service

kubectl mittens on my-service -n my-namespace5. **Modify**: View, filter, and modify requests/responses in real-time



# 4. Interactive mitmproxy TUI opens6. **Cleanup**: Press `q` to exit, mittens automatically removes the sidecar and cleans up

# - Use arrow keys to navigate requests

# - Press 'e' to edit a request

# - Press 'd' to change request details

# - Press 'space' to inspect response details## RequirementsTraditional approach:<img src='img/complex-k8s-webapp-testing.png' class='img'/>

# - Type 'q' to exit (and auto-cleanup)



# 5. Already exited? Manually cleanup with:

kubectl mittens off my-service -n my-namespace- **Kubernetes**: 1.19 or later- Deploy debugger locally and somehow connect to cluster

```

- **kubectl**: 1.19 or later with plugin support

## See Also

- **Container Runtime**: Docker, containerd, or similar- Add logging statements and redeployFor environments that use Kubernetes, Kubetap is altering the status quo.

- **[Full Usage Guide](getting_started/usage.md)** - All commands and options

- **[K9s Integration](getting_started/k9s-integration.md)** - Tap from k9s with Ctrl+M- **Network Access**: Can connect to mitmproxy listening port (default: 7777)

- **[Contributing](mittens_development/contributing.md)** - Help improve mittens

- Use tools like Telepresence (heavyweight and not always practical)

## License

## Key Features

Mittens is licensed under the Apache License 2.0. For original attribution, see [ATTRIBUTION.md](../ATTRIBUTION.md).

Kubetap allows testers to select a target Service and intercept all traffic

1. **Zero Configuration** - Works out of the box with sensible defaults

2. **Interactive Debugging** - Full mitmproxy TUI for inspecting traffic**Mittens approach:**that is destined for that Service, regardless of where the requests originate.

3. **Request Modification** - Change headers, bodies, response codes on the fly

4. **Automatic Cleanup** - All resources removed when you exit- Deploy mitmproxy sidecar in seconds

5. **K9s Integration** - One-key tapping from k9s dashboard (Ctrl+M)

6. **Multi-Protocol** - HTTP, HTTPS, WebSockets, gRPC (mitmproxy features)- Inspect live traffic in interactive TUIThe transparency and visibility afforded by Kubetap allows testers to better

7. **Persistent Sessions** - Stay connected even if your terminal disconnects (tmux)

8. **Pod and Service Support** - Tap services or specific pods- See exact request/response format, headers, timingunderstand and exercise the Service without the prohibitively (expensive) time

9. **Custom Images** - Use your own mitmproxy image or configuration

10. **No Code Changes** - Inspect services without redeployment- Modify responses for testing edge casescost of configuring and deploying a proxy manually. **Microservices deep in a



## Architecture- No code changes or redeployment neededtechnology stack that were once inaccessible to testers can now be proxied with ease.**



Unlike traditional approaches that use port forwarding:



```### Microservices Integration Testing<img src='img/kubetap-proxying.png' class='img'/>

Traditional Approach:

Service Pod           Tunnel              Local Machine

┌──────────────┐      ========            ┌──────────────┐

│ App          │◄──────port-fwd─────────►│ Proxy Tool   │When testing microservices that communicate with each other, you need to:### Developer debugging

│ Container    │                         │ (BurpSuite)  │

└──────────────┘                         └──────────────┘- Monitor what downstream services are sending



Mittens Approach (All in Cluster):- Test error conditions and edge casesWhen an application or microservice is exhibiting unintended behavior,

┌──────────────────────────────────────┐

│  Kubernetes Cluster                  │- Validate request formats and timingdevelopers must debug the application through a debugger, printf statements,

│  ┌──────────────────────────────────┐│

│  │  Service Pod                     ││or static code analysis. This is often because infrastructure architecture looks

│  │  ┌────────────────────────────┐ ││

│  │  │ App Container              │ ││Mittens gives you real-time visibility into this cross-service traffic without complex setup.something like this:

│  │  └────────────────────────────┘ ││

│  │  ┌────────────────────────────┐ ││

│  │  │ mitmproxy Sidecar          │ ││

│  │  └─────────────┬──────────────┘ ││## Getting Started<img src='img/typical-architecture.png' class='img'/>

│  └───────────────┼────────────────┘│

│                  │ tmux session    │

│    ┌─────────────▼──────────────┐  │

│    │ Your Terminal              │  │- **[Quick Start Guide](getting_started/quick-start.md)** - Get up and running in 5 minutesWhat happens when the bug is only exhibited when deployed to staging and

│    │ (mitmproxy TUI)            │  │

│    └────────────────────────────┘  │- **[Installation](getting_started/installation.md)** - Install mittens (recommended: krew)production, and not in a local development environment?

└──────────────────────────────────────┘

```- **[Usage Guide](getting_started/usage.md)** - Complete command reference and workflows



**Benefits:**- **[K9s Integration](getting_started/k9s-integration.md)** - Use mittens directly from k9s with Ctrl+M![xkcd-979][xkcd]

- All traffic is transparently captured (no client changes needed)

- Interactive session stays local to the pod (no network tunneling overhead)[xkcd-979][xkcd]

- Automatic cleanup when session ends

- Works across cluster boundaries (pod can be on any node)## How It Works



## Quick ExampleWhile there are tools like [Telepresence][telepresence] that allow developers to



```sh1. **Enable**: Run `kubectl mittens on <service>` in your target Kubernetes clustermove containers running in a cluster on to their local machines for debugging,

# 1. Install mittens

kubectl krew install mittens2. **Deploy**: Mittens deploys an mitmproxy sidecar container next to your target podthis is a heavy-handed approach and not practical in many situations. Often



# 2. List available services3. **Intercept**: All traffic destined for the Service flows through mitmproxydevelopers just need to visually inspect Service inputs, such as JSON objects,

kubectl mittens list

4. **Inspect**: Open an interactive mitmproxy TUI session using tmux + kubectl execthat originate from other microservices.

# 3. Start tapping a service

kubectl mittens on my-service -n my-namespace5. **Modify**: View, filter, and modify requests/responses in real-time



# 4. Interactive mitmproxy TUI opens6. **Cleanup**: Press `q` to exit, mittens automatically removes the sidecar and cleans up<img src='img/proverb.png' class='img'/>

# - Use arrow keys to navigate requests

# - Press 'e' to edit a request

# - Press 'd' to change request details

# - Press 'space' to inspect response details## RequirementsKubetap allows developers to deploy MITMproxy in front of a Service, enabling

# - Type 'q' to exit (and auto-cleanup)

visibility for all incoming HTTP traffic to that Service. With this, developers

# 5. Already exited? Manually cleanup with:

kubectl mittens off my-service -n my-namespace- **Kubernetes**: 1.19 or latercan inspect and debug a Service without unnecessary printf debugging code-pushes

```

- **kubectl**: 1.19 or later with plugin support

## See Also

- **Container Runtime**: Docker, containerd, or similar## Getting Started

- **[Full Usage Guide](getting_started/usage.md)** - All commands and options

- **[K9s Integration](getting_started/k9s-integration.md)** - Tap from k9s with Ctrl+M- **Network Access**: Can connect to mitmproxy listening port (default: 7777)

- **[Contributing](../kubetap_development/contributing.md)** - Help improve mittens

- [Quick Start Guide](getting_started/quick-start.md) - Get up and running in minutes

## License

## Key Features- [Installation](getting_started/installation.md) - Install Kubetap

Mittens is licensed under the Apache License 2.0. For original attribution, see [ATTRIBUTION.md](../../ATTRIBUTION.md).

- [Usage Guide](getting_started/usage.md) - Learn how to use Kubetap

1. **Zero Configuration** - Works out of the box with sensible defaults- [K9s Integration](getting_started/k9s-integration.md) - Use Kubetap directly from k9s

2. **Interactive Debugging** - Full mitmproxy TUI for inspecting traffic

3. **Request Modification** - Change headers, bodies, response codes on the fly[kubetapGH]: https://github.com/Lappihuan/kubetap

4. **Automatic Cleanup** - All resources removed when you exit[soluble]: https://www.soluble.ai

5. **K9s Integration** - One-key tapping from k9s dashboard (Ctrl+M)[kubetapDemo]: https://www.youtube.com/watch?v=hBroFtlxvkM

6. **Multi-Protocol** - HTTP, HTTPS, WebSockets, gRPC (mitmproxy features)[erinerGH]: https://github.com/eriner

7. **Persistent Sessions** - Stay connected even if your terminal disconnects (tmux)[telepresence]: https://github.com/telepresenceio/telepresence

8. **Pod and Service Support** - Tap services or specific pods[xkcd]: https://imgs.xkcd.com/comics/wisdom_of_the_ancients.png

9. **Custom Images** - Use your own mitmproxy image or configuration
10. **No Code Changes** - Inspect services without redeployment

## Architecture

Unlike traditional approaches that use port forwarding:

```
Traditional Approach:
Service Pod           Tunnel              Local Machine
┌──────────────┐      ========            ┌──────────────┐
│ App          │◄──────port-fwd─────────►│ Proxy Tool   │
│ Container    │                         │ (BurpSuite)  │
└──────────────┘                         └──────────────┘

Mittens Approach (All in Cluster):
┌──────────────────────────────────────┐
│  Kubernetes Cluster                  │
│  ┌──────────────────────────────────┐│
│  │  Service Pod                     ││
│  │  ┌────────────────────────────┐ ││
│  │  │ App Container              │ ││
│  │  └────────────────────────────┘ ││
│  │  ┌────────────────────────────┐ ││
│  │  │ mitmproxy Sidecar          │ ││
│  │  └─────────────┬──────────────┘ ││
│  └───────────────┼────────────────┘│
│                  │ tmux session    │
│    ┌─────────────▼──────────────┐  │
│    │ Your Terminal              │  │
│    │ (mitmproxy TUI)            │  │
│    └────────────────────────────┘  │
└──────────────────────────────────────┘
```

**Benefits:**
- All traffic is transparently captured (no client changes needed)
- Interactive session stays local to the pod (no network tunneling overhead)
- Automatic cleanup when session ends
- Works across cluster boundaries (pod can be on any node)

## Quick Example

```sh
# 1. Install mittens
kubectl krew install mittens

# 2. List available services
kubectl mittens list

# 3. Start tapping a service
kubectl mittens on my-service -n my-namespace

# 4. Interactive mitmproxy TUI opens
# - Use arrow keys to navigate requests
# - Press 'e' to edit a request
# - Press 'd' to change request details
# - Press 'space' to inspect response details
# - Type 'q' to exit (and auto-cleanup)

# 5. Already exited? Manually cleanup with:
kubectl mittens off my-service -n my-namespace
```

## See Also

- **[Full Usage Guide](getting_started/usage.md)** - All commands and options
- **[K9s Integration](getting_started/k9s-integration.md)** - Tap from k9s with Ctrl+M
- **[Contributing](../kubetap_development/contributing.md)** - Help improve mittens

## License

Mittens is licensed under the Apache License 2.0. For original attribution, see [ATTRIBUTION.md](../../ATTRIBUTION.md).
