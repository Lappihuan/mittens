# Attribution

## Original Project: Kubetap

Mittens is a fork of the [kubetap](https://github.com/soluble-ai/kubetap) project, originally created by [Soluble](https://www.soluble.ai/).

### Original License

Kubetap is licensed under the [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0), and **Mittens continues to be licensed under the same Apache License 2.0**.

### Original Authors & Contributors

- Created by Soluble Inc.
- Original repository: https://github.com/soluble-ai/kubetap
- License header: Copyright 2020 Soluble Inc

### Major Changes in Mittens

Mittens is a heavily modified fork that changes the core architecture and approach:

1. **From mitmweb to mitmproxy**: Replaced mitmweb browser UI with mitmproxy TUI (Terminal User Interface)
2. **Direct TUI Attachment**: Users attach directly to an interactive tmux session with mitmproxy instead of using port forwarding and opening a browser
3. **Simplified Architecture**: Removed port forwarding complexity entirely in favor of direct `kubectl exec` integration
4. **Automated Cleanup**: Automatic cleanup of taps when user exits the tmux session
5. **K9s Integration**: Added "Mittens" k9s plugin for quick tapping from the k9s dashboard


### Respect for Original Work

This project maintains:
- Full Apache 2.0 license compliance
- Original copyright headers in source files
- Clear attribution to Soluble and the original kubetap project
- All contributions continue under Apache 2.0

### Compatibility Note

Mittens is NOT backward compatible with kubetap as it fundamentally changes how the tool works. It is a spiritual successor and fork with a different philosophy on how to provide Kubernetes service proxying.
