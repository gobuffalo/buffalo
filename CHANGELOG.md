# Change Log

## [v0.7.1](https://github.com/gobuffalo/buffalo/tree/v0.7.1) (2017-01-13)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.7.0...v0.7.1)

**Closed issues:**

- Channel for community chat [\#126](https://github.com/gobuffalo/buffalo/issues/126)
- build\_path does not work for Windows [\#124](https://github.com/gobuffalo/buffalo/issues/124)
- Installation error [\#120](https://github.com/gobuffalo/buffalo/issues/120)
- Default js and css file when webpack is skipped [\#116](https://github.com/gobuffalo/buffalo/issues/116)
- go.rice requirement in build [\#115](https://github.com/gobuffalo/buffalo/issues/115)
- Warn user about SESSION\_SECRET variable? [\#114](https://github.com/gobuffalo/buffalo/issues/114)
- Buffalo sub-commands unavailable [\#113](https://github.com/gobuffalo/buffalo/issues/113)
- Error installing buffalo.  [\#111](https://github.com/gobuffalo/buffalo/issues/111)
- In `production` mode Buffalo/Velvet templates show traces [\#110](https://github.com/gobuffalo/buffalo/issues/110)
- buffalo new \<project\> fails looking for golang.org/x/tools/go/gcimporter [\#108](https://github.com/gobuffalo/buffalo/issues/108)
- missing "public/assets" box on new app without webpack [\#104](https://github.com/gobuffalo/buffalo/issues/104)
- SHOULDERS update & grift task [\#99](https://github.com/gobuffalo/buffalo/issues/99)
- \[minor\] dev mode on non-buffalo project results in panic [\#91](https://github.com/gobuffalo/buffalo/issues/91)
- typo in generated database.yml? [\#87](https://github.com/gobuffalo/buffalo/issues/87)
- Buffalo dev not starting [\#86](https://github.com/gobuffalo/buffalo/issues/86)
- Export fileResolver Field in Render Options Struct  [\#84](https://github.com/gobuffalo/buffalo/issues/84)
- `buffalo task` should forward to `grift` [\#59](https://github.com/gobuffalo/buffalo/issues/59)
- generate a default .codeclimate.yml file for new projects [\#37](https://github.com/gobuffalo/buffalo/issues/37)
- generate a README.md for new projects [\#35](https://github.com/gobuffalo/buffalo/issues/35)
- add a form generator to helper [\#19](https://github.com/gobuffalo/buffalo/issues/19)
- Don't write test.log files when running tests [\#17](https://github.com/gobuffalo/buffalo/issues/17)
- Add an "actions" generator [\#16](https://github.com/gobuffalo/buffalo/issues/16)

**Merged pull requests:**

- Add badge for Go Report Card to README [\#132](https://github.com/gobuffalo/buffalo/pull/132) ([stuartellis](https://github.com/stuartellis))
- Makes our tests run on Go 1.7 and 1.8 [\#131](https://github.com/gobuffalo/buffalo/pull/131) ([apaganobeleno](https://github.com/apaganobeleno))
- build\_path does not work for Windows closes \#124 [\#130](https://github.com/gobuffalo/buffalo/pull/130) ([markbates](https://github.com/markbates))
- Edit some typo [\#129](https://github.com/gobuffalo/buffalo/pull/129) ([IvanMenshykov](https://github.com/IvanMenshykov))
- Passing some issues from codeclimate [\#122](https://github.com/gobuffalo/buffalo/pull/122) ([apaganobeleno](https://github.com/apaganobeleno))
- Provide a mechanism to map status codes to error handles. Closes \#110 [\#121](https://github.com/gobuffalo/buffalo/pull/121) ([markbates](https://github.com/markbates))
- Warn user about SESSION\_SECRET variable? closes \#114 [\#119](https://github.com/gobuffalo/buffalo/pull/119) ([markbates](https://github.com/markbates))
- point people to npm docs if there is an issue running npm [\#118](https://github.com/gobuffalo/buffalo/pull/118) ([markbates](https://github.com/markbates))
- Default css js files and no logo in assets wo webpack [\#117](https://github.com/gobuffalo/buffalo/pull/117) ([fooflare](https://github.com/fooflare))
- Return the RouteInfo when mapping an endpoint. Also make it available in the request context [\#109](https://github.com/gobuffalo/buffalo/pull/109) ([markbates](https://github.com/markbates))
- missing "public/assets" box on new app without webpack closes \#104 [\#107](https://github.com/gobuffalo/buffalo/pull/107) ([markbates](https://github.com/markbates))
- Using gentronics to generate the templates and the actions [\#106](https://github.com/gobuffalo/buffalo/pull/106) ([apaganobeleno](https://github.com/apaganobeleno))
- Adds an Actions generator to the cmd package [\#103](https://github.com/gobuffalo/buffalo/pull/103) ([apaganobeleno](https://github.com/apaganobeleno))
- Cleaning up some docs [\#102](https://github.com/gobuffalo/buffalo/pull/102) ([CodyOss](https://github.com/CodyOss))
- Avoiding generating log folder when running tests. [\#101](https://github.com/gobuffalo/buffalo/pull/101) ([apaganobeleno](https://github.com/apaganobeleno))
- \[grift\] changing the task to be pointing the buffalo repo [\#100](https://github.com/gobuffalo/buffalo/pull/100) ([apaganobeleno](https://github.com/apaganobeleno))
- making buffalo call grift for the tasks [\#98](https://github.com/gobuffalo/buffalo/pull/98) ([apaganobeleno](https://github.com/apaganobeleno))
- Update build.go [\#93](https://github.com/gobuffalo/buffalo/pull/93) ([arifemre](https://github.com/arifemre))
- Update new.go [\#92](https://github.com/gobuffalo/buffalo/pull/92) ([arifemre](https://github.com/arifemre))
- Typo in readme [\#90](https://github.com/gobuffalo/buffalo/pull/90) ([CodyOss](https://github.com/CodyOss))
- Typos in readme [\#88](https://github.com/gobuffalo/buffalo/pull/88) ([kennygrant](https://github.com/kennygrant))

## [v0.7.0](https://github.com/gobuffalo/buffalo/tree/v0.7.0) (2017-01-04)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.6.0...v0.7.0)

**Closed issues:**

- Error When Generating Goth Code [\#79](https://github.com/gobuffalo/buffalo/issues/79)
- Goth generator does not install required dependencies [\#75](https://github.com/gobuffalo/buffalo/issues/75)
- Export goGet and goInstall from cmd package [\#73](https://github.com/gobuffalo/buffalo/issues/73)

**Merged pull requests:**

- issue-37: initial commit for codeclimate yml generation [\#83](https://github.com/gobuffalo/buffalo/pull/83) ([briandowns](https://github.com/briandowns))
- issue-35: readme generator [\#82](https://github.com/gobuffalo/buffalo/pull/82) ([briandowns](https://github.com/briandowns))
- 0.7.0 [\#81](https://github.com/gobuffalo/buffalo/pull/81) ([markbates](https://github.com/markbates))
- test the goth generator in docker [\#78](https://github.com/gobuffalo/buffalo/pull/78) ([markbates](https://github.com/markbates))
- added the CopyWebpackPlugin to copy files from assets directory [\#77](https://github.com/gobuffalo/buffalo/pull/77) ([markbates](https://github.com/markbates))
- Install required dependencies when using the Goth generator [\#76](https://github.com/gobuffalo/buffalo/pull/76) ([intabulas](https://github.com/intabulas))
- Make GoGet and GoInstall useable from generators [\#74](https://github.com/gobuffalo/buffalo/pull/74) ([intabulas](https://github.com/intabulas))

## [v0.6.0](https://github.com/gobuffalo/buffalo/tree/v0.6.0) (2016-12-29)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.5.1...v0.6.0)

## [v0.5.1](https://github.com/gobuffalo/buffalo/tree/v0.5.1) (2016-12-22)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.5.0...v0.5.1)

**Closed issues:**

- groups on groups should inherit the prefix of the previous group [\#72](https://github.com/gobuffalo/buffalo/issues/72)
- Improve resource generator to insert the resource into actions/app.go [\#43](https://github.com/gobuffalo/buffalo/issues/43)

## [v0.5.0](https://github.com/gobuffalo/buffalo/tree/v0.5.0) (2016-12-21)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.4.7...v0.5.0)

**Closed issues:**

- Log should output host:port when starting up [\#70](https://github.com/gobuffalo/buffalo/issues/70)
- add web pack to the "new" generator [\#18](https://github.com/gobuffalo/buffalo/issues/18)

**Merged pull requests:**

- add web pack to the "new" generator closes \#18 [\#71](https://github.com/gobuffalo/buffalo/pull/71) ([markbates](https://github.com/markbates))

## [v0.4.7](https://github.com/gobuffalo/buffalo/tree/v0.4.7) (2016-12-19)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.4.6...v0.4.7)

**Closed issues:**

- Add a generator for Goth [\#65](https://github.com/gobuffalo/buffalo/issues/65)
- Add a REPL/console [\#63](https://github.com/gobuffalo/buffalo/issues/63)

**Merged pull requests:**

- auto mount a generated resource [\#69](https://github.com/gobuffalo/buffalo/pull/69) ([markbates](https://github.com/markbates))
- add Host as an option to the App [\#68](https://github.com/gobuffalo/buffalo/pull/68) ([markbates](https://github.com/markbates))
- Add a generator for Goth closes \#65 [\#66](https://github.com/gobuffalo/buffalo/pull/66) ([markbates](https://github.com/markbates))
- Add a REPL/console closes \#63 [\#64](https://github.com/gobuffalo/buffalo/pull/64) ([markbates](https://github.com/markbates))

## [v0.4.6](https://github.com/gobuffalo/buffalo/tree/v0.4.6) (2016-12-15)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.4.5...v0.4.6)

**Closed issues:**

- Go Get Errors [\#61](https://github.com/gobuffalo/buffalo/issues/61)
- `buffalo db` should forward to `soda` [\#58](https://github.com/gobuffalo/buffalo/issues/58)
- markdown rendering breaks some raymond parsing [\#55](https://github.com/gobuffalo/buffalo/issues/55)
- add template caching [\#54](https://github.com/gobuffalo/buffalo/issues/54)

**Merged pull requests:**

- Added a resolvers package to help find and resolve files. [\#62](https://github.com/gobuffalo/buffalo/pull/62) ([markbates](https://github.com/markbates))
- `buffalo db` should forward to `soda` closes \#58 [\#60](https://github.com/gobuffalo/buffalo/pull/60) ([markbates](https://github.com/markbates))
- small bug fixes to template caching [\#57](https://github.com/gobuffalo/buffalo/pull/57) ([markbates](https://github.com/markbates))
- add template caching closes \#54 [\#56](https://github.com/gobuffalo/buffalo/pull/56) ([markbates](https://github.com/markbates))

## [v0.4.5](https://github.com/gobuffalo/buffalo/tree/v0.4.5) (2016-12-13)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.4.4...v0.4.5)

**Closed issues:**

- generate SHOULDERS.md when deploying a new release. [\#52](https://github.com/gobuffalo/buffalo/issues/52)
- generate js & css files into sub directory of assets [\#49](https://github.com/gobuffalo/buffalo/issues/49)
- PORT should be settable via an ENV var [\#47](https://github.com/gobuffalo/buffalo/issues/47)
- define ENV at the top of actions/app.go [\#46](https://github.com/gobuffalo/buffalo/issues/46)
- Add NewRelic middleware [\#45](https://github.com/gobuffalo/buffalo/issues/45)
- Procfile should use the project name by default [\#44](https://github.com/gobuffalo/buffalo/issues/44)
- Add a "resource" generator [\#41](https://github.com/gobuffalo/buffalo/issues/41)
- Add "bootstrap" to the "new" generator \(optional\) [\#24](https://github.com/gobuffalo/buffalo/issues/24)
- add a markdown renderer [\#13](https://github.com/gobuffalo/buffalo/issues/13)

**Merged pull requests:**

- Shoulders [\#53](https://github.com/gobuffalo/buffalo/pull/53) ([markbates](https://github.com/markbates))
- add a markdown renderer closes \#13 [\#51](https://github.com/gobuffalo/buffalo/pull/51) ([markbates](https://github.com/markbates))
- generate js & css files into sub directory of assets  [\#50](https://github.com/gobuffalo/buffalo/pull/50) ([markbates](https://github.com/markbates))
- A bunch of fixes [\#48](https://github.com/gobuffalo/buffalo/pull/48) ([markbates](https://github.com/markbates))

## [v0.4.4](https://github.com/gobuffalo/buffalo/tree/v0.4.4) (2016-12-11)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.4.3.1...v0.4.4)

**Closed issues:**

- generated code should pass govet and golint [\#39](https://github.com/gobuffalo/buffalo/issues/39)
- Run gofmt on generated code [\#38](https://github.com/gobuffalo/buffalo/issues/38)

**Merged pull requests:**

- Generators [\#40](https://github.com/gobuffalo/buffalo/pull/40) ([markbates](https://github.com/markbates))

## [v0.4.3.1](https://github.com/gobuffalo/buffalo/tree/v0.4.3.1) (2016-12-11)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.4.3...v0.4.3.1)

## [v0.4.3](https://github.com/gobuffalo/buffalo/tree/v0.4.3) (2016-12-10)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.4.2...v0.4.3)

**Closed issues:**

- Add some functions from the strings package as helpers [\#33](https://github.com/gobuffalo/buffalo/issues/33)
- Add the functions from https://github.com/markbates/inflect as helpers [\#32](https://github.com/gobuffalo/buffalo/issues/32)
- Add support for mapping "Resources" [\#31](https://github.com/gobuffalo/buffalo/issues/31)
- "grift release" should run the "docker CI" first. [\#29](https://github.com/gobuffalo/buffalo/issues/29)

**Merged pull requests:**

- added more helpers. closes \#32 closes \#33 [\#34](https://github.com/gobuffalo/buffalo/pull/34) ([markbates](https://github.com/markbates))
- Working on adding support for a Resource interface [\#30](https://github.com/gobuffalo/buffalo/pull/30) ([markbates](https://github.com/markbates))

## [v0.4.2](https://github.com/gobuffalo/buffalo/tree/v0.4.2) (2016-12-09)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.4.1...v0.4.2)

## [v0.4.1](https://github.com/gobuffalo/buffalo/tree/v0.4.1) (2016-12-09)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.4.1.pre...v0.4.1)

## [v0.4.1.pre](https://github.com/gobuffalo/buffalo/tree/v0.4.1.pre) (2016-12-09)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.4.0...v0.4.1.pre)

## [v0.4.0](https://github.com/gobuffalo/buffalo/tree/v0.4.0) (2016-12-09)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.4.0.pre...v0.4.0)

## [v0.4.0.pre](https://github.com/gobuffalo/buffalo/tree/v0.4.0.pre) (2016-12-09)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/0.4.0...v0.4.0.pre)

## [0.4.0](https://github.com/gobuffalo/buffalo/tree/0.4.0) (2016-12-09)
**Closed issues:**

- replace httprouter with gorilla mux [\#27](https://github.com/gobuffalo/buffalo/issues/27)
- Add "jQuery" to the "new" generator \(optional\) [\#23](https://github.com/gobuffalo/buffalo/issues/23)
- Generate a Procfile in the `new` generator [\#21](https://github.com/gobuffalo/buffalo/issues/21)
- Add a pop transaction middleware to the "new" generator [\#15](https://github.com/gobuffalo/buffalo/issues/15)
- add a cmd to run the app in "dev" w/ refresh [\#12](https://github.com/gobuffalo/buffalo/issues/12)
- Add refresh to the generator [\#11](https://github.com/gobuffalo/buffalo/issues/11)
- Add pop/soda to the generator [\#10](https://github.com/gobuffalo/buffalo/issues/10)
- Add grift to the generator [\#9](https://github.com/gobuffalo/buffalo/issues/9)
- add a wrapHandlerFunc fund [\#8](https://github.com/gobuffalo/buffalo/issues/8)
- add a wrapHandler function [\#7](https://github.com/gobuffalo/buffalo/issues/7)
- Add template caching [\#6](https://github.com/gobuffalo/buffalo/issues/6)
- Serve static files [\#5](https://github.com/gobuffalo/buffalo/issues/5)
- Add Websocket support [\#4](https://github.com/gobuffalo/buffalo/issues/4)
- Need `bind` function [\#3](https://github.com/gobuffalo/buffalo/issues/3)
- Need README [\#2](https://github.com/gobuffalo/buffalo/issues/2)
- Need GoDoc [\#1](https://github.com/gobuffalo/buffalo/issues/1)

**Merged pull requests:**

- replace httprouter with gorilla mux closes closes \#27 [\#28](https://github.com/gobuffalo/buffalo/pull/28) ([markbates](https://github.com/markbates))
- added some helpers and event source support [\#26](https://github.com/gobuffalo/buffalo/pull/26) ([markbates](https://github.com/markbates))
- switched over to using gentronics for generating templates [\#25](https://github.com/gobuffalo/buffalo/pull/25) ([markbates](https://github.com/markbates))



\* *This Change Log was automatically generated by [github_changelog_generator](https://github.com/skywinder/Github-Changelog-Generator)*