## Development Environment

This folder contains the files which are used for bootstrapping the Inklabsfoundation Inkchain development environment.
and so far, we recommend using ubuntu 16.04.

We provide a [env setup script](./setup-env.sh) to config your environment.

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
REPOSITORY                     TAG                 IMAGE ID            CREATED             SIZE
inkchain/inkchain-tools        latest              7b18be3a82fe        4 days ago          1.34GB
inkchain/inkchain-tools        x86_64-0.10         7b18be3a82fe        4 days ago          1.34GB
inkchain/inkchain-couchdb      latest              eccce4288cf5        4 days ago          1.51GB
inkchain/inkchain-couchdb      x86_64-0.10         eccce4288cf5        4 days ago          1.51GB
inkchain/inkchain-kafka        latest              1bb4cd5bba69        4 days ago          1.3GB
inkchain/inkchain-kafka        x86_64-0.10         1bb4cd5bba69        4 days ago          1.3GB
inkchain/inkchain-zookeeper    latest              4eb257be4062        4 days ago          1.33GB
inkchain/inkchain-zookeeper    x86_64-0.10         4eb257be4062        4 days ago          1.33GB
inkchain/inkchain-orderer      latest              85cab7e9c4cf        4 days ago          180MB
inkchain/inkchain-orderer      x86_64-0.10         85cab7e9c4cf        4 days ago          180MB
inkchain/inkchain-peer         latest              0a7d33a7f6f0        4 days ago          183MB
inkchain/inkchain-peer         x86_64-0.10         0a7d33a7f6f0        4 days ago          183MB
inkchain/inkchain-javaenv      latest              09231aa0c8f0        4 days ago          1.42GB
inkchain/inkchain-javaenv      x86_64-0.10         09231aa0c8f0        4 days ago          1.42GB
inkchain/inkchain-ccenv        latest              53566e6b8ae2        4 days ago          1.29GB
inkchain/inkchain-ccenv        x86_64-0.10         53566e6b8ae2        4 days ago          1.29GB
inkchain/inkchain-baseimage    x86_64-0.3.1        c92d9fdee998        2 months ago        1.26GB
inkchain/inkchain-baseos       x86_64-0.3.1        4b0cab202084        6 months ago        157MB
```

## Run the examples

We provide an example [token issue](../examples/issue-token) which based on inkchain.
