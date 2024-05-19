# Duplik8s

🚧 Work in progress!

---

**Duplicate** 🔁 kubectl plugin to duplicate resources in a Kubernetes cluster.

<p>
    <a href="https://github.com/Telemaco019/duplik8s/releases"><img src="https://img.shields.io/github/release/Telemaco019/duplik8s.svg" alt="Latest Release"></a>
    <a href="https://github.com/Telemaco019/duplik8s/actions"><img src="https://github.com/Telemaco019/duplik8s/actions/workflows/ci.yaml/badge.svg" alt="Build Status"></a>
</p>

---

![](./docs/demo.gif)

`duplik8s` allows you to easily duplicate Kubernetes pods with overridden commands and configurations.
This is useful for testing, debugging, and development purposes.

As you might have guessed, `duplik8s` shines when used in combination with the
amazing [k9s](https://github.com/derailed/k9s) ✨.
Check out the installation instructions below to easily load it as a k9s plugin.

## Installation

### Install with Go

```sh
$ go install github.com/telemaco019/duplik8s/kubectl-duplicate@latest
```

### Use as k9s plugin

After installing `duplik8s`, you can add it to your k9s plugins by adding the following to
your `$XDG_CONFIG_HOME/k9s/plugins.yml` file.

After reloading k9s, you should be able to duplicate Pods with `Ctrl-T`.

```yaml
# $XDG_CONFIG_HOME/k9s/plugins.yaml
plugins:
  duplik8s:
    shortCut: Ctrl-T
    description: Duplicate Pod
    scopes:
      - po
    command: kubectl
    background: true
    args:
      - duplicate
      - pod
      - $NAME
      - -n
      - $NAMESPACE
      - --context
      - $CONTEXT
```

On MacOS, you can find the `plugins.yml` file at `~/Library/Application Support/k9s/plugins.yaml`.

For more information on k9s plugins, you can refer to the [official documentation](https://k9scli.io/topics/plugins).

## Examples

Duplicate a Pod:

```sh
$ kubectl duplicate pod my-pod
```


--- 

## License 

This project is licensed under the Apache License. See the [LICENSE](./LICENSE) file for details.

