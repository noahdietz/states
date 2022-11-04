# protoc-gen-states

Collects all enums and their values defined within the input proto package
that (may) represent a [State](https://google.aip.dev/216). Output is CSV format
and will be written to `stderr` if the `out_file` parameter isn't specified.
Multiple runs using the same `out_file` appends the results of the subsequent
runs to the existing contents.

Run `./analyze.sh` in the projects directory to get results. This script is
specifically meant to analyze all proto packages in `google/cloud` of
[googleapis][].

Clone [googleapis] and export the variable `GOOGLEAPIS` in your shell to avoid
repetitive downloads. The script will download and set the variable itself
if unset.

[googleapis]: https://github.com/googleapis
