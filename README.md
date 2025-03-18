# lazyenv

cli tool to non-destructively set `.env` vars @ P2P.me from contracts deployments to the UI repos, the lazy way.

# TODOs

- [x] Setup config structure
- [x] Setup assert internals
- [x] Read from a config
- [ ] Read the env vars from source repo's using what method? REGEX? MANUAL PARSING? OR WHAT?
- [ ] Have it in some kind of structure in memory. PROLLY A MAP.
- [ ] Read the dest .env files one by one. Find the best way to structure the content.
- [ ] Comment out the currently used env variables in the envMapping
- [ ] Add a new env entry for all the vars in the envMapping - w/ some metadata like chain, date, etc.
