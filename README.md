# vault-hello-world

This is an example of how to write a secret plugin for [HashiCorp Vault](https://vault.io). I wasn't able to find any other example of a secret plugin online so hopefully this might be of use to others.

## Building

The plugin can be built as any other Golang project by using `go build` to produce the plugin binary.

## Deployment

Make sure that you have set the `plugin_directory` key in your Vault configuration. Move the binary to that directory.

Calculate the checksum of the binary, and add it to the plugin registry:

```
SHASUM=$(shasum -a 256 "/tmp/vault-plugins/vault-hello-world" | cut -d " " -f1)
vault write sys/plugins/catalog/example-plugin \
  sha_256="$SHASUM" \
  command="vault-hello-world"
```

Mount the plugin to your Vault:

```
vault mount -path=example -plugin-name=example-plugin plugin
```

Now test it out :) !

```
vault write example/hello name=John
```