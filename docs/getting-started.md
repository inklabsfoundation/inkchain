## Development Environment

This folder contains the files which are used for bootstrapping the INK Consortium chain development environment.
And so far, we recommend using ubuntu 16.04.

We provide an [env setup script](../devops/setup-env.sh) to config your environment.

### Clone code

Clone the Inkchain code base into this path 

```bash
$ mkdir -p $GOPATH/src/github.com/inklabsfoundation

$ cd $GOPATH/src/github.com/inklabsfoundation

$ git clone https://github.com/inklabsfoundation/inkchain.git

```

### Build the binaries and images

Now make the platform-specific binaries for `cryptogen` , `configtxgen` , `configtxlator` ,
`peer`, `orderer`.

```bash
$ cd $GOPATH/src/github.com/inklabsfoundation/inkchain

$ make release
```

This will output platform-specific binaries into the ``inkchain/release`` folder.

Next, make the `Inkchain` images. This will takes between ten to twenty minutes, so be patient.

```bash
# make sure you are in the /inkchain directory

$ make docker

```

> recommendation: You can also download those images at `https://hub.docker.com/u/inkchain/`

Execute a `docker images` command in yout terminal, If the images compiled successfully, you should
see output similar to the following:

```bash
REPOSITORY                             TAG                    IMAGE ID            CREATED             SIZE
inklabsfoundation/inkchain-tools       latest                 b94ae27d0032        8 days ago          1.342 GB
inklabsfoundation/inkchain-tools       x86_64-0.10.3          b94ae27d0032        8 days ago          1.342 GB
inklabsfoundation/inkchain-couchdb     latest                 71e7a843daac        8 days ago          1.514 GB
inklabsfoundation/inkchain-couchdb     x86_64-0.10.3          71e7a843daac        8 days ago          1.514 GB
inklabsfoundation/inkchain-kafka       latest                 1448be36a1d7        9 days ago          1.309 GB
inklabsfoundation/inkchain-kafka       x86_64-0.10.3          1448be36a1d7        9 days ago          1.309 GB
inklabsfoundation/inkchain-zookeeper   latest                 e9e3cc6b2883        9 days ago          1.327 GB
inklabsfoundation/inkchain-zookeeper   x86_64-0.10.3          e9e3cc6b2883        9 days ago          1.327 GB
inklabsfoundation/inkchain-orderer     latest                 fbb3b9b060a6        9 days ago          180.1 MB
inklabsfoundation/inkchain-orderer     x86_64-0.10.3          fbb3b9b060a6        9 days ago          180.1 MB
inklabsfoundation/inkchain-peer        latest                 ae96a865e213        9 days ago          183.4 MB
inklabsfoundation/inkchain-peer        x86_64-0.10.3          ae96a865e213        9 days ago          183.4 MB
inklabsfoundation/inkchain-javaenv     latest                 ef62109f2ad1        9 days ago          1.425 GB
inklabsfoundation/inkchain-javaenv     x86_64-0.10.3          ef62109f2ad1        9 days ago          1.425 GB
inklabsfoundation/inkchain-ccenv       latest                 dcaef40a4d5a        9 days ago          1.293 GB
inklabsfoundation/inkchain-ccenv       x86_64-0.10.3          dcaef40a4d5a        9 days ago          1.293 GB
inklabsfoundation/inkchain-ca          latest                 4f610a5b8393        7 weeks ago         238.2 MB
inklabsfoundation/inkchain-ca          x86_64-0.10.3          4f610a5b8393        7 weeks ago         238.2 MB
inklabsfoundation/inkchain-baseimage   x86_64-0.3.1           9f2e9ec7c527        7 months ago        1.268 GB
inklabsfoundation/inkchain-baseos      x86_64-0.3.1           4b0cab202084        7 months ago        156.6 MB
```

## Run the examples

We provide an example [token issue](../examples/issue-token) which based on inkchain.
