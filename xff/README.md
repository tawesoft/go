# xff - DirectX (.x) file format decoder

## About

Package xff implements a decoder for the DirectX (.x) file format in
Go (Golang).

THIS IS A PREVIEW RELEASE. A few templates are missing and the returned data
is not decoded from bytes to types properly. The testdata is missing. It needs
a refactor. This will be fixed early next week (beginning Feb 24th 2020).

This parser is (or aims to be) a complete implementation. It supports user-defined
templates and should be able to parse *every* well-formed DirectX (.x) file.

There are a few features not yet implemented - not because they are difficult
to implement, but because I can't find any real-world examples using that
feature to test against. Things like object referencing by UIID instead of
name, or the binary encoded DirectX .x file format, arrays of strings, or
multidimensional arrays. Send in a sample if you find something that doesn't
work (and should) and it'll get fixed quickly!

This is free and open source software. It took about five days to write - so
if you find this module useful enough that you use it in a commercial product then
there is a polite expectation that you donate a sum to the author equal to the
retail cost of one copy of your product. To do so, make a payment at
https://www.paypal.me/TawesoftLtd and forward a copy of your receipt to
"xff-thanks@tawesoft.co.uk". Alternatively you can purchase commercial support
through open-source@tawesoft.co.uk. If you're making a computer game I'd
appreciate one free copy, too. This is, of course, entirely optional.

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