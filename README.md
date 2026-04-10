# Proxy for MATLAB Parallel Server&trade;

Use the parallelserverproxy to start a proxy server which can proxy all traffic between a MATLAB&reg; client and a MATLAB Parallel Server cluster.
This can be used to create a single access point for the cluster on the host where the parallelserverproxy is run and via the port specified.

The parallelserverproxy uses the SOCKS protocol to allow connecting clients to specify the destination to connect to. Authentication of connecting clients is provided by
mutual TLS (mTLS) using client certificates. Certificate files can be generated using the mjssetup tool available from (https://github.com/mathworks/mjssetup).

You do not require access to a MATLAB installation.

## Installation

You can download pre-compiled binaries for Linux&reg; and Windows&reg; from the [Releases](https://github.com/mathworks/parallelserverproxy/releases) page.

### Installation on Linux

1. Navigate to the [Releases](https://github.com/mathworks/parallelserverproxy/releases) page.
2. Download the latest `parallelserverproxy-glnxa64.tar.gz` file from the assets section of the latest release.
3. To extract the binary, in the terminal, run `tar -xzf parallelserverproxy-glnxa64.tar.gz`.

### Installation on Windows

1. Navigate to the [Releases](https://github.com/mathworks/parallelserverproxy/releases) page.
2. Download the latest `parallelserverproxy-win64.zip` file from the assets section of the latest release.
3. To extract the binary, unzip the `parallelserverproxy-win64.zip` file.

## Usage

`parallelserverproxy [<args>]` starts a proxy server using the specified input arguments.
- `args` - Inputs to the proxy server.

Once connected, the url clients use to connect to the proxy is printed. For example:
```
parallelserverproxy -port 1080 -certificate "proxy-certificate.json"
Starting SOCKS5 proxy on socks5s://*:1080
```
Connecting clients should then specify the proxy url as `socks5s://<proxyHost>:1080` to access the cluster via the proxy, where
proxyHost should be replaced with the resolvable hostname or IP address of the machine running the proxy server.

To display the help text for parallelserverproxy, run
```
parallelserverproxy -help
```

### Examples

Start a `parallelserverproxy` on all network interfaces on the default port (1080).
Pass the certificate file to use during the mTLS handshake with the `-certificate` argument. 
```
parallelserverproxy -certificate "proxy-certificate.json"
```
For example, the "proxy-certificate.json" file is generated using the `mjssetup` tool as follows:
```
mjssetup generate-certificate -secretfile "secret.json" -outfile "proxy-certificate.json"
```
Any connecting client must have a corresponding client certificate file signed with the same secret file to authenticate
during the mTLS handshake.

Start a `parallelserverproxy` in insecure mode on all network interfaces on the default port (1080).
Disable both encryption and authentication on client-to-proxy connections using the `-disableMutualTLS` argument.
This will issue a warning since there is now no authentication of clients and exposed services may now be insecure. Only use
this option in secure closed networks with trusted clients to avoid any performance impact of mTLS.
```
parallelserverproxy -disableMutualTLS
```
Users must either provide a certificate file using the `-certificate` argument to use mTLS authentication or explicitly
disable mTLS using the `-disableMutualTLS` argument.

Start a `parallelserverproxy` on a specified network interface and port.
```
parallelserverproxy -host localhost -port 1080 -certificate "proxy-certificate.json"
```

## Build Proxy from Source Code

To download a zip file of this repository, at the top of this repository page, select Code > Download ZIP.
Alternatively, to clone this repository to your computer with Git installed, run the following command on your operating system's command line:

```
git clone https://github.com/mathworks/parallelserverproxy
```

To compile the parallelserverproxy executable from the source code, you must use Go version 1.23 or later.
Use Go to compile the parallelserverproxy executable:
```
go build -o parallelserverproxy main.go
```

## License

The license is available in the [license.txt](license.txt) file in this repository.

## Community Support

[MATLAB Central](https://www.mathworks.com/matlabcentral)

## Technical Support

If you require assistance or have a request for additional features or capabilities, contact [MathWorks Technical Support](https://www.mathworks.com/support/contact_us.html).

Copyright 2024-2025 The MathWorks, Inc.
