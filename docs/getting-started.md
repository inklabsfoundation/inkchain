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

> recommendation: Or you can use the [script](../scripts/inkchain_images_download.sh) to download the necessary images.

Execute a `docker images` command in yout terminal, If the images compiled successfully, you should
see output similar to the following:

```bash
REPOSITORY                             TAG                    IMAGE ID            CREATED             SIZE
inklabsfoundation/inkchain-tools         latest              1626ee3e191e        4 days ago          1.45GB
inklabsfoundation/inkchain-tools         x86_64-0.13.0       1626ee3e191e        4 days ago          1.45GB
inklabsfoundation/inkchain-peer          latest              a024514b8fd3        4 days ago          279MB
inklabsfoundation/inkchain-peer          x86_64-0.13.0       a024514b8fd3        4 days ago          279MB
inklabsfoundation/inkchain-ccenv         latest              6c53a39ef9fb        4 days ago          1.29GB
inklabsfoundation/inkchain-ccenv         x86_64-0.13.0       6c53a39ef9fb        4 days ago          1.29GB
inklabsfoundation/inkchain-orderer       latest              3cadc0059ccc        4 days ago          275MB
inklabsfoundation/inkchain-orderer       x86_64-0.13.0       3cadc0059ccc        4 days ago          275MB
inklabsfoundation/inkchain-javaenv       latest              cc66f198b653        4 days ago          1.43GB
inklabsfoundation/inkchain-javaenv       x86_64-0.13.0       cc66f198b653        4 days ago          1.43GB
inklabsfoundation/inkchain-ca            latest              c216e43656f1        2 months ago        238MB
inklabsfoundation/inkchain-ca            x86_64-0.13.0       c216e43656f1        2 months ago        238MB
inklabsfoundation/inkchain-kafka         latest              01c020e46c87        3 months ago        1.31GB
inklabsfoundation/inkchain-kafka         x86_64-0.13.0       01c020e46c87        3 months ago        1.31GB
inklabsfoundation/inkchain-couchdb       latest              56babe466471        3 months ago        1.51GB
inklabsfoundation/inkchain-couchdb       x86_64-0.13.0       56babe466471        3 months ago        1.51GB
inklabsfoundation/inkchain-zookeeper     latest              720c4368c2f5        3 months ago        1.33GB
inklabsfoundation/inkchain-zookeeper     x86_64-0.13.0       720c4368c2f5        3 months ago        1.33GB
inklabsfoundation/inkchain-baseimage     x86_64-0.3.1        9f2e9ec7c527        7 months ago        1.268 GB
inklabsfoundation/inkchain-baseos        x86_64-0.3.1        4b0cab202084        7 months ago        156.6 MB
```

## Run the examples

We provide an example [token issue](../examples/cli_test) which based on inkchain.
