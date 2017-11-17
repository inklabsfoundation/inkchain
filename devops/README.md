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
inkchain/inkchain-ccenv        latest                          249e938219b0        29 hours ago        1.28GB
inkchain/inkchain-ccenv        x86_64-0.9.1                    249e938219b0        29 hours ago        1.28GB
inkchain/inkchain-tools        latest                          b19f13f153ef        46 hours ago        1.34GB
inkchain/inkchain-tools        x86_64-0.9.1                    b19f13f153ef        46 hours ago        1.34GB
inkchain/inkchain-testenv      latest                          41b6a64231b7        46 hours ago        1.41GB
inkchain/inkchain-testenv      x86_64-0.9.1                    41b6a64231b7        46 hours ago        1.41GB
inkchain/inkchain-buildenv     latest                          4d64e0249a72        46 hours ago        1.32GB
inkchain/inkchain-buildenv     x86_64-0.9.1                    4d64e0249a72        46 hours ago        1.32GB
inkchain/inkchain-orderer      latest                          967234ea4d84        46 hours ago        180MB
inkchain/inkchain-orderer      x86_64-0.9.1                    967234ea4d84        46 hours ago        180MB
inkchain/inkchain-peer         latest                          e046ec9cc594        46 hours ago        183MB
inkchain/inkchain-peer         x86_64-0.9.1                    e046ec9cc594        46 hours ago        183MB
inkchain/inkchain-javaenv      latest                          5616849dda4a        46 hours ago        1.42GB
inkchain/inkchain-javaenv      x86_64-0.9.1                    5616849dda4a        46 hours ago        1.42GB
inkchain/inkchain-baseimage    x86_64-0.3.1                    9f2e9ec7c527        6 months ago        1.27GB
inkchain/inkchain-baseos       x86_64-0.3.1                    4b0cab202084        6 months ago        157MB
inkchain/inkchain-couchdb      latest                          f960e952b103        2 weeks ago         1.51GB
inkchain/inkchain-couchdb      x86_64-0.9.1                    f960e952b103        2 weeks ago         1.51GB
inkchain/inkchain-kafka        latest                          e73bcd0d91ab        2 weeks ago         1.3GB
inkchain/inkchain-kafka        x86_64-0.9.1                    e73bcd0d91ab        2 weeks ago         1.3GB
inkchain/inkchain-zookeeper    latest                          9da7d54b7d6f        2 weeks ago         1.31GB
inkchain/inkchain-zookeeper    x86_64-0.9.1                    9da7d54b7d6f        2 weeks ago         1.31GB
```

## Run the examples

We provide an example [token issue](../examples/issue-token) which based on inkchain.
