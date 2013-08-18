# m

Metrics server designed to make storing metrics easy. Simple set up graphite-like tool. Eventually will allow for running multiple instances and sharing load.

## Status

Tool is designed for use with Riak as a backend. This is due to the fact the author is currently developing the Go Riak client (rgo). The idea is to make the backend plugable but this will not be the focus of the first release.

## Usage

Run the tool using `./m`

Flags : 

	-interval	:	interval for the servers to use
	-port		:	signal which port you want to listen on

## License

MIT

