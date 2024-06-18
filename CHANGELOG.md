# Changelog

## Unreleased

### New features

* Show prompt for selecting a Pod in the current namespace when no Pod name is provided as argument.
* Properly handle Pods that mount persistent storage.
  Pods that mount a PersistentVolume with exclusive access modes (`ReadWriteOnce`, `ReadWriteOncePod`) are cloned
  on the same node as the original. This ensures that the duplicate can also mount the same volume.

### Chores

* Refactoring to make code testable.

## v0.1.0

Initial release.