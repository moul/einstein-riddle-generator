# einstein-riddle-generator
:game_die: Einstein Riddle Generator

[![GoDoc](https://godoc.org/github.com/moul/einstein-riddle-generator?status.svg)](https://godoc.org/github.com/moul/einstein-riddle-generator)

Check out the [demo on Appspot](http://einsteins-riddle.appspot.com) !

## Install

```console
$ go get github.com/moul/einstein-riddle-generator/...
```

## Example

```console
$ einstein-riddle
Riddle
------
- job:cop == pet:shark
- room:kitchen == nationality:scottish
- house-color:red == pet:beaver
- nationality:spannish == house-color:pink
- job:scientist == room:bedroom
- room:restroom == job:designer
- job:teacher == pet:poney
- pet:bird is direct on the right of job:architect
- nationality:american == room:bathroom
- nationality:french == room:garden
- house-color:yellow is on the left of nationality:english
- house-color:blue is in the middle
- room:bedroom == nationality:english
- pet:beaver is on the right of room:bathroom
- job:architect is direct on the right of nationality:french
- job:cop is on the far left
- pet:shark is on the far left
- nationality:american is on the far left

- where is pet:snake ?
- where is house-color:magenta ?

Answer
------
              1             2            3              4             5
room          2: bathroom   1: garden    1: kitchen     1: restroom   2: bedroom
job           2: cop        1: teacher   2: architect   1: designer   1: scientist
pet           2: shark      1: poney     0: snake       1: bird       2: beaver
house-color   0: magenta    1: yellow    1: blue        1: pink       1: red
nationality   2: american   2: french    1: scottish    1: spannish   2: english
```

## License

MIT
