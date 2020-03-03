# :lock_with_ink_pen: Entities

Modern go identity manager for UNIX systems.

Entities parses includes file to generate UNIX-compliant `/etc/passwd` , `/etc/shadow` and `/etc/groups` files.
It can be used to handle identities management and honors already existing entities in the system.


```

$> entities apply -f <policy.yaml>
$> entities delete -f <policy.yaml>

```