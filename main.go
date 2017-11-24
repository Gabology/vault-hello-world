// Based on: https://www.hashicorp.com/blog/building-a-vault-secure-plugin
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/vault/logical/framework"

	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/plugin"
)

func main() {
	apiClientMeta := &pluginutil.APIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args)

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := pluginutil.VaultPluginTLSProvider(tlsConfig)

	// Builds all the required plugin APIs, TLS connections, and RPC server.
	// The BackendFactoryFunc calls the local Factory function.
	if err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: Factory,
		TLSProviderFunc:    tlsProviderFunc,
	}); err != nil {
		log.Fatal(err)
	}
}

// Factory ...
// The factory is responsible for setting up and configuring the plugin (sometimes called a "backend" internally),
// returning any errors that occur during setup.
func Factory(c *logical.BackendConfig) (logical.Backend, error) {
	b := Backend(c)

	if err := b.Setup(c); err != nil {
		return nil, err
	}

	return b, nil
}

// The backend struct embeds the standard framework.Backend.
// This allows our backend to inherit almost all the required functions and properties without writing more boilerplate code.
type backend struct {
	*framework.Backend
}

func (b *backend) pathSayHello(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return &logical.Response{
		Data: map[string]interface{}{
			"message": fmt.Sprintf("Hello %s!", d.Get("name")),
		},
	}, nil
}

// Backend is an implementation of logical.Backend
// that allows the implementer to code a backend using a much more programmer-friendly framework
// that handles a lot of the routing and validation for you.
// This is recommended over implementing logical.Backend directly.
func Backend(c *logical.BackendConfig) *backend {
	var b backend

	b.Backend = &framework.Backend{
		BackendType: logical.TypeLogical,
		// Tell Vault that this is a special path that does not require any Vault token
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{"hello"},
		},
		Paths: []*framework.Path{
			&framework.Path{
				Pattern: "hello",
				// Define a user input field of type string
				Fields: map[string]*framework.FieldSchema{
					"name": &framework.FieldSchema{
						Type: framework.TypeString,
					},
				},
				Callbacks: map[logical.Operation]framework.OperationFunc{
					// This is the callback that will be invoked for the HTTP PUT & POST verb.
					logical.UpdateOperation: b.pathSayHello,
				},
			},
		},
	}

	return &b
}
