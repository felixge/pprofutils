# Notes to Self

The following env secrets need to be set on the agent app:

- `DD_API_KEY`
- `DD_SITE`

And for the pprofutils http serve process:

- `DD_AGENT_HOST`: Host of the agent, e.g. dd-agent-pprof-to.internal.
- `DD_ENV`: Name of the env, e.g. 'prod'.
- `DD_SERVICE` Name of the service, e.g. 'pprof.to'.
