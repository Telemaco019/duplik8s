# Changelog

## v0.3.0

### New features

* Add command `list` for showing duplicated resources. Example:

```shell
kubectl duplicate list -n  my-namespace
```

* Add command `cleanup` for deleting duplicated resources. Example:

```shell
kubectl duplicate cleanup -n  my-namespace
```

### Chores

* Switch
  to [Dynamic Client](https://github.com/kubernetes/client-go/blob/master/examples/dynamic-create-update-delete-deployment/README.md),
  opening the door for duplicating any resource type.

## v0.2.1

### Fixes

* In interactive selection, don't list resources that have already been duplicated.

## v0.2.0

### New features

* Properly handle Pods that mount persistent storage.
  Pods that mount a PersistentVolume with exclusive access modes (`ReadWriteOnce`, `ReadWriteOncePod`) are cloned
  on the same node as the original. This ensures that the duplicate can also mount the same volume.
* Add support for duplicating Deployments and StatefulSets.
* Interactively select Pods, Deployments, or StatefulSets to duplicate when no name is provided as an argument.

### Chores

* Refactoring to make code testable.
* Update demo GIFs in the README.

## v0.1.0

Initial release.