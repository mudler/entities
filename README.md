# :lock_with_ink_pen: Entities

Modern go identity manager for UNIX systems.

Entities parses includes file to generate UNIX-compliant `/etc/passwd` , `/etc/shadow` and `/etc/groups` files.
It can be used to handle identities management and honors already existing entities in the system.


```

$> entities apply <entity.yaml>
$> entities delete <entity.yaml>
$> entities create <entity.yaml>

```

## Entities file format

### Passwd

```yaml
kind: "user"
username: "foo"
password: "pass"
uid: 0
gid: 0
info: "Foo!"
homedir: "/home/foo"
shell: "/bin/bash"
```

### Gshadow

```yaml
kind: "gshadow"
name: "postmaster"
password: "foo"
administrators: "barred"
members: "baz"
```

### Shadow

```yaml
kind: "shadow"
username: "foo"
password: "bar"
last_changed: 1
minimum_changed: 2
maximum_changed: 3
warn: 4
inactive: 5
expire: 6
```

To define `last_changed` with a value equal to current days from 1970 use `now`.

### Group

```yaml
kind: "group"
group_name: "sddm"
password: "xx"
gid: 1
users: "one,two,tree"
```
