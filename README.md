# go-docs

Go package documentation for internal modules. Self hosted server that automatically generates go doc pages like the ones in https://pkg.go.dev/, but for private repositories.

For more information on go doc comments, check out https://go.dev/doc/comment.

System diagram and other useful information can be found [here](./docs/system.md).

## Add go doc to your internal go Repo

If your repository is a valid go module within `github.com/private/` and you have comments in go doc format (https://go.dev/doc/comment), `go-docs` should work out-of-the-box for you after the following steps:

1. Update any insances of `github.com/private` -> `github.com/name-of-your-private-org`
2. Update or add auth as directed in the comments throughout this project; this may require different or extra steps depending on your organization's setup. A good way to confirm this would be [running locally](#running-locally)
3. Deploy go-docs server behind an endpoint accessible only within your VPN range. You can now access your go documentation:

For example, the repo https://github.com/private/go-module defines the module `github.com/private/go-module`:
 - go docs automatically generated at https://endpoint-that-you-set-up-for-go-docs-server.com/**github.com/private/go-module**

## Running Locally

 - Copy your `.netrc` to the repository root
 - Copy your `id_rsa` to the `ssh-keys` directory
 - Run `docker-compose up`
 - Visit doc site at `127.0.0.1:8080/<go module>`
   - For example, `127.0.0.1:8080/github.com/private/go-module`
   - You can find the go module in the `go.mod` file of the go project: `module <go module>`

### Viewing Documentation for your repository not on @latest/feature branches

If you want to view go docs for a specific version, simply add that to the module path. For example, `https://endpoint-that-you-set-up-for-go-docs-server.com/github.com/private/go-module@v1.2.3-<checksum>`

This can be used to view how the documentation will look for feature branches before merging.

To find the version of a feature branch `dev-branch` for the module `go-module`, you can find the version with:
```
> go list -m -json github.com/private/go-module@dev-branch
{
	"Path": "github.com/private/go-module",
	"Version": "v1.2.3-<checksum>",
	"Query": "dev-branch",
	"Time": "2024-10-11T19:09:40Z",
	"Origin": {
		"Ref": "refs/heads/dev-branch",
		...
	}
}
```
