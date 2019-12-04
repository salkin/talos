<!-- markdownlint-disable -->
## osctl gen ca

Generates a self-signed X.509 certificate authority

### Synopsis

Generates a self-signed X.509 certificate authority

```
osctl gen ca [flags]
```

### Options

```
  -h, --help                  help for ca
      --hours int             the hours from now on which the certificate validity period ends (default 87600)
      --organization string   X.509 distinguished name for the Organization
      --rsa                   generate in RSA format
```

### Options inherited from parent commands

```
      --talosconfig string   The path to the Talos configuration file (default "/root/.talos/config")
  -t, --target strings       target the specificed node
```

### SEE ALSO

* [osctl gen](osctl_gen.md)	 - Generate CAs, certificates, and private keys

###### Auto generated by spf13/cobra on 13-Nov-2019