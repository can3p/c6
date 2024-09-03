C6
===========================

C6 is a SASS 3.2 compatible implementation written in Go.

This is a fork of the original project with no intention for backwards compatibility.
Current goal is to have at least some success against running the official spec.

If you want to help, the process is the following:

- Decide which spec you want to work on in the `sass-spec`
- Remove the corresponding entries from [blacklist](testspec/failures_list.go)
- run the spec and see it fail with:
  ```
  IGNORE_BLACKLISTED=true go test testspec/*
  ```
- Fix the code and create a pull request

Game rules are simple:
- If a spec was successfully removed from the blacklist with the previous changes it cannot be added back
- That's it

## Working in progress

- [ ] Lexing
  - [x] `@import`
  - [x] simple selector.
    - [x] type selector.
    - [x] child selector.
    - [x] attribute selector.
    - [x] adjacent selector.
    - [x] descendant selector.
    - [x] class selector.
    - [x] ID selector.
  - [x] Ruleset
  - [x] Sub-ruleset
  - [x] Interpolation
  - [x] Property name
  - [x] Property value list
  - [x] Nested properties.
  - [x] Comma-separated list
  - [x] Space-separated list
  - [x] `@if` , `@else` , `@else if`
  - [x] `@for`, `from`, `through` statement
  - [x] `@while`
  - [x] `@mixin`
  - [x] `@mixin` with arguments
  - [x] `@include`
  - [x] Flags: `!default`, `!important`, `!optional`, `!global`
  - [x] Hex color
  - [x] Functions
  - [x] Vendor prefix properties
  - [x] MS filter.  `progid:DXImageTransform.Microsoft....`
  - [x] Variable names
  - [x] Variable assignment
  - [x] Time unit support. `s` second, `ms` ... etc
  - [x] Angle unit support.
  - [x] Resolution unit support.
  - [x] Unicode Range support: <https://developer.mozilla.org/en-US/docs/Web/CSS/unicode-range>
  - [x] Media Query
- [ ] Built-in Functions
  - .... to be listed
- [ ] Parser
  - [x] Parse `@import`
  - [x] Parse Expr
  - [x] Parse Space-Sep List
  - [x] Parse Comma-Sep List
  - [x] Parse Map (tests required)
  - [x] Parse Selector
  - [ ] Parse Selector with interpolation
  - [x] Parse RuleSet
  - [x] Parse DeclBlock
  - [x] Parse Variable Assignment Stmt
  - [x] Parse PropertyName
  - [x] Parse PropertyName with interpolation
  - [x] Parse PropertyValue
  - [x] Parse PropertyValue with interpolation
  - [x] Parse conditions
  - [x] Parse `@media` statement
  - [x] Parse Nested RuleSet
  - [x] Parse Nested Properties
  - [x] Parse options: `!default`, `!global`, `!optional`
  - [ ] Parse CSS Hack for different browser (support more syntax sugar for this)
  - [x] Parse `@font-face` block
  - [x] Parse `@if` statement
  - [x] Parse `@for` statement
  - [x] Parse `@while` statement
  - [x] Parse `@mixin` statement
  - [x] Parse `@include` statement
  - [x] Parse `@function` statement
  - [x] Parse `@return` statement
  - [x] Parse `@extend` statement
  - [x] Parse keyword arguments for `@function`
  - [ ] Parse `@switch` statement
  - [ ] Parse `@case` statement
  - [ ] Parse `@use` statement
- [ ] Building AST
  - [x] RuleSet
  - [x] DeclBlock
  - [x] PropertyName
  - [x] PropertyValue
  - [x] Comma-Separated List
  - [x] Space-Separated List
  - [x] Basic Exprs
  - [x] FunctionCall
  - [x] Expr with interpolation
  - [x] Variable statements
  - [x] Built-in color keyword table
  - [x] Hex Color computation
  - [x] Number operation: add, sub, mul, div
  - [x] Length operation: number operation for px, pt, em, rem, cm ...etc
  - [x] Expr evaluation
  - [x] Boolean expression evaluation
  - [x] Media Query conditions
  - [x] `@if` If Condition
  - [x] `@else if` If Else If
  - [x] `@else` else condition
  - [x] `@while` statement node
  - [x] `@function` statement node
  - [x] `@mixin` statement node
  - [x] `@include` statement node
  - [x] `@return` statement node
  - [ ] `@each` statement node

- [ ] Runtime
  - [ ] HSL Color computation
  - [ ] Function Call Invoke mech
  - [ ] Mixin Include
  - [ ] Import

