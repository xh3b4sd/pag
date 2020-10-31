package typescript

const indexTemplate = `//
// Do not edit. This file was generated via the "pag" command line tool. More
// information about the tool can be found at github.com/xh3b4sd/pag.
//
//     pag generate typescript
//
{{ range $r := . }}
// -------------------------------------------------------------------------- //

import * as {{ $r.Dir | ToResource }}Client  from "./{{ $r.Dir }}/ApiServiceClientPb";
import * as {{ $r.Dir | ToResource }}Create  from "./{{ $r.Dir }}/create_pb";
import * as {{ $r.Dir | ToResource }}Delete  from "./{{ $r.Dir }}/delete_pb";
import * as {{ $r.Dir | ToResource }}Search  from "./{{ $r.Dir }}/search_pb";
import * as {{ $r.Dir | ToResource }}Update  from "./{{ $r.Dir }}/update_pb";

export const {{ $r.Dir | ToResource }} = {
  Client:  {{ $r.Dir | ToResource }}Client.APIClient,
  Create: {
    I: {{ $r.Dir | ToResource }}Create.CreateI,
    O: {{ $r.Dir | ToResource }}Create.CreateO,
  },
  Delete: {
    I: {{ $r.Dir | ToResource }}Delete.DeleteI,
    O: {{ $r.Dir | ToResource }}Delete.DeleteO,
  },
  Search: {
    I: {{ $r.Dir | ToResource }}Search.SearchI,
    O: {{ $r.Dir | ToResource }}Search.SearchO,
  },
  Update: {
    I: {{ $r.Dir | ToResource }}Update.UpdateI,
    O: {{ $r.Dir | ToResource }}Update.UpdateO,
  },
}

// -------------------------------------------------------------------------- //


{{ end -}}
`
