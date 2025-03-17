# lazyenv

cli tool to non-destructively set `.env` vars @ P2P.me from contracts deployments to the UI repos.

This tool should:
- read from a JSON config file that sets various configurations, do the needful, set the env vars.
- the config should consist of
 - src repo
 - cmd to run in src repo
 - dst repos
 - env var mapping
- easily extendable in the future
