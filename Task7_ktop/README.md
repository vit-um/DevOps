## Розширення kubectl (Extending kubectl)
1. Для виконання задачі розберемось як працює [приклад](https://github.com/vladimirvivien/ktop)
- Оберемо метод встановлення за допомогою вже існуючого в системі go:
```sh
$ go install github.com/vladimirvivien/ktop@latest
$ export PATH="/root/go/bin:$PATH"
$ ktop             
Connected to: https://0.0.0.0:43297
 _    _ 
| | _| |_ ___  _ __
| |/ / __/ _ \| '_ \
|   <| || (_) | |_) |
|_|\_\\__\___/| .__/
              |_|
```

![ktop](image.png)

2. Тепер розберемось з механізмом встановлення [плагінів для `kubectl](https://krew.sigs.k8s.io/plugins/)  
```sh
(
  set -x; cd "$(mktemp -d)" &&
  OS="$(uname | tr '[:upper:]' '[:lower:]')" &&
  ARCH="$(uname -m | sed -e 's/x86_64/amd64/' -e 's/\(arm\)\(64\)\?.*/\1\2/' -e 's/aarch64$/arm64/')" &&
  KREW="krew-${OS}_${ARCH}" &&
  curl -fsSLO "https://github.com/kubernetes-sigs/krew/releases/latest/download/${KREW}.tar.gz" &&
  tar zxvf "${KREW}.tar.gz" &&
  ./"${KREW}" install krew
)

$ nano ~/.zshrc

# Add this row: export PATH="${KREW_ROOT:-$HOME/.krew}/bin:$PATH"

$ source ~/.zshrc

$ k krew      
krew is the kubectl plugin manager.
You can invoke krew through kubectl: "kubectl krew [command]..."

Usage:
  kubectl krew [command]

Available Commands:
  help        Help about any command
  index       Manage custom plugin indexes
  info        Show information about an available plugin
  install     Install kubectl plugins
  list        List installed kubectl plugins
  search      Discover kubectl plugins
  uninstall   Uninstall plugins
  update      Update the local copy of the plugin index
  upgrade     Upgrade installed plugins to newer versions
  version     Show krew version and diagnostics

Flags:
  -h, --help      help for krew
  -v, --v Level   number for the log level verbosity

Use "kubectl krew [command] --help" for more information about a command.

# Спробуємо встановити планін ns для перемикання між namespaces
$ k krew install ns
Updated the local copy of plugin index.
Installing plugin: ns
Installed plugin: ns
\
 | Use this plugin:
 |      kubectl ns
 | Documentation:
 |      https://github.com/ahmetb/kubectx
 | Caveats:
 | \
 |  | If fzf is installed on your machine, you can interactively choose
 |  | between the entries using the arrow keys, or by fuzzy searching
 |  | as you type.
 | /
/
WARNING: You installed plugin "ns" from the krew-index plugin repository.
   These plugins are not audited for security by the Krew maintainers.
   Run them at your own risk.

$ k ns             
kube-system
default
kube-public
kube-node-lease
argocd
demo

(⎈|k3d-argo:default)➜  ~ k ns demo
Context "k3d-argo" modified.
Active namespace is "demo".
(⎈|k3d-argo:demo)➜  ~ 
```