- [ ] SASS Built-in Functions
  - [ ] RGB functions
    - [ ] `rgb($red, $green, $blue)`
    - [ ] `rgba($red, $green, $blue, $alpha)`
    - [ ] `red($color)`
    - [ ] `green($color)`
    - [ ] `blue($color)`
    - [ ] `mix($color1, $color2, [$weight])`
  - [ ] HSL Functions
    - [ ] `hsl($hue, $saturation, $lightness)`
    - [ ] `hsla($hue, $saturation, $lightness, $alpha)`
    - [ ] `hue($color)`
    - [ ] `saturation($color)`
    - [ ] `lightness($color)`
    - [ ] `adjust-hue($color, $degrees)`
    - [ ] `lighten($color, $amount)`
    - [ ] `darken($color, $amount)`
    - [ ] `saturate($color, $amount)`
    - [ ] `desaturate($color, $amount)`
    - [ ] `grayscale($color)`
    - [ ] `complement($color)`
    - [ ] `invert($color)`
  - [ ] Opacity Functions
    - [ ] `alpha($color) / opacity($color)`
    - [ ] `rgba($color, $alpha)`
    - [ ] `opacify($color, $amount) / fade-in($color, $amount)`
    - [ ] `transparentize($color, $amount) / fade-out($color, $amount)`
  - [ ] Other Color Functions
    - [ ] `adjust-color($color, [$red], [$green], [$blue], [$hue], [$saturation], [$lightness], [$alpha])`
    - [ ] `scale-color($color, [$red], [$green], [$blue], [$saturation], [$lightness], [$alpha])`
    - [ ] `change-color($color, [$red], [$green], [$blue], [$hue], [$saturation], [$lightness], [$alpha])`
    - [ ] `ie-hex-str($color)`
  - [ ] String Functions
    - [ ] `unquote($string)`
    - [ ] `quote($string)`
    - [ ] `str-length($string)`
    - [ ] `str-insert($string, $insert, $index)`
    - [ ] `str-index($string, $substring)`
    - [ ] `str-slice($string, $start-at, [$end-at])`
    - [ ] `to-upper-case($string)`
    - [ ] `to-lower-case($string)`
  - [ ] Number Functions
    - [ ] `percentage($number)`
    - [ ] `round($number)`
    - [ ] `ceil($number)`
    - [ ] `floor($number)`
    - [ ] `abs($number)`
    - [ ] `min($numbers…)`
    - [ ] `max($numbers…)`
    - [ ] `random([$limit])`
  - [ ] List Functions
    - [ ] `length($list)`
    - [ ] `nth($list, $n)`
    - [ ] `set-nth($list, $n, $value)`
    - [ ] `join($list1, $list2, [$separator])`
    - [ ] `append($list1, $val, [$separator])`
    - [ ] `zip($lists…)`
    - [ ] `index($list, $value)`
    - [ ] `list-separator(#list)`
  - [ ] Map Functions
    - [ ] `map-get($map, $key)`
    - [ ] `map-merge($map1, $map2)`
    - [ ] `map-remove($map, $keys…)`
    - [ ] `map-keys($map)`
    - [ ] `map-values($map)`
    - [ ] `map-has-key($map, $key)`
    - [ ] `keywords($args)`
  - [ ] Selector Functions
    - .... to be expanded ...

- [ ] CodeGen
  - [ ] CompactCompiler
    - [ ] CompileCssImportStmt: `@import url(...);`
    - [ ] CompileRuleSet
    - [ ] CompileSelectors
      - [ ] CoimpileParentSelector
    - [ ] CompileSubRuleSet
    - [ ] CompileCommentBlock
    - [ ] CompileDeclBlock
    - [ ] CompileMediaQuery: `@media`
    - [ ] CompileSupportrQuery: `@support`
    - [ ] CompileFontFace: `@support`
    - [ ] CompileForStmt
    - [ ] CompileIfStmt
      - [ ] CompileElseIfStmt
    - [ ] CompileWhileStmt
    - [ ] CompileEachStmt
    - [ ] ... list more ast nodes here ...

- [ ] Syntax
  - [ ] built-in `@import-once`

<!--
## Features

- [ ] import directory: https://github.com/sass/sass/issues/690
- [ ] import css as sass: https://github.com/sass/sass/issues/556
- [ ] import once: https://github.com/sass/sass/issues/139
- [ ] namespace and alias: https://github.com/sass/sass/issues/353
- [ ] `@use` directive: https://github.com/nex3/sass/issues/353#issuecomment-5146513
- [ ] conditional import: https://github.com/sass/sass/issues/451
- [ ] `@sprite` syntax sugar
-->

## Ambiguity

The original design of SASS contains a lot of grammar ambiguity.

for example, as SASS uses interpolation:

```
#{$name}:before {

}
```

Since nested properties are allowed, in the above code, we don't know if it's a
selector or a property namespace if we don't know the `$name` variable.

Where `before` might be a property value or a part of the selector.

links:
- <https://www.facebook.com/cindylinz/posts/10202186527405801?hc_location=ufi>
- <https://www.facebook.com/yoan.lin/posts/10152968537931715?_rdr>

To handle this kind of interpolation, we define a type of token named template:

```
#{$name}:before a {

}
```

In the above code, `#{$name}:before a` is treated as `T_SELECTOR_TEMPLATE` token, which
type will be resolved at the runtime.


## Reference

SASS Reference <http://sass-lang.com/documentation/file.SASS_REFERENCE.html>


A feature check list from libsass:

- <https://github.com/sass/libsass/releases/tag/3.2.0>
- <https://github.com/sass/sass/issues/1094>


Standards:

- CSS Syntax Module Level 3 <http://www.w3.org/TR/css-syntax-3>
- CSS 3 Selector <http://www.w3.org/TR/css3-selectors/#grouping>
- CSS Font <http://www.w3.org/TR/css3-fonts/#basic-font-props>
- Selectors API <http://www.w3.org/TR/selectors-api/>
- At-Page Rule <http://dev.w3.org/csswg/css-page-3/#at-page-rule>
- Railroad diagram <https://github.com/tabatkins/railroad-diagrams>
- CSS 2.1 Grammar <http://www.w3.org/TR/CSS21/grammar.html>

SASS/CSS Frameworks, libraries:

- Bourbon <http://bourbon.io/>
- Marx <https://github.com/mblode/marx>
- FormHack <http://formhack.io/>
- Susy <http://susy.oddbird.net/>
- Gumby <http://www.gumbyframework.com/>

Articles:

- Logic in media queries - <https://css-tricks.com/logic-in-media-queries/>

## Credits

Original source code is available at https://github.com/c9s/c6

## License

MPL License <https://www.mozilla.org/MPL/2.0/>

(MPL is like LGPL but with static/dynamic linking exception, which allows you
to either dynamic/static link this library without permissions)
