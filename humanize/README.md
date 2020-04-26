# humanize - lightweight human-readable numbers

## About

Package humanize implements lightweight human-readable numbers for Go.

Why use this package?

* Much more efficient than the leading package (parse input about 5x faster, fewer memory allocations)

* Very flexible formatting and internationalisation means better numbers for your humans e.g.
this package even supports the Indian System of Numeration for lakh and crore digit grouping. 

Compare tawesoft/go/humanize vs dustin/go-humanize: Tawesoft's parses input about 5 times faster with fewer memory
allocations, and formats output about 25% quicker than dustin's. But dustin's is so far more complete, older, has a
stable API, and has been tested by more people.

|  Links  | License | Stable? | 
|:-------:|:-------:|:-------:| 
| [home][home_] ∙ [docs][docs_] ∙ [src][src_] | [MIT-0][copy_] | ✘ **no** |

[home_]: https://tawesoft.co.uk/go/humanize
[src_]:  https://github.com/tawesoft/go/tree/master/humanize
[docs_]: https://godoc.org/tawesoft.co.uk/go/humanize
[copy_]: https://github.com/tawesoft/go/tree/master/humanize/_COPYING.md

## Download

```shell script
go get -u tawesoft.co.uk/go
```

## Import

```
import tawesoft.co.uk/go/humanize
```