package typescript

const indexTemplate = `#
# Do not edit. This file was generated via the "pag" command line tool. More
# information about the tool can be found at github.com/xh3b4sd/pag.
#
#     pag generate typescript
#
{{ range $r := . }}
// -------------------------------------------------------------------------- //

import * as {{ $r.Dir | ToExport }}Client  from "./{{ $r.Dir }}/ApiServiceClientPb";
import * as {{ $r.Dir | ToExport }}Create  from "./{{ $r.Dir }}/create_pb";
import * as {{ $r.Dir | ToExport }}Delete  from "./{{ $r.Dir }}/delete_pb";
import * as {{ $r.Dir | ToExport }}Search  from "./{{ $r.Dir }}/search_pb";
import * as {{ $r.Dir | ToExport }}Update  from "./{{ $r.Dir }}/update_pb";

export const {{ $r.Dir | ToExport }} = {
  Client:  {{ $r.Dir | ToExport }}Client.APIClient,
  Create: {
    I: {{ $r.Dir | ToExport }}Create.CreateI,
    O: {{ $r.Dir | ToExport }}Create.CreateO,
  },
  Delete: {
    I: {{ $r.Dir | ToExport }}Delete.DeleteI,
    O: {{ $r.Dir | ToExport }}Delete.DeleteO,
  },
  Search: {
    I: {{ $r.Dir | ToExport }}Search.SearchI,
    O: {{ $r.Dir | ToExport }}Search.SearchO,
  },
  Update: {
    I: {{ $r.Dir | ToExport }}Update.UpdateI,
    O: {{ $r.Dir | ToExport }}Update.UpdateO,
  },
}

// -------------------------------------------------------------------------- //


{{ end -}}
`
