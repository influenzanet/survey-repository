

# 0.4

- Add a dedicated login endpoint providing an short lived token for others endpoints


# O.3.1

- Update dependencies fiber & x/crypto
- Add explore page in UI with survey rendering
- Fix return status
- Use provided name when import survey (allow to override survey key)

# 0.3

- Add basic web UI
- Add basic stat endpoint /namespace/:namespace/surveys/stats
- Add modelType in meta model (D=Definition, P=Preview), both version can be stored (with same version)
- Add shorter structure to list imported versions at /namespace/:namespace/surveys/versions (`ShortVersionMeta`)
- Import returns more information (as `ShortVersionMeta`)

# 0.2

- Update dependencies
- Refactor config with sections
- Fix query surveys accepting platforms and names
- allow to provide a version numer in request, but mandatory if survey version is not in json

