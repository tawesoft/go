# xff - DirectX (.x) file format decoder

## About

Package xff implements a decoder for the DirectX (.x) file format in
Go (Golang).

THIS IS A PREVIEW RELEASE. Currently the parser can just validate a document
but doesn't return useful data. The full implementation will follow over
the next day or so.

This parser is a complete implementation that supports user-defined templates.
It should be able to parse *every* well-formed DirectX (.x) file. Get in touch
if you find any that don't parse correctly.

This is free and open source software. If you find this module useful enough
that you use it in a commercial product, it is suggested that you donate a
sum to the author equal to the cost of one copy of your product. To do so,
make a donation at https://www.paypal.me/TawesoftLtd and forward a copy of your
receipt to "ben+xff-donation@tawesoft.co.uk". Alternatively you can purchase
commercial support through open-source@tawesoft.co.uk. This is, of course,
entirely optional.

DirectX is a registered trademark of Microsoft Corporation in the United States
and/or other countries.

|  Links  | License | Stable? | 
|:-------:|:-------:|:-------:| 
| [home][home_] ∙ [docs][docs_] ∙ [src][src_] | [MIT][copy_] | ✘ **no** |

[home_]: https://tawesoft.co.uk/go/xff
[src_]:  https://github.com/tawesoft/go/tree/master/xff
[docs_]: https://godoc.org/tawesoft.co.uk/go/xff
[copy_]: https://github.com/tawesoft/go/tree/master/xff/_COPYING.md

## Download

```shell script
go get -u tawesoft.co.uk/go
```

## Import

```
import tawesoft.co.uk/go/xff
```