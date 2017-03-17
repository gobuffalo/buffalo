# Change Log

## [v0.7.4](https://github.com/gobuffalo/buffalo/tree/v0.7.4) (2017-03-03)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.7.3...v0.7.4)

**Implemented enhancements:**

- Add yarn as alternative --with-yarn [\#256](https://github.com/gobuffalo/buffalo/pull/256) ([fooflare](https://github.com/fooflare))

**Fixed bugs:**

- Buffalo needs to be updated to support Webpack 2.2.x [\#195](https://github.com/gobuffalo/buffalo/issues/195)

**Closed issues:**

- buffalo command throwing errors [\#266](https://github.com/gobuffalo/buffalo/issues/266)
- v0.7.3 webpack2.0 release breaks default configuration [\#262](https://github.com/gobuffalo/buffalo/issues/262)
- models starting with a b produce errors [\#261](https://github.com/gobuffalo/buffalo/issues/261)
- App.Group incorrectly builds paths on Windows machines [\#258](https://github.com/gobuffalo/buffalo/issues/258)
- Can't register custom error handler [\#255](https://github.com/gobuffalo/buffalo/issues/255)
- install fails when sqllite fails to build. [\#253](https://github.com/gobuffalo/buffalo/issues/253)
- Add Support for MongoDB via mgo package [\#252](https://github.com/gobuffalo/buffalo/issues/252)
- Documentation Needs to State Requirement for Go \>= 1.7 [\#251](https://github.com/gobuffalo/buffalo/issues/251)

**Merged pull requests:**

- Inching towards being able to use different template engines in Buffalo [\#265](https://github.com/gobuffalo/buffalo/pull/265) ([markbates](https://github.com/markbates))
- support binding of html multipart post requests [\#264](https://github.com/gobuffalo/buffalo/pull/264) ([tsauter](https://github.com/tsauter))
- Webpack v2.2.1 upgrade [\#263](https://github.com/gobuffalo/buffalo/pull/263) ([fooflare](https://github.com/fooflare))
- fixed build so it also builds non-db apps [\#260](https://github.com/gobuffalo/buffalo/pull/260) ([markbates](https://github.com/markbates))
- updated router Group function to acount for Windows path cruft [\#259](https://github.com/gobuffalo/buffalo/pull/259) ([schigh](https://github.com/schigh))
- Copy ErrorHandlers to Group. [\#257](https://github.com/gobuffalo/buffalo/pull/257) ([drlogout](https://github.com/drlogout))
- Add Go version requirement to README.md [\#254](https://github.com/gobuffalo/buffalo/pull/254) ([gillchristian](https://github.com/gillchristian))

## [v0.7.3](https://github.com/gobuffalo/buffalo/tree/v0.7.3) (2017-02-15)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.7.2...v0.7.3)

**Implemented enhancements:**

- Add a Redirect function to the Router [\#245](https://github.com/gobuffalo/buffalo/issues/245)
- Add a `Clear` function to Session [\#230](https://github.com/gobuffalo/buffalo/issues/230)
- Run tasks from the built binary [\#224](https://github.com/gobuffalo/buffalo/issues/224)
-  create a new buffalo app in the current directory [\#206](https://github.com/gobuffalo/buffalo/issues/206)

**Closed issues:**

- Best way to wrap or dispatch to http.Handler [\#241](https://github.com/gobuffalo/buffalo/issues/241)
- Allow for new binders to be registered with Buffalo [\#222](https://github.com/gobuffalo/buffalo/issues/222)
- nulls package types can't be decoded with gorilla [\#221](https://github.com/gobuffalo/buffalo/issues/221)

**Merged pull requests:**

- Run tasks from the built binary closes \#224 [\#248](https://github.com/gobuffalo/buffalo/pull/248) ([markbates](https://github.com/markbates))
- Use envy instead of defaults for new apps [\#247](https://github.com/gobuffalo/buffalo/pull/247) ([markbates](https://github.com/markbates))
- Add a Redirect function to the Router closes \#245 [\#246](https://github.com/gobuffalo/buffalo/pull/246) ([markbates](https://github.com/markbates))
- Add an empty Helpers map to actions/render.go [\#244](https://github.com/gobuffalo/buffalo/pull/244) ([markbates](https://github.com/markbates))
- Content type ranging - extra test [\#238](https://github.com/gobuffalo/buffalo/pull/238) ([philipithomas](https://github.com/philipithomas))
- fixed the generation of the refresh file if it didn't exist [\#237](https://github.com/gobuffalo/buffalo/pull/237) ([markbates](https://github.com/markbates))
- buffalo db should now print out `buffalo db` and not `buffalo soda`. [\#236](https://github.com/gobuffalo/buffalo/pull/236) ([markbates](https://github.com/markbates))
- removed the no longer existing docs for LogDir and added some for [\#235](https://github.com/gobuffalo/buffalo/pull/235) ([markbates](https://github.com/markbates))
- content types need to be ranged over in case of ones with a ';' in them [\#234](https://github.com/gobuffalo/buffalo/pull/234) ([markbates](https://github.com/markbates))
- create a new buffalo app in the current directory closes \#206 [\#233](https://github.com/gobuffalo/buffalo/pull/233) ([markbates](https://github.com/markbates))
- put the current\_path in the context closes \#207 [\#232](https://github.com/gobuffalo/buffalo/pull/232) ([markbates](https://github.com/markbates))
- Add a `Clear` function to Session closes \#230 [\#231](https://github.com/gobuffalo/buffalo/pull/231) ([markbates](https://github.com/markbates))
- Update usage of validate in html-crud example [\#228](https://github.com/gobuffalo/buffalo/pull/228) ([srt32](https://github.com/srt32))
- removed the multilogger since it wasn't providing any real benefit [\#227](https://github.com/gobuffalo/buffalo/pull/227) ([markbates](https://github.com/markbates))
- removed the used MethodOverride var and a duplicate check for setting the MethodOverride [\#226](https://github.com/gobuffalo/buffalo/pull/226) ([markbates](https://github.com/markbates))
- Custom binders [\#223](https://github.com/gobuffalo/buffalo/pull/223) ([markbates](https://github.com/markbates))

## [v0.7.2](https://github.com/gobuffalo/buffalo/tree/v0.7.2) (2017-02-03)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.7.1...v0.7.2)

**Implemented enhancements:**

- Buffalo new command unsopported multiple GOPATH [\#203](https://github.com/gobuffalo/buffalo/issues/203)
- Resource generator shouldn't generate pointers [\#198](https://github.com/gobuffalo/buffalo/issues/198)
- Resource generator should generate a better variable name [\#197](https://github.com/gobuffalo/buffalo/issues/197)
- Make sure a new Buffalo app doesn't have Context\#Get deprecation warnings [\#194](https://github.com/gobuffalo/buffalo/issues/194)
- Add file:line info to Context\#Get deprecation warnings [\#193](https://github.com/gobuffalo/buffalo/issues/193)
- Logrus middleware for Buffalo [\#188](https://github.com/gobuffalo/buffalo/issues/188)
- Allow build command to split assets bundle and binary app [\#170](https://github.com/gobuffalo/buffalo/issues/170)
- Buffalo fails to start/build/test on Windows [\#166](https://github.com/gobuffalo/buffalo/issues/166)
- Buffalo new "welcome" output [\#164](https://github.com/gobuffalo/buffalo/issues/164)
- Flash [\#139](https://github.com/gobuffalo/buffalo/issues/139)
- Need a way to easily add "routes" when generating files. [\#105](https://github.com/gobuffalo/buffalo/issues/105)
- Add a form builder helper  [\#96](https://github.com/gobuffalo/buffalo/issues/96)
- Buffalo Docker Image [\#89](https://github.com/gobuffalo/buffalo/issues/89)
- generate a default .travis.yml file for new projects [\#36](https://github.com/gobuffalo/buffalo/issues/36)
- add a `grift test` task [\#20](https://github.com/gobuffalo/buffalo/issues/20)
- Add a default 500 page [\#14](https://github.com/gobuffalo/buffalo/issues/14)
- Updated Unknown Error message [\#162](https://github.com/gobuffalo/buffalo/pull/162) ([bscott](https://github.com/bscott))
- Skip webpack install if already in PATH [\#158](https://github.com/gobuffalo/buffalo/pull/158) ([mdhender](https://github.com/mdhender))

**Fixed bugs:**

- Buffalo fails to start/build/test on Windows [\#166](https://github.com/gobuffalo/buffalo/issues/166)
- running multiple `buffalo` commands causes a "timeout" [\#67](https://github.com/gobuffalo/buffalo/issues/67)

**Closed issues:**

- Remove all licensing from the top of `cmd` files [\#208](https://github.com/gobuffalo/buffalo/issues/208)
- Templates when Rendering from Resource [\#204](https://github.com/gobuffalo/buffalo/issues/204)
- Silent failure of unknown db type [\#183](https://github.com/gobuffalo/buffalo/issues/183)
- Action generator should append new tests instead of clearing test file [\#174](https://github.com/gobuffalo/buffalo/issues/174)
- Generating invalid .codeclimate.yml [\#173](https://github.com/gobuffalo/buffalo/issues/173)
- Binary name should have ".exe" implicitly appended to it on Windows [\#169](https://github.com/gobuffalo/buffalo/issues/169)
- There was a problem starting the dev server [\#156](https://github.com/gobuffalo/buffalo/issues/156)
- buffalo.Context should implement the context.Context interface [\#153](https://github.com/gobuffalo/buffalo/issues/153)
- Cross-compiling fails from 64-bit linux to arm linux [\#142](https://github.com/gobuffalo/buffalo/issues/142)
- Prevent new app creation from outside the Go workspace [\#140](https://github.com/gobuffalo/buffalo/issues/140)
- Problem starting "buffalo dev" server [\#135](https://github.com/gobuffalo/buffalo/issues/135)
- Installation on macOS fails at go-sqlite3 step [\#128](https://github.com/gobuffalo/buffalo/issues/128)
- Tests only running in latest stable Go version [\#123](https://github.com/gobuffalo/buffalo/issues/123)
- Add a test helper equivalent of the PopTransaction middleware [\#22](https://github.com/gobuffalo/buffalo/issues/22)

**Merged pull requests:**

- make the Flash\#Persist function private [\#219](https://github.com/gobuffalo/buffalo/pull/219) ([markbates](https://github.com/markbates))
- Buffalo new "welcome" output closes \#164 [\#218](https://github.com/gobuffalo/buffalo/pull/218) ([markbates](https://github.com/markbates))
- added generator tests for goth and changed where it adds the routes [\#215](https://github.com/gobuffalo/buffalo/pull/215) ([markbates](https://github.com/markbates))
- Adds Test cases for the travis/none ci-provider generation [\#212](https://github.com/gobuffalo/buffalo/pull/212) ([apaganobeleno](https://github.com/apaganobeleno))
- Remove all licensing from the top of `cmd` files closes \#208 [\#209](https://github.com/gobuffalo/buffalo/pull/209) ([markbates](https://github.com/markbates))
- Fixed issue \#203: multiple GOPATH are not supported by buffalo new. [\#205](https://github.com/gobuffalo/buffalo/pull/205) ([stanislas-m](https://github.com/stanislas-m))
- \[\#193\] adding line number to the Context\#Get warning [\#202](https://github.com/gobuffalo/buffalo/pull/202) ([apaganobeleno](https://github.com/apaganobeleno))
- \[\#194\] Avoid Context\#Get warnings on the newly created app. [\#201](https://github.com/gobuffalo/buffalo/pull/201) ([apaganobeleno](https://github.com/apaganobeleno))
- \[\#198\] moving resource generator to avoid pointers [\#200](https://github.com/gobuffalo/buffalo/pull/200) ([apaganobeleno](https://github.com/apaganobeleno))
- \[\#197\] generates better code inside actions [\#199](https://github.com/gobuffalo/buffalo/pull/199) ([apaganobeleno](https://github.com/apaganobeleno))
- need to peg webpack to 1.14.0 because 2.2.x doesn't work with the generated config [\#196](https://github.com/gobuffalo/buffalo/pull/196) ([markbates](https://github.com/markbates))
- Fixed actions generator: imports were missing. [\#192](https://github.com/gobuffalo/buffalo/pull/192) ([stanislas-m](https://github.com/stanislas-m))
- remove unused code [\#191](https://github.com/gobuffalo/buffalo/pull/191) ([philipithomas](https://github.com/philipithomas))
- convert Version to constant [\#190](https://github.com/gobuffalo/buffalo/pull/190) ([philipithomas](https://github.com/philipithomas))
- Add Golint and fix all issues [\#189](https://github.com/gobuffalo/buffalo/pull/189) ([philipithomas](https://github.com/philipithomas))
- better error printing. [\#186](https://github.com/gobuffalo/buffalo/pull/186) ([markbates](https://github.com/markbates))
- fixes \#183 [\#185](https://github.com/gobuffalo/buffalo/pull/185) ([amedeiros](https://github.com/amedeiros))
- \[feature\] adds .travis.yml when generating a new app [\#184](https://github.com/gobuffalo/buffalo/pull/184) ([apaganobeleno](https://github.com/apaganobeleno))
- Allow buffalo build to extract assets and put them into a zip file [\#180](https://github.com/gobuffalo/buffalo/pull/180) ([stanislas-m](https://github.com/stanislas-m))
- set a default session name based on the project when it is created [\#179](https://github.com/gobuffalo/buffalo/pull/179) ([markbates](https://github.com/markbates))
- fixed an issue where flash messages where not clearing properly [\#178](https://github.com/gobuffalo/buffalo/pull/178) ([markbates](https://github.com/markbates))
- Fixes \#173 [\#177](https://github.com/gobuffalo/buffalo/pull/177) ([amedeiros](https://github.com/amedeiros))
- Skipping test generation if test exists on action generation [\#176](https://github.com/gobuffalo/buffalo/pull/176) ([apaganobeleno](https://github.com/apaganobeleno))
- Fixed action generator which erased previously defined tests [\#175](https://github.com/gobuffalo/buffalo/pull/175) ([stanislas-m](https://github.com/stanislas-m))
- Fix rendering if alternative layout is used  closes \#167 [\#171](https://github.com/gobuffalo/buffalo/pull/171) ([markbates](https://github.com/markbates))
- fix webpack to run locally on windows [\#168](https://github.com/gobuffalo/buffalo/pull/168) ([markbates](https://github.com/markbates))
- Added jetbrains IDE workspace directory for such IDE's as Gogland [\#163](https://github.com/gobuffalo/buffalo/pull/163) ([bscott](https://github.com/bscott))
- Creating a new app with webpack requires admin privileges closes \#157 [\#161](https://github.com/gobuffalo/buffalo/pull/161) ([markbates](https://github.com/markbates))
- updated contributers list [\#155](https://github.com/gobuffalo/buffalo/pull/155) ([markbates](https://github.com/markbates))
- buffalo.Context should implement the context.Context interface close … [\#154](https://github.com/gobuffalo/buffalo/pull/154) ([markbates](https://github.com/markbates))
- changed a few pointers in the render package to not pointers [\#152](https://github.com/gobuffalo/buffalo/pull/152) ([markbates](https://github.com/markbates))
- removed a few debug statements [\#151](https://github.com/gobuffalo/buffalo/pull/151) ([markbates](https://github.com/markbates))
- Fix spelling mistake in PopTransaction documentation [\#147](https://github.com/gobuffalo/buffalo/pull/147) ([DanielHeckrath](https://github.com/DanielHeckrath))
- fix browser typo [\#145](https://github.com/gobuffalo/buffalo/pull/145) ([dankleiman](https://github.com/dankleiman))
- Minor spelling/grammar fixes [\#144](https://github.com/gobuffalo/buffalo/pull/144) ([mrxinu](https://github.com/mrxinu))
- Implementing the `flash` helper [\#143](https://github.com/gobuffalo/buffalo/pull/143) ([apaganobeleno](https://github.com/apaganobeleno))
- Prevent new app creation from outside the Go workspace [\#141](https://github.com/gobuffalo/buffalo/pull/141) ([markbates](https://github.com/markbates))
- UnWrap HttpErrors in pop middleware and return them. [\#138](https://github.com/gobuffalo/buffalo/pull/138) ([lumost](https://github.com/lumost))
- preventing wrapping errors with point in error handler stack [\#137](https://github.com/gobuffalo/buffalo/pull/137) ([lumost](https://github.com/lumost))
- import buffalo when generating a new action [\#136](https://github.com/gobuffalo/buffalo/pull/136) ([lumost](https://github.com/lumost))
- First Attempt at \#105  [\#112](https://github.com/gobuffalo/buffalo/pull/112) ([apaganobeleno](https://github.com/apaganobeleno))

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
- Don't write test.log files when running tests [\#17](https://github.com/gobuffalo/buffalo/issues/17)
- Add an "actions" generator [\#16](https://github.com/gobuffalo/buffalo/issues/16)

**Merged pull requests:**

- Add badge for Go Report Card to README [\#132](https://github.com/gobuffalo/buffalo/pull/132) ([stuartellis](https://github.com/stuartellis))
- Makes our tests run on Go 1.7 and 1.8 [\#131](https://github.com/gobuffalo/buffalo/pull/131) ([apaganobeleno](https://github.com/apaganobeleno))
- build\_path does not work for Windows closes \#124 [\#130](https://github.com/gobuffalo/buffalo/pull/130) ([markbates](https://github.com/markbates))
- Edit some typo [\#129](https://github.com/gobuffalo/buffalo/pull/129) ([janczer](https://github.com/janczer))
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
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/0.4.0...v0.4.1)

## [0.4.0](https://github.com/gobuffalo/buffalo/tree/0.4.0) (2016-12-09)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.4.0...0.4.0)

## [v0.4.0](https://github.com/gobuffalo/buffalo/tree/v0.4.0) (2016-12-09)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.4.0.pre...v0.4.0)

## [v0.4.0.pre](https://github.com/gobuffalo/buffalo/tree/v0.4.0.pre) (2016-12-09)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.4.1.pre...v0.4.0.pre)

## [v0.4.1.pre](https://github.com/gobuffalo/buffalo/tree/v0.4.1.pre) (2016-12-09)
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