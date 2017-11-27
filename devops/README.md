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
inklabsfoundation/inkchain-tools       latest              5420d9596cbe        4 days ago          1.342 GB
inklabsfoundation/inkchain-tools       x86_64-0.10         5420d9596cbe        4 days ago          1.342 GB
inklabsfoundation/inkchain-couchdb     latest              b5c33a9fe7a6        4 days ago          1.514 GB
inklabsfoundation/inkchain-couchdb     x86_64-0.10         b5c33a9fe7a6        4 days ago          1.514 GB
inklabsfoundation/inkchain-kafka       latest              92008d73b9d1        4 days ago          1.309 GB
inklabsfoundation/inkchain-kafka       x86_64-0.10         92008d73b9d1        4 days ago          1.309 GB
inklabsfoundation/inkchain-zookeeper   latest              f6a70097a8f7        4 days ago          1.327 GB
inklabsfoundation/inkchain-zookeeper   x86_64-0.10         f6a70097a8f7        4 days ago          1.327 GB
inklabsfoundation/inkchain-orderer     latest              4fc46e39dce7        4 days ago          180.1 MB
inklabsfoundation/inkchain-orderer     x86_64-0.10         4fc46e39dce7        4 days ago          180.1 MB
inklabsfoundation/inkchain-peer        latest              a4e2ecb84b3b        4 days ago          183.4 MB
inklabsfoundation/inkchain-peer        x86_64-0.10         a4e2ecb84b3b        4 days ago          183.4 MB
inklabsfoundation/inkchain-javaenv     latest              f71fe027949a        4 days ago          1.425 GB
inklabsfoundation/inkchain-javaenv     x86_64-0.10         f71fe027949a        4 days ago          1.425 GB
inklabsfoundation/inkchain-ccenv       latest              09a38d154647        4 days ago          1.293 GB
inklabsfoundation/inkchain-ccenv       x86_64-0.10         09a38d154647        4 days ago          1.293 GB
inkchain/inkchain-baseimage            x86_64-0.3.1        c92d9fdee998        2 months ago        1.26GB
inkchain/inkchain-baseos               x86_64-0.3.1        4b0cab202084        6 months ago        157MB
```

## Run the examples

We provide an example [token issue](../examples/issue-token) which based on inkchain.
