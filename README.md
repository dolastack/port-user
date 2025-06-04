# port-user

🔍 A cross-platform CLI utility to determine which process is using one or more TCP/UDP ports.

Supports:
- ✅ Linux (amd64/arm64)
- ✅ macOS (Intel & Apple Silicon)
- ✅ Windows (amd64)

Features:
- 🔍 Find processes by port(s)
- 🔄 Supports multiple ports: `80,443,22`
- 📦 Built-in `.deb`, `.rpm`, `.pkg` packaging
- 🎨 Colorized output (`--color`)
- 🧾 JSON format support (`--json`)
- 🐍 Written in Go + [Cobra CLI](https://github.com/spf13/cobra) 

---

## 🚀 Quick Start

### Install via Go

```bash
go install github.com/dolastack/port-user@latest
```

### Build from Source
```
git clone https://github.com/dolastack/port-user.git 
cd port-user
go run main.go port-user --help
```

Or build binaries for all platforms:
```
./build.sh
```
## Usage
```
# Basic usage
port-user 8080

# Multiple ports
port-user 80,443,22

# JSON output
port-user --json 8080

# Colorized output
port-user -c 8080

# Help
port-user --help
```
Output Example
```
Proto     PID     Executable             Port
------------------------------------------------------------
TCP       1234    /usr/bin/python3       8080
UDP       5678    /usr/sbin/dhcpd        67
```
JSON Output
```
[
  {
    "pid": "1234",
    "executable": "/usr/bin/python3",
    "protocol": "TCP",
    "port": 8080
  }
]
```

#### Shell Autocompletion
To enable autocompletion:

Bash
```
source <(port-user completion bash)
```

Zsh
```
source <(port-user completion bash)
```
Fish
```
port-user completion fish | source
```
Powershell
```
port-user completion powershell | Out-String | Invoke-Expression
```

To persist across sessions:

Bash (Linux/macOS)
```
echo 'source <(port-user completion bash)' >> ~/.bashrc
source ~/.bashrc
```

Zsh
```
echo 'source <(port-user completion zsh)' >> ~/.zshrc
source ~/.zshrc
```
Fish
```
port-user completion fish >> ~/.config/fish/completions/port-user.fish
```
PowerShell
Add to profile:
```
Add-Content -Path $PROFILE.CurrentUserAllHosts -Value "`nport-user completion powershell | Out-String | Invoke-Expression"
```