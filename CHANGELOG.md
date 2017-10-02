# Change Log

## [v0.9.5](https://github.com/gobuffalo/buffalo/tree/v0.9.5) (2017-09-30)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.9.4...v0.9.5)

**Fixed bugs:**

- Templates and migrations are bundled in binary with the `-e` flag [\#567](https://github.com/gobuffalo/buffalo/issues/567)

**Closed issues:**

- upgrade build causes could not find folder for box:\(projectpath\)/actions/templates [\#660](https://github.com/gobuffalo/buffalo/issues/660)
- Mounting a resource on a resource produces bad routes [\#656](https://github.com/gobuffalo/buffalo/issues/656)
- Using Goth for login/password auth? [\#652](https://github.com/gobuffalo/buffalo/issues/652)
- CSS for error pages is broken in v0.9.5 [\#650](https://github.com/gobuffalo/buffalo/issues/650)
- Import error generating new app as API only [\#648](https://github.com/gobuffalo/buffalo/issues/648)
- Any simple way to use custom layout? [\#640](https://github.com/gobuffalo/buffalo/issues/640)
- Need an invite to slack channel? [\#639](https://github.com/gobuffalo/buffalo/issues/639)
- building linux binaries from the Mac with sqlite [\#638](https://github.com/gobuffalo/buffalo/issues/638)
- mount at group router not work [\#637](https://github.com/gobuffalo/buffalo/issues/637)
- Duplicate Session Set-Cookie Header [\#631](https://github.com/gobuffalo/buffalo/issues/631)
- Error: you need to be inside your buffalo project path to run this command [\#629](https://github.com/gobuffalo/buffalo/issues/629)
- A few questions [\#628](https://github.com/gobuffalo/buffalo/issues/628)
- deployed binary cannot pass in argument to task [\#625](https://github.com/gobuffalo/buffalo/issues/625)
- Build on Windows produces incorrect extension [\#623](https://github.com/gobuffalo/buffalo/issues/623)
- `buffalo dev` code auto rebuild doesn't work with NFS [\#620](https://github.com/gobuffalo/buffalo/issues/620)

**Merged pull requests:**

- Fix links in SHOULDERS.md [\#669](https://github.com/gobuffalo/buffalo/pull/669) ([jorinvo](https://github.com/jorinvo))
- PopTransaction middleware should roll back if there is a non success status code [\#668](https://github.com/gobuffalo/buffalo/pull/668) ([markbates](https://github.com/markbates))
- Fixes default error pages showing with broken style. [\#667](https://github.com/gobuffalo/buffalo/pull/667) ([paganotoni](https://github.com/paganotoni))
- added support for compressing assets in binaries [\#666](https://github.com/gobuffalo/buffalo/pull/666) ([markbates](https://github.com/markbates))
- cleaned up the generated resource html a bit [\#665](https://github.com/gobuffalo/buffalo/pull/665) ([markbates](https://github.com/markbates))
- Added go-i18n to the giants [\#662](https://github.com/gobuffalo/buffalo/pull/662) ([fdonzello](https://github.com/fdonzello))
- get the content type using the httpx package [\#661](https://github.com/gobuffalo/buffalo/pull/661) ([markbates](https://github.com/markbates))
- Mounting a resource on a resource produces bad routes fixes \#656 [\#657](https://github.com/gobuffalo/buffalo/pull/657) ([markbates](https://github.com/markbates))
- fix: remove mixed tabs and spaces from .travis.yml template [\#655](https://github.com/gobuffalo/buffalo/pull/655) ([DanielHeckrath](https://github.com/DanielHeckrath))
- Resource error changes [\#654](https://github.com/gobuffalo/buffalo/pull/654) ([markbates](https://github.com/markbates))
- don't make an assetsBox if in api mode [\#653](https://github.com/gobuffalo/buffalo/pull/653) ([markbates](https://github.com/markbates))
- make the app available to the context [\#651](https://github.com/gobuffalo/buffalo/pull/651) ([markbates](https://github.com/markbates))
- Fix API mode import errors [\#649](https://github.com/gobuffalo/buffalo/pull/649) ([tsmith1024](https://github.com/tsmith1024))
- Templates and migrations are bundled in binary with the `-e` flag fix… [\#646](https://github.com/gobuffalo/buffalo/pull/646) ([markbates](https://github.com/markbates))
- changed where the assetsBox is initialized [\#645](https://github.com/gobuffalo/buffalo/pull/645) ([markbates](https://github.com/markbates))
- sessions lookups should happen up the app parent chain [\#644](https://github.com/gobuffalo/buffalo/pull/644) ([markbates](https://github.com/markbates))
- mount at group router not work fixes \#637 [\#641](https://github.com/gobuffalo/buffalo/pull/641) ([markbates](https://github.com/markbates))
- deployed binary cannot pass in argument to task fixes \#625 [\#636](https://github.com/gobuffalo/buffalo/pull/636) ([markbates](https://github.com/markbates))
- Feature assets fingerprint [\#635](https://github.com/gobuffalo/buffalo/pull/635) ([paganotoni](https://github.com/paganotoni))
- do a better job of capturing signals for buffalo dev [\#634](https://github.com/gobuffalo/buffalo/pull/634) ([markbates](https://github.com/markbates))
- middleware.SessionSaver is now built-in. [\#633](https://github.com/gobuffalo/buffalo/pull/633) ([markbates](https://github.com/markbates))
- added deprecation notices around middleware.CSRF as well as prevent the CSRF middleware from running in test mode [\#632](https://github.com/gobuffalo/buffalo/pull/632) ([markbates](https://github.com/markbates))
- gracefully capture and handle signal's during buffalo build [\#630](https://github.com/gobuffalo/buffalo/pull/630) ([markbates](https://github.com/markbates))
- deprecated "Automatic" in favor of just "New". [\#627](https://github.com/gobuffalo/buffalo/pull/627) ([markbates](https://github.com/markbates))
- cache yarn packages in the base docker image. also made the list of [\#626](https://github.com/gobuffalo/buffalo/pull/626) ([markbates](https://github.com/markbates))
- Check GOOS from runtime rather than env var [\#624](https://github.com/gobuffalo/buffalo/pull/624) ([hut8](https://github.com/hut8))

## [v0.9.4](https://github.com/gobuffalo/buffalo/tree/v0.9.4) (2017-09-05)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.9.3...v0.9.4)

**Fixed bugs:**

- Problems with reserved SQL identifier [\#577](https://github.com/gobuffalo/buffalo/issues/577)

**Closed issues:**

- Creating a new project fails at goimports with Buffalo 0.9.3 and Go 1.9 [\#612](https://github.com/gobuffalo/buffalo/issues/612)
- Ability to set up pre-app handlers [\#609](https://github.com/gobuffalo/buffalo/issues/609)
- running db create and db seed and db migrate on server without go [\#607](https://github.com/gobuffalo/buffalo/issues/607)
- Provide way to mount other application [\#602](https://github.com/gobuffalo/buffalo/issues/602)
- Move the Goth generator to a plugin and out of "core" [\#600](https://github.com/gobuffalo/buffalo/issues/600)
- `buffalo db create -e development` asking for password [\#596](https://github.com/gobuffalo/buffalo/issues/596)
- buffalo new error message incorrect [\#593](https://github.com/gobuffalo/buffalo/issues/593)
- Did my setup fail or are instructions incomplete?  [\#590](https://github.com/gobuffalo/buffalo/issues/590)
- \(&DefaultContext{}\).Redirect\(\) has surprising format-string behavior [\#586](https://github.com/gobuffalo/buffalo/issues/586)
- Asset Fingerprinting [\#578](https://github.com/gobuffalo/buffalo/issues/578)
- Please add example for login forms [\#149](https://github.com/gobuffalo/buffalo/issues/149)

**Merged pull requests:**

- fix for gorilla context memory leak [\#619](https://github.com/gobuffalo/buffalo/pull/619) ([nutbunnies](https://github.com/nutbunnies))
- reverts assets-figerprint [\#618](https://github.com/gobuffalo/buffalo/pull/618) ([paganotoni](https://github.com/paganotoni))
- v0.9.4 [\#617](https://github.com/gobuffalo/buffalo/pull/617) ([markbates](https://github.com/markbates))
- use bootstrap-sass not regular bootstrap [\#616](https://github.com/gobuffalo/buffalo/pull/616) ([markbates](https://github.com/markbates))
- Fix 'formular' typo [\#615](https://github.com/gobuffalo/buffalo/pull/615) ([slomek](https://github.com/slomek))
- Enable plugin subcommands [\#614](https://github.com/gobuffalo/buffalo/pull/614) ([ryanfaerman](https://github.com/ryanfaerman))
- Ability to set up pre-app handlers fixes \#609 [\#610](https://github.com/gobuffalo/buffalo/pull/610) ([markbates](https://github.com/markbates))
- Fixes order for yarn command on setup.go [\#608](https://github.com/gobuffalo/buffalo/pull/608) ([mayra-cabrera](https://github.com/mayra-cabrera))
- use a null session for APIs [\#605](https://github.com/gobuffalo/buffalo/pull/605) ([markbates](https://github.com/markbates))
- Provide way to mount other application fixes \#602 [\#604](https://github.com/gobuffalo/buffalo/pull/604) ([markbates](https://github.com/markbates))
- improved the travis yaml [\#603](https://github.com/gobuffalo/buffalo/pull/603) ([markbates](https://github.com/markbates))
- Move the Goth generator to a plugin and out of "core" fixes \#600 [\#601](https://github.com/gobuffalo/buffalo/pull/601) ([markbates](https://github.com/markbates))
- added a middleware for forcing ssl redirect [\#599](https://github.com/gobuffalo/buffalo/pull/599) ([markbates](https://github.com/markbates))
- handle trailing slash properly [\#598](https://github.com/gobuffalo/buffalo/pull/598) ([markbates](https://github.com/markbates))
- Wrong sqlite3 path: buffalo db create puts the database file into wrong directory fixes \#436 [\#597](https://github.com/gobuffalo/buffalo/pull/597) ([markbates](https://github.com/markbates))
- buffalo new error message incorrect fixes \#593 [\#595](https://github.com/gobuffalo/buffalo/pull/595) ([markbates](https://github.com/markbates))
- do a better job of making sure signals are handled better [\#594](https://github.com/gobuffalo/buffalo/pull/594) ([markbates](https://github.com/markbates))
- \(&DefaultContext{}\).Redirect\(\) has surprising format-string behavior … [\#592](https://github.com/gobuffalo/buffalo/pull/592) ([markbates](https://github.com/markbates))
- stop using actions\_test in favor of actions [\#591](https://github.com/gobuffalo/buffalo/pull/591) ([markbates](https://github.com/markbates))
- When fetching packages with go get, remove redundant go install [\#588](https://github.com/gobuffalo/buffalo/pull/588) ([bradleyfalzon](https://github.com/bradleyfalzon))
- Assets fingerprint [\#583](https://github.com/gobuffalo/buffalo/pull/583) ([paganotoni](https://github.com/paganotoni))

## [v0.9.3](https://github.com/gobuffalo/buffalo/tree/v0.9.3) (2017-08-14)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.9.2...v0.9.3)

**Merged pull requests:**

- 0.9.3 [\#587](https://github.com/gobuffalo/buffalo/pull/587) ([markbates](https://github.com/markbates))
- fixed missing error if the user deletes the local logo.svg [\#585](https://github.com/gobuffalo/buffalo/pull/585) ([markbates](https://github.com/markbates))
- set up git for new applications [\#584](https://github.com/gobuffalo/buffalo/pull/584) ([markbates](https://github.com/markbates))

## [v0.9.2](https://github.com/gobuffalo/buffalo/tree/v0.9.2) (2017-08-12)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.9.1.2...v0.9.2)

**Implemented enhancements:**

- Buffalo can't bind on fixed IP [\#518](https://github.com/gobuffalo/buffalo/issues/518)
- colored logs [\#360](https://github.com/gobuffalo/buffalo/issues/360)
- New Generator - SystemD Unit [\#356](https://github.com/gobuffalo/buffalo/issues/356)
- Ship webpack with testing set up [\#97](https://github.com/gobuffalo/buffalo/issues/97)
- Add support for mail [\#94](https://github.com/gobuffalo/buffalo/issues/94)

**Fixed bugs:**

- Last logrus version breaks Buffalo logger\_formatter [\#568](https://github.com/gobuffalo/buffalo/issues/568)

**Closed issues:**

- buffalo new: dep ensure fails [\#580](https://github.com/gobuffalo/buffalo/issues/580)
- Verbose log of workers [\#579](https://github.com/gobuffalo/buffalo/issues/579)
- Some problem on paginator helper for mostly empty or large number of objects. [\#573](https://github.com/gobuffalo/buffalo/issues/573)
- Linux: Buffalo tries to install yarn globally without checking for root rights [\#570](https://github.com/gobuffalo/buffalo/issues/570)
- A "catch-all" route causes a panic [\#561](https://github.com/gobuffalo/buffalo/issues/561)
- dep isn't being automatically installed in 0.9.2 [\#551](https://github.com/gobuffalo/buffalo/issues/551)
- skip dep automatically if it's not installed on `new` [\#550](https://github.com/gobuffalo/buffalo/issues/550)
- New Command generates invalid app settings [\#540](https://github.com/gobuffalo/buffalo/issues/540)
- webpack error in Buffalo 0.9.1.2 [\#532](https://github.com/gobuffalo/buffalo/issues/532)
- \[Question\] db.Transaction and controller unittest [\#528](https://github.com/gobuffalo/buffalo/issues/528)
- Mac issue [\#524](https://github.com/gobuffalo/buffalo/issues/524)
- Buffalo dev runtime error upon trying to access local host page [\#521](https://github.com/gobuffalo/buffalo/issues/521)
- Demo/Example Buffalo Application [\#520](https://github.com/gobuffalo/buffalo/issues/520)
- \[Question\] showing sql queries when testing or log to test.log [\#519](https://github.com/gobuffalo/buffalo/issues/519)
- Assets missing on new v0.9.1.2 app without npm [\#517](https://github.com/gobuffalo/buffalo/issues/517)
- Environment variables in .buffalo.dev.yml? [\#516](https://github.com/gobuffalo/buffalo/issues/516)
- How to start app at 0.0.0.0:3000? [\#515](https://github.com/gobuffalo/buffalo/issues/515)
- Weird bug regarding $GOPATH bin directory [\#513](https://github.com/gobuffalo/buffalo/issues/513)
- Provide an error if -\(dash\) in a helper name. [\#512](https://github.com/gobuffalo/buffalo/issues/512)
- Running buffalo commands outside of the directory should error out [\#508](https://github.com/gobuffalo/buffalo/issues/508)
- Add explicit version \(or range\) to the Node packages [\#503](https://github.com/gobuffalo/buffalo/issues/503)
- Getting started:  running buffalo new coke with 0.9.0 and go 1.8.3 fails on windows [\#500](https://github.com/gobuffalo/buffalo/issues/500)
- Unsuccessful: buffalo new \<name\> --skip-pop [\#497](https://github.com/gobuffalo/buffalo/issues/497)
- Create a new context.Cookie wrapper for gorilla Cookiestore [\#391](https://github.com/gobuffalo/buffalo/issues/391)
- Investigate GopherJS integration [\#358](https://github.com/gobuffalo/buffalo/issues/358)
- New Generator - Caddyfile [\#355](https://github.com/gobuffalo/buffalo/issues/355)
- Mail Registration/Authentication. [\#217](https://github.com/gobuffalo/buffalo/issues/217)

**Merged pull requests:**

- Skip the windows directory when looking for plugins [\#582](https://github.com/gobuffalo/buffalo/pull/582) ([ryanfaerman](https://github.com/ryanfaerman))
- changed log level for perform method of simple worker. [\#581](https://github.com/gobuffalo/buffalo/pull/581) ([sio4](https://github.com/sio4))
- Fix GitLab CI config generator [\#576](https://github.com/gobuffalo/buffalo/pull/576) ([stanislas-m](https://github.com/stanislas-m))
- Add Prefix to RootPath \(BaseUrl\) [\#572](https://github.com/gobuffalo/buffalo/pull/572) ([markbates](https://github.com/markbates))
- Fix \#568: logger\_formatter was broken after logrus update [\#571](https://github.com/gobuffalo/buffalo/pull/571) ([stanislas-m](https://github.com/stanislas-m))
- Add Context.Cookies [\#569](https://github.com/gobuffalo/buffalo/pull/569) ([bigblind](https://github.com/bigblind))
- cleaned up logging by removing the big giant empty string. pop also now [\#564](https://github.com/gobuffalo/buffalo/pull/564) ([markbates](https://github.com/markbates))
- Added support for aliasing a route [\#563](https://github.com/gobuffalo/buffalo/pull/563) ([markbates](https://github.com/markbates))
- Fixing a panic that happens when implementing a Catch-All route. [\#562](https://github.com/gobuffalo/buffalo/pull/562) ([paganotoni](https://github.com/paganotoni))
- Routes escape extra args [\#560](https://github.com/gobuffalo/buffalo/pull/560) ([paganotoni](https://github.com/paganotoni))
- Removes extra ? when the route helper doesn't have extra args [\#559](https://github.com/gobuffalo/buffalo/pull/559) ([paganotoni](https://github.com/paganotoni))
- override the default `go` path with ENV [\#558](https://github.com/gobuffalo/buffalo/pull/558) ([markbates](https://github.com/markbates))
- yarn is now the default instead of npm [\#555](https://github.com/gobuffalo/buffalo/pull/555) ([markbates](https://github.com/markbates))
- make sure yarn uses --no-progress [\#554](https://github.com/gobuffalo/buffalo/pull/554) ([markbates](https://github.com/markbates))
- reversed the policy on dep. it is now off by default. [\#552](https://github.com/gobuffalo/buffalo/pull/552) ([markbates](https://github.com/markbates))
- WIP 0.9.2 [\#549](https://github.com/gobuffalo/buffalo/pull/549) ([markbates](https://github.com/markbates))
- ensure the dep installs the correct version of buffalo [\#546](https://github.com/gobuffalo/buffalo/pull/546) ([markbates](https://github.com/markbates))
- added a -d flag to buffalo dev to run the app with delve [\#545](https://github.com/gobuffalo/buffalo/pull/545) ([markbates](https://github.com/markbates))
- Fix resource generation example message [\#544](https://github.com/gobuffalo/buffalo/pull/544) ([slomek](https://github.com/slomek))
- Adding versions to the NPM packages [\#543](https://github.com/gobuffalo/buffalo/pull/543) ([paganotoni](https://github.com/paganotoni))
- updated to a regex to check for illegal characters [\#542](https://github.com/gobuffalo/buffalo/pull/542) ([markbates](https://github.com/markbates))
- Bugfix special characters [\#541](https://github.com/gobuffalo/buffalo/pull/541) ([paganotoni](https://github.com/paganotoni))
- removed the installation of dep in the integration tests [\#539](https://github.com/gobuffalo/buffalo/pull/539) ([markbates](https://github.com/markbates))
- removed the unused docker environments [\#538](https://github.com/gobuffalo/buffalo/pull/538) ([markbates](https://github.com/markbates))
- added an `info` command to print of env/versions to help with debugging [\#537](https://github.com/gobuffalo/buffalo/pull/537) ([markbates](https://github.com/markbates))
- use `--skip-dep` in the integration tests [\#536](https://github.com/gobuffalo/buffalo/pull/536) ([markbates](https://github.com/markbates))
- webpack should rely on webpack config not -p flag [\#535](https://github.com/gobuffalo/buffalo/pull/535) ([tonyalaribe](https://github.com/tonyalaribe))
- updated the Dockerfile.build to use 1.8.3 and fix the dep install [\#534](https://github.com/gobuffalo/buffalo/pull/534) ([markbates](https://github.com/markbates))
- added support for plugin aliases [\#533](https://github.com/gobuffalo/buffalo/pull/533) ([markbates](https://github.com/markbates))
- enhanced the simple worker with logging [\#531](https://github.com/gobuffalo/buffalo/pull/531) ([markbates](https://github.com/markbates))
- Use i18n in index template [\#529](https://github.com/gobuffalo/buffalo/pull/529) ([slomek](https://github.com/slomek))
- Buffalo can't bind on fixed IP fixes \#518 [\#527](https://github.com/gobuffalo/buffalo/pull/527) ([markbates](https://github.com/markbates))
- time out bad plugins fixes \#513 [\#526](https://github.com/gobuffalo/buffalo/pull/526) ([markbates](https://github.com/markbates))
- Assets missing on new v0.9.1.2 app without npm fixes \#517 [\#525](https://github.com/gobuffalo/buffalo/pull/525) ([markbates](https://github.com/markbates))
- added --no-progress to npm install [\#523](https://github.com/gobuffalo/buffalo/pull/523) ([markbates](https://github.com/markbates))
- Prevents buffalo from executing commands outside the buffalo project [\#514](https://github.com/gobuffalo/buffalo/pull/514) ([paganotoni](https://github.com/paganotoni))
- Update string and plain renderer with the charset [\#511](https://github.com/gobuffalo/buffalo/pull/511) ([ryanfaerman](https://github.com/ryanfaerman))

## [v0.9.1.2](https://github.com/gobuffalo/buffalo/tree/v0.9.1.2) (2017-07-12)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.9.1.1...v0.9.1.2)

**Closed issues:**

- Pagination and auto-generated route functions [\#501](https://github.com/gobuffalo/buffalo/issues/501)

**Merged pull requests:**

- Binary migrations [\#507](https://github.com/gobuffalo/buffalo/pull/507) ([markbates](https://github.com/markbates))
- use goroutines to more quickly look through paths for plugins [\#506](https://github.com/gobuffalo/buffalo/pull/506) ([markbates](https://github.com/markbates))
- Generate pagination [\#505](https://github.com/gobuffalo/buffalo/pull/505) ([markbates](https://github.com/markbates))
- cleaned up per request logging [\#504](https://github.com/gobuffalo/buffalo/pull/504) ([markbates](https://github.com/markbates))
- Allowing to pass extra arguments to named routes [\#502](https://github.com/gobuffalo/buffalo/pull/502) ([paganotoni](https://github.com/paganotoni))

## [v0.9.1.1](https://github.com/gobuffalo/buffalo/tree/v0.9.1.1) (2017-07-11)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.9.1...v0.9.1.1)

## [v0.9.1](https://github.com/gobuffalo/buffalo/tree/v0.9.1) (2017-07-10)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.9.0...v0.9.1)

**Implemented enhancements:**

- Add support for plugins for the buffalo binary [\#495](https://github.com/gobuffalo/buffalo/issues/495)
- Resource generator doesn't support "nested" values [\#364](https://github.com/gobuffalo/buffalo/issues/364)
- WIP 0.9.1 [\#499](https://github.com/gobuffalo/buffalo/pull/499) ([markbates](https://github.com/markbates))

**Closed issues:**

- error with flag --api [\#492](https://github.com/gobuffalo/buffalo/issues/492)
- No Benchmark ? Seriously ! [\#490](https://github.com/gobuffalo/buffalo/issues/490)
- New project tests fail \(on buffalo 0.9.0\) [\#489](https://github.com/gobuffalo/buffalo/issues/489)
- v2 of github.com/monoculum/formam always get installed [\#488](https://github.com/gobuffalo/buffalo/issues/488)
- template helper accessing context value [\#486](https://github.com/gobuffalo/buffalo/issues/486)
- cockroachdb compatibility [\#478](https://github.com/gobuffalo/buffalo/issues/478)
- Grifts generator [\#477](https://github.com/gobuffalo/buffalo/issues/477)
- Destroy Model is inconsistent with the generator for the model [\#475](https://github.com/gobuffalo/buffalo/issues/475)
- running deployed binary task argument missing [\#473](https://github.com/gobuffalo/buffalo/issues/473)
- dep wrong version of github.com/monoculum/formam [\#472](https://github.com/gobuffalo/buffalo/issues/472)
- adding mailer [\#292](https://github.com/gobuffalo/buffalo/issues/292)
- Provide alternative webpack starters [\#214](https://github.com/gobuffalo/buffalo/issues/214)

**Merged pull requests:**

- added a render.Plain Renderer [\#498](https://github.com/gobuffalo/buffalo/pull/498) ([markbates](https://github.com/markbates))
- Add support for plugins for the buffalo binary closes \#495 [\#496](https://github.com/gobuffalo/buffalo/pull/496) ([markbates](https://github.com/markbates))
- fix case-insensitive import collision [\#494](https://github.com/gobuffalo/buffalo/pull/494) ([s7v7nislands](https://github.com/s7v7nislands))
- error with flag --api fixes \#492 [\#493](https://github.com/gobuffalo/buffalo/pull/493) ([markbates](https://github.com/markbates))
- Update README.md Spelling Mistake [\#491](https://github.com/gobuffalo/buffalo/pull/491) ([asqrd](https://github.com/asqrd))
- allows n number of yields in templates [\#485](https://github.com/gobuffalo/buffalo/pull/485) ([markbates](https://github.com/markbates))
- Vdemeester fix binding parse time [\#484](https://github.com/gobuffalo/buffalo/pull/484) ([markbates](https://github.com/markbates))
- Making some fronted optimisation [\#481](https://github.com/gobuffalo/buffalo/pull/481) ([egorikas](https://github.com/egorikas))
- Moves destroy model command [\#480](https://github.com/gobuffalo/buffalo/pull/480) ([paganotoni](https://github.com/paganotoni))
- Feature grift generator [\#479](https://github.com/gobuffalo/buffalo/pull/479) ([paganotoni](https://github.com/paganotoni))
- improved the connection with grifts, including the ability automatically add grifts to a buffalo.App [\#476](https://github.com/gobuffalo/buffalo/pull/476) ([markbates](https://github.com/markbates))
- Session once [\#474](https://github.com/gobuffalo/buffalo/pull/474) ([markbates](https://github.com/markbates))
- Nested resources with actions namespacing [\#469](https://github.com/gobuffalo/buffalo/pull/469) ([paganotoni](https://github.com/paganotoni))

## [v0.9.0](https://github.com/gobuffalo/buffalo/tree/v0.9.0) (2017-06-19)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.8.2...v0.9.0)

**Implemented enhancements:**

- Add an HTTP basic auth middleware [\#432](https://github.com/gobuffalo/buffalo/issues/432)
- Proposal: Add a `config` package [\#394](https://github.com/gobuffalo/buffalo/issues/394)
- Remove or rebuild the apps in the examples folder [\#376](https://github.com/gobuffalo/buffalo/issues/376)
- Destroy CLI [\#277](https://github.com/gobuffalo/buffalo/issues/277)
- Create a Docker container to run Buffalo [\#250](https://github.com/gobuffalo/buffalo/issues/250)
- Buffalo - Heroku Buildpack [\#159](https://github.com/gobuffalo/buffalo/issues/159)
- Go package dependencies are untracked [\#133](https://github.com/gobuffalo/buffalo/issues/133)
- Add support for background jobs [\#95](https://github.com/gobuffalo/buffalo/issues/95)
- 0.9.0 Release [\#468](https://github.com/gobuffalo/buffalo/pull/468) ([markbates](https://github.com/markbates))

**Fixed bugs:**

- Worker package is not compatible with go1.7 [\#448](https://github.com/gobuffalo/buffalo/issues/448)
- Actions generator inserts test & route for already existing actions [\#412](https://github.com/gobuffalo/buffalo/issues/412)
- Fix 428: use proper case for ID in generated templates [\#431](https://github.com/gobuffalo/buffalo/pull/431) ([stanislas-m](https://github.com/stanislas-m))

**Closed issues:**

- Add an `--api` flag to `buffalo new` [\#462](https://github.com/gobuffalo/buffalo/issues/462)
- Need to set path of database.yml when production mode. [\#460](https://github.com/gobuffalo/buffalo/issues/460)
- When running go imports on generates it's slow with vendor directories [\#459](https://github.com/gobuffalo/buffalo/issues/459)
- `buffalo dev` stuck when running dev server for 1st time [\#454](https://github.com/gobuffalo/buffalo/issues/454)
- Test ticket from Code Climate [\#453](https://github.com/gobuffalo/buffalo/issues/453)
- problem on middleware skipping [\#445](https://github.com/gobuffalo/buffalo/issues/445)
- Allow action generator to specify the method to use on app.go [\#438](https://github.com/gobuffalo/buffalo/issues/438)
- Allow action generator to skip views. [\#437](https://github.com/gobuffalo/buffalo/issues/437)
- Wrong sqlite3 path: buffalo db create puts the database file into wrong directory [\#436](https://github.com/gobuffalo/buffalo/issues/436)
- Can I store an Object into Session? [\#433](https://github.com/gobuffalo/buffalo/issues/433)
- Documentation for the Recent Background Job [\#429](https://github.com/gobuffalo/buffalo/issues/429)
- treatment of field name '\_id' is different between model and template. [\#428](https://github.com/gobuffalo/buffalo/issues/428)
- If no database is being used the `nosqlite` build tag should be used [\#425](https://github.com/gobuffalo/buffalo/issues/425)
- Export `buffalo.Response` so it can be used outside of Buffalo [\#423](https://github.com/gobuffalo/buffalo/issues/423)
- Remove deprecated functions slated for 0.9.0 [\#418](https://github.com/gobuffalo/buffalo/issues/418)
- Applicaton.scss does not update built application.css [\#414](https://github.com/gobuffalo/buffalo/issues/414)
- Slow tests for `cmd/generate` [\#407](https://github.com/gobuffalo/buffalo/issues/407)
- Replace schema package [\#335](https://github.com/gobuffalo/buffalo/issues/335)

**Merged pull requests:**

- add fall throughs for finding the schema to set up the test db [\#471](https://github.com/gobuffalo/buffalo/pull/471) ([markbates](https://github.com/markbates))
- Added a flag to generate different Dockerfile types, as well as a sub-command to generate docker [\#470](https://github.com/gobuffalo/buffalo/pull/470) ([markbates](https://github.com/markbates))
- generate a multi-stage Dockerfile [\#467](https://github.com/gobuffalo/buffalo/pull/467) ([markbates](https://github.com/markbates))
- Added a --static flag to buffalo build [\#466](https://github.com/gobuffalo/buffalo/pull/466) ([markbates](https://github.com/markbates))
- read the database.yml from inside the binary fixes \#460 [\#465](https://github.com/gobuffalo/buffalo/pull/465) ([markbates](https://github.com/markbates))
- added a --skip-dep flag to not use dep [\#464](https://github.com/gobuffalo/buffalo/pull/464) ([markbates](https://github.com/markbates))
- Add an `--api` flag to `buffalo new` closes \#462 [\#463](https://github.com/gobuffalo/buffalo/pull/463) ([markbates](https://github.com/markbates))
- When running go imports on generates it's slow with vendor directories closes \#459 [\#461](https://github.com/gobuffalo/buffalo/pull/461) ([markbates](https://github.com/markbates))
- fixed an issue with binding if the time was an empty string it would raise an error [\#458](https://github.com/gobuffalo/buffalo/pull/458) ([markbates](https://github.com/markbates))
- App.Stop should take an error. [\#457](https://github.com/gobuffalo/buffalo/pull/457) ([markbates](https://github.com/markbates))
- cleaned up the basicauth package to prevent "stuttering" as well as make it nicer to use. [\#456](https://github.com/gobuffalo/buffalo/pull/456) ([markbates](https://github.com/markbates))
- Adds Formam as a replacement for Schema [\#455](https://github.com/gobuffalo/buffalo/pull/455) ([paganotoni](https://github.com/paganotoni))
- Set up github.com/golang/dep for new applications fixes \#133 [\#452](https://github.com/gobuffalo/buffalo/pull/452) ([markbates](https://github.com/markbates))
- hopefully this fixes skipping of middleware for resources. closes \#445 [\#451](https://github.com/gobuffalo/buffalo/pull/451) ([markbates](https://github.com/markbates))
- ignore csrf in test mode [\#450](https://github.com/gobuffalo/buffalo/pull/450) ([markbates](https://github.com/markbates))
- cleaned up the errors page output to show less "pointer garbage" [\#447](https://github.com/gobuffalo/buffalo/pull/447) ([markbates](https://github.com/markbates))
- generate a Dockerfile and .dockerignore file for new apps [\#446](https://github.com/gobuffalo/buffalo/pull/446) ([markbates](https://github.com/markbates))
- \[feature\] implementing some UI for the error and welcome pages [\#444](https://github.com/gobuffalo/buffalo/pull/444) ([paganotoni](https://github.com/paganotoni))
- fixed `http: multiple response.WriteHeader calls` [\#443](https://github.com/gobuffalo/buffalo/pull/443) ([jnoortheen](https://github.com/jnoortheen))
- cleaned up some non-idiomatic code [\#442](https://github.com/gobuffalo/buffalo/pull/442) ([markbates](https://github.com/markbates))
- implementing skip templates and method [\#440](https://github.com/gobuffalo/buffalo/pull/440) ([paganotoni](https://github.com/paganotoni))
- Moving action generator tests to use filetest [\#439](https://github.com/gobuffalo/buffalo/pull/439) ([paganotoni](https://github.com/paganotoni))
- Fix HTTP Basic browser authentication prompt [\#435](https://github.com/gobuffalo/buffalo/pull/435) ([peteretelej](https://github.com/peteretelej))
- Add HTTP Basic authentication Middleware [\#434](https://github.com/gobuffalo/buffalo/pull/434) ([peteretelej](https://github.com/peteretelej))
- Destroy resource CLI [\#430](https://github.com/gobuffalo/buffalo/pull/430) ([paganotoni](https://github.com/paganotoni))
- If no database is being used the `nosqlite` build tag should be used \#425 [\#426](https://github.com/gobuffalo/buffalo/pull/426) ([markbates](https://github.com/markbates))
- Export `buffalo.Response` so it can be used outside of Buffalo \#423 [\#424](https://github.com/gobuffalo/buffalo/pull/424) ([markbates](https://github.com/markbates))
- RequestLogger sets request\_id in context [\#421](https://github.com/gobuffalo/buffalo/pull/421) ([nzoschke](https://github.com/nzoschke))
- Add support for background jobs \#95 [\#420](https://github.com/gobuffalo/buffalo/pull/420) ([markbates](https://github.com/markbates))
- Remove deprecations [\#419](https://github.com/gobuffalo/buffalo/pull/419) ([markbates](https://github.com/markbates))
- Remove or rebuild the apps in the examples folder \#376 [\#417](https://github.com/gobuffalo/buffalo/pull/417) ([markbates](https://github.com/markbates))
- Fixes \#412  action generation [\#416](https://github.com/gobuffalo/buffalo/pull/416) ([paganotoni](https://github.com/paganotoni))
- Fix function name [\#415](https://github.com/gobuffalo/buffalo/pull/415) ([markbates](https://github.com/markbates))
- Fix typos in GoDoc [\#413](https://github.com/gobuffalo/buffalo/pull/413) ([nubunto](https://github.com/nubunto))

## [v0.8.2](https://github.com/gobuffalo/buffalo/tree/v0.8.2) (2017-05-01)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.8.1.2...v0.8.2)

**Implemented enhancements:**

- Resource generator should be able to generate XML [\#388](https://github.com/gobuffalo/buffalo/issues/388)
- Resource generator should be able to generate JSON [\#365](https://github.com/gobuffalo/buffalo/issues/365)
- Add support for "named routes" [\#210](https://github.com/gobuffalo/buffalo/issues/210)

**Closed issues:**

- Command failing on MacOS Sierra 10.12 [\#404](https://github.com/gobuffalo/buffalo/issues/404)
- Gracefully handle panics [\#401](https://github.com/gobuffalo/buffalo/issues/401)
- i18n should search for the language in the session as well [\#397](https://github.com/gobuffalo/buffalo/issues/397)
- babel deprecation warning during creation [\#395](https://github.com/gobuffalo/buffalo/issues/395)
- DB reverse tools  [\#393](https://github.com/gobuffalo/buffalo/issues/393)
- Generator New - can't get package [\#385](https://github.com/gobuffalo/buffalo/issues/385)
- field name got cut off if it ends with s when generate resource [\#384](https://github.com/gobuffalo/buffalo/issues/384)
- I've got "undefined: models in models.User" [\#383](https://github.com/gobuffalo/buffalo/issues/383)

**Merged pull requests:**

- Added a ldflags option to the build command [\#409](https://github.com/gobuffalo/buffalo/pull/409) ([tbasnett](https://github.com/tbasnett))
- specified a version for certain deprecations [\#408](https://github.com/gobuffalo/buffalo/pull/408) ([markbates](https://github.com/markbates))
- 0.8.2 [\#405](https://github.com/gobuffalo/buffalo/pull/405) ([markbates](https://github.com/markbates))
- gcc, node and npm dependecies [\#403](https://github.com/gobuffalo/buffalo/pull/403) ([danielsousaio](https://github.com/danielsousaio))
- Gracefully handle panics fixes \#401 [\#402](https://github.com/gobuffalo/buffalo/pull/402) ([markbates](https://github.com/markbates))
- add forbidden app names validation check [\#400](https://github.com/gobuffalo/buffalo/pull/400) ([marythought](https://github.com/marythought))
- i18n should search for the language in the session as well fixes \#397 [\#398](https://github.com/gobuffalo/buffalo/pull/398) ([markbates](https://github.com/markbates))
- babel deprecation warning during creation fixes \#395 [\#396](https://github.com/gobuffalo/buffalo/pull/396) ([markbates](https://github.com/markbates))
- build tags can now be passed in. oh, and sqlite will not be compiled [\#392](https://github.com/gobuffalo/buffalo/pull/392) ([markbates](https://github.com/markbates))
- Added XML resource generator, fixes \#388 [\#390](https://github.com/gobuffalo/buffalo/pull/390) ([stanislas-m](https://github.com/stanislas-m))
- Csrf cleanup [\#387](https://github.com/gobuffalo/buffalo/pull/387) ([markbates](https://github.com/markbates))
- Implements JSON resource generator [\#386](https://github.com/gobuffalo/buffalo/pull/386) ([stanislas-m](https://github.com/stanislas-m))
- Named Routes implementation. [\#332](https://github.com/gobuffalo/buffalo/pull/332) ([paganotoni](https://github.com/paganotoni))

## [v0.8.1.2](https://github.com/gobuffalo/buffalo/tree/v0.8.1.2) (2017-04-13)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.8.1...v0.8.1.2)

**Closed issues:**

- Default routes listing link to the route [\#381](https://github.com/gobuffalo/buffalo/issues/381)
- Add CSRF middleware handler to TestSuite [\#380](https://github.com/gobuffalo/buffalo/issues/380)
- migration with model [\#379](https://github.com/gobuffalo/buffalo/issues/379)
- Running outside of GOPATH fails [\#375](https://github.com/gobuffalo/buffalo/issues/375)
- Warning when `go get`-ing buffalo [\#374](https://github.com/gobuffalo/buffalo/issues/374)

**Merged pull requests:**

- updated the examples apps to use c.Value instead of the old c.Get function [\#377](https://github.com/gobuffalo/buffalo/pull/377) ([markbates](https://github.com/markbates))

## [v0.8.1](https://github.com/gobuffalo/buffalo/tree/v0.8.1) (2017-04-10)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.8.0...v0.8.1)

**Fixed bugs:**

- --skip-pop doesn't work [\#343](https://github.com/gobuffalo/buffalo/issues/343)
- Install in Arch Linux fails [\#311](https://github.com/gobuffalo/buffalo/issues/311)

**Closed issues:**

- sqlite cant store image blobs [\#372](https://github.com/gobuffalo/buffalo/issues/372)
- go get on macOS Sierra exits with Killed: 9 [\#371](https://github.com/gobuffalo/buffalo/issues/371)
- node deprecation warning for loader-utils [\#369](https://github.com/gobuffalo/buffalo/issues/369)
- Using Laravel Mix as a webpack wrapper for easier usage and sensible defaults [\#368](https://github.com/gobuffalo/buffalo/issues/368)
- Resource names don't show up correctly in routes tables [\#362](https://github.com/gobuffalo/buffalo/issues/362)
- Create project issue [\#359](https://github.com/gobuffalo/buffalo/issues/359)
- New Generator - System Init [\#357](https://github.com/gobuffalo/buffalo/issues/357)
- stdout on windows console. [\#348](https://github.com/gobuffalo/buffalo/issues/348)
- actions-packr.go: path error in windows [\#347](https://github.com/gobuffalo/buffalo/issues/347)
- Upgrading to V0.8X [\#346](https://github.com/gobuffalo/buffalo/issues/346)
- context.Logger\(\).InfoF or Errorf does not have color [\#344](https://github.com/gobuffalo/buffalo/issues/344)
- Generate the `middleware.CSRF` middleware in new applications [\#342](https://github.com/gobuffalo/buffalo/issues/342)
- Patch for a new project, "buffalo new" cmd on windows 10 [\#340](https://github.com/gobuffalo/buffalo/issues/340)
- Upgrade of gobuffalo and running go setup -u wants to reference a resolvers package that doesn't exist anymore [\#339](https://github.com/gobuffalo/buffalo/issues/339)
- Getting "Killed: 9" response from buffalo binary [\#338](https://github.com/gobuffalo/buffalo/issues/338)
- reflect.Set: value of type uuid.UUID is not assignable to type uuid.UUID [\#337](https://github.com/gobuffalo/buffalo/issues/337)
- Windows 10 new buffalo site sqlite3 "\<no value\>" in database.yml [\#336](https://github.com/gobuffalo/buffalo/issues/336)
- Integration test to resources fails, when basic implementation contains a model. [\#328](https://github.com/gobuffalo/buffalo/issues/328)
- Add basic functions into the default implementation of resource functions [\#297](https://github.com/gobuffalo/buffalo/issues/297)
- I18n \(internationalisation\) support [\#249](https://github.com/gobuffalo/buffalo/issues/249)

**Merged pull requests:**

- Merge Development branch for v0.8.1 [\#373](https://github.com/gobuffalo/buffalo/pull/373) ([markbates](https://github.com/markbates))
- generate app.go to have translations available in actions [\#370](https://github.com/gobuffalo/buffalo/pull/370) ([markbates](https://github.com/markbates))
- Typo in setup command long description [\#367](https://github.com/gobuffalo/buffalo/pull/367) ([jmhobbs](https://github.com/jmhobbs))
- Resource names don't show up correctly in routes tables fixes \#362 [\#363](https://github.com/gobuffalo/buffalo/pull/363) ([markbates](https://github.com/markbates))
- Translate in actions [\#361](https://github.com/gobuffalo/buffalo/pull/361) ([markbates](https://github.com/markbates))
- I18n \(internationalisation\) support closes \#249 [\#354](https://github.com/gobuffalo/buffalo/pull/354) ([markbates](https://github.com/markbates))
- Generate the `middleware.CSRF` middleware in new applications closes \#342 [\#353](https://github.com/gobuffalo/buffalo/pull/353) ([markbates](https://github.com/markbates))
- Resource generator fixes \#328 [\#352](https://github.com/gobuffalo/buffalo/pull/352) ([markbates](https://github.com/markbates))
- Build webpack assets for production in `buffalo build` [\#351](https://github.com/gobuffalo/buffalo/pull/351) ([mattdoran](https://github.com/mattdoran))
- Setup webpack to enable full ES6 transpile [\#350](https://github.com/gobuffalo/buffalo/pull/350) ([mattdoran](https://github.com/mattdoran))
- turn off colored output on windows in dev mode. fixes \#348 [\#349](https://github.com/gobuffalo/buffalo/pull/349) ([markbates](https://github.com/markbates))
- --skip-pop doesn't work fixes \#343 [\#345](https://github.com/gobuffalo/buffalo/pull/345) ([markbates](https://github.com/markbates))
- ignore case when checking if the user is in the GOPATH. fixes \#340 [\#341](https://github.com/gobuffalo/buffalo/pull/341) ([markbates](https://github.com/markbates))
- added a ParameterLogger middleware [\#334](https://github.com/gobuffalo/buffalo/pull/334) ([markbates](https://github.com/markbates))
- Resource generator --use-model \#297 fixes \#328 [\#333](https://github.com/gobuffalo/buffalo/pull/333) ([as27](https://github.com/as27))

## [v0.8.0](https://github.com/gobuffalo/buffalo/tree/v0.8.0) (2017-03-27)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.8.0.dev...v0.8.0)

## [v0.8.0.dev](https://github.com/gobuffalo/buffalo/tree/v0.8.0.dev) (2017-03-27)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.7.4.1...v0.8.0.dev)

**Implemented enhancements:**

- Potential enhancement -- version cmd [\#322](https://github.com/gobuffalo/buffalo/issues/322)
- Generate a default .gitlab-ci.yml file for new projects [\#181](https://github.com/gobuffalo/buffalo/issues/181)
- React+Babel+webpack generator [\#125](https://github.com/gobuffalo/buffalo/issues/125)

**Closed issues:**

- Test fails in v0.8.0 buffalo/cmd/generate on Windows [\#325](https://github.com/gobuffalo/buffalo/issues/325)
- Asset paths returning 404 [\#324](https://github.com/gobuffalo/buffalo/issues/324)
- Having trouble installing -- prob easy fix [\#317](https://github.com/gobuffalo/buffalo/issues/317)
- Update docs to reflect current state of template context. [\#312](https://github.com/gobuffalo/buffalo/issues/312)

**Merged pull requests:**

- Adds a version command to the CLI [\#330](https://github.com/gobuffalo/buffalo/pull/330) ([SpencerCDixon](https://github.com/SpencerCDixon))
- Issue \#325 fixed failing windows tests [\#327](https://github.com/gobuffalo/buffalo/pull/327) ([as27](https://github.com/as27))
- fix output path [\#323](https://github.com/gobuffalo/buffalo/pull/323) ([mario-seidel](https://github.com/mario-seidel))
- Params doc fix [\#321](https://github.com/gobuffalo/buffalo/pull/321) ([markbates](https://github.com/markbates))
- rewrote the default "routes" task so it doesn't need to import the tablewriter package any more [\#320](https://github.com/gobuffalo/buffalo/pull/320) ([markbates](https://github.com/markbates))
- no longer need to shell out to the grifts binary [\#319](https://github.com/gobuffalo/buffalo/pull/319) ([markbates](https://github.com/markbates))
- Gen seed script [\#318](https://github.com/gobuffalo/buffalo/pull/318) ([markbates](https://github.com/markbates))
- Added a "buffalo setup" command to make sure applications are setup and [\#316](https://github.com/gobuffalo/buffalo/pull/316) ([markbates](https://github.com/markbates))
- spelling fix - chosing -\> choosing [\#314](https://github.com/gobuffalo/buffalo/pull/314) ([corylanou](https://github.com/corylanou))
- V0.8.0 [\#269](https://github.com/gobuffalo/buffalo/pull/269) ([markbates](https://github.com/markbates))

## [v0.7.4.1](https://github.com/gobuffalo/buffalo/tree/v0.7.4.1) (2017-03-22)
[Full Changelog](https://github.com/gobuffalo/buffalo/compare/v0.7.4...v0.7.4.1)

**Implemented enhancements:**

- Drop go.rice support [\#290](https://github.com/gobuffalo/buffalo/issues/290)
- Provide a CSRF middleware [\#243](https://github.com/gobuffalo/buffalo/issues/243)
- Proposal: ActionSuite [\#242](https://github.com/gobuffalo/buffalo/issues/242)
- Clean up the generate package [\#165](https://github.com/gobuffalo/buffalo/issues/165)
- Improve resource generator to also generate the model and migration [\#42](https://github.com/gobuffalo/buffalo/issues/42)
- Added CSRF middleware \(fix issue \#243\) [\#271](https://github.com/gobuffalo/buffalo/pull/271) ([stanislas-m](https://github.com/stanislas-m))

**Fixed bugs:**

- The migrate command on a built binary doesn't work [\#289](https://github.com/gobuffalo/buffalo/issues/289)

**Closed issues:**

- New project - webpack error [\#310](https://github.com/gobuffalo/buffalo/issues/310)
- Lower case package dependency [\#308](https://github.com/gobuffalo/buffalo/issues/308)
- Context Session [\#306](https://github.com/gobuffalo/buffalo/issues/306)
- Best practices of model error handling  [\#304](https://github.com/gobuffalo/buffalo/issues/304)
- fatal error: concurrent map writes [\#303](https://github.com/gobuffalo/buffalo/issues/303)
- Test fails on Windows in router\_test.go [\#295](https://github.com/gobuffalo/buffalo/issues/295)
- sql: Scan error [\#288](https://github.com/gobuffalo/buffalo/issues/288)
- context.HasParam\(name\) ? [\#285](https://github.com/gobuffalo/buffalo/issues/285)
- Documented video link is broken [\#284](https://github.com/gobuffalo/buffalo/issues/284)
- Set up test db before running `buffalo test` [\#281](https://github.com/gobuffalo/buffalo/issues/281)
- Removed the dependency on TT \(and BoltDB\) [\#280](https://github.com/gobuffalo/buffalo/issues/280)
- Skip vendor directory during build when looking for rice files [\#278](https://github.com/gobuffalo/buffalo/issues/278)
- Remove the new relic middleware from the 'middleware' package [\#225](https://github.com/gobuffalo/buffalo/issues/225)

**Merged pull requests:**

- run templates through the markdown engine first before going through plush [\#307](https://github.com/gobuffalo/buffalo/pull/307) ([markbates](https://github.com/markbates))
- Fix \#181: generate default gitlab ci config [\#305](https://github.com/gobuffalo/buffalo/pull/305) ([stanislas-m](https://github.com/stanislas-m))
- Action and Model Suites [\#302](https://github.com/gobuffalo/buffalo/pull/302) ([markbates](https://github.com/markbates))
- Removed Resolvers and Cleaned up actions/render.go [\#301](https://github.com/gobuffalo/buffalo/pull/301) ([markbates](https://github.com/markbates))
- cleaned up a bunch of code that did stuff with multiple go paths [\#300](https://github.com/gobuffalo/buffalo/pull/300) ([markbates](https://github.com/markbates))
- use schema's to set up the test database [\#299](https://github.com/gobuffalo/buffalo/pull/299) ([markbates](https://github.com/markbates))
- Correcting router\_test.go fixes \#295 [\#296](https://github.com/gobuffalo/buffalo/pull/296) ([as27](https://github.com/as27))
- Typo in Resource documentation fixed [\#294](https://github.com/gobuffalo/buffalo/pull/294) ([as27](https://github.com/as27))
- Drop go.rice support in favor of Packr [\#293](https://github.com/gobuffalo/buffalo/pull/293) ([markbates](https://github.com/markbates))
- Remove new relic [\#283](https://github.com/gobuffalo/buffalo/pull/283) ([markbates](https://github.com/markbates))
- Removed the dependency on TT \(and BoltDB\) closes \#280 and closes \#281 [\#282](https://github.com/gobuffalo/buffalo/pull/282) ([markbates](https://github.com/markbates))
- Skip vendor directory during build when looking for rice files [\#279](https://github.com/gobuffalo/buffalo/pull/279) ([markbates](https://github.com/markbates))
- Get valid gopath if multiple [\#275](https://github.com/gobuffalo/buffalo/pull/275) ([fooflare](https://github.com/fooflare))
- Deplist [\#274](https://github.com/gobuffalo/buffalo/pull/274) ([markbates](https://github.com/markbates))
- Switched to using github.com/gobuffalo/makr for generators [\#273](https://github.com/gobuffalo/buffalo/pull/273) ([markbates](https://github.com/markbates))
- Clean up the generate package [\#272](https://github.com/gobuffalo/buffalo/pull/272) ([markbates](https://github.com/markbates))
- \[feature\] implementing model generation for resources [\#270](https://github.com/gobuffalo/buffalo/pull/270) ([paganotoni](https://github.com/paganotoni))
- now that 1.8 is released we should be testing against the latest Go as well [\#268](https://github.com/gobuffalo/buffalo/pull/268) ([markbates](https://github.com/markbates))
- Convert Buffalo to use Plush by default [\#267](https://github.com/gobuffalo/buffalo/pull/267) ([markbates](https://github.com/markbates))

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
- Adds Test cases for the travis/none ci-provider generation [\#212](https://github.com/gobuffalo/buffalo/pull/212) ([paganotoni](https://github.com/paganotoni))
- Remove all licensing from the top of `cmd` files closes \#208 [\#209](https://github.com/gobuffalo/buffalo/pull/209) ([markbates](https://github.com/markbates))
- Fixed issue \#203: multiple GOPATH are not supported by buffalo new. [\#205](https://github.com/gobuffalo/buffalo/pull/205) ([stanislas-m](https://github.com/stanislas-m))
- \[\#193\] adding line number to the Context\#Get warning [\#202](https://github.com/gobuffalo/buffalo/pull/202) ([paganotoni](https://github.com/paganotoni))
- \[\#194\] Avoid Context\#Get warnings on the newly created app. [\#201](https://github.com/gobuffalo/buffalo/pull/201) ([paganotoni](https://github.com/paganotoni))
- \[\#198\] moving resource generator to avoid pointers [\#200](https://github.com/gobuffalo/buffalo/pull/200) ([paganotoni](https://github.com/paganotoni))
- \[\#197\] generates better code inside actions [\#199](https://github.com/gobuffalo/buffalo/pull/199) ([paganotoni](https://github.com/paganotoni))
- need to peg webpack to 1.14.0 because 2.2.x doesn't work with the generated config [\#196](https://github.com/gobuffalo/buffalo/pull/196) ([markbates](https://github.com/markbates))
- Fixed actions generator: imports were missing. [\#192](https://github.com/gobuffalo/buffalo/pull/192) ([stanislas-m](https://github.com/stanislas-m))
- remove unused code [\#191](https://github.com/gobuffalo/buffalo/pull/191) ([philipithomas](https://github.com/philipithomas))
- convert Version to constant [\#190](https://github.com/gobuffalo/buffalo/pull/190) ([philipithomas](https://github.com/philipithomas))
- Add Golint and fix all issues [\#189](https://github.com/gobuffalo/buffalo/pull/189) ([philipithomas](https://github.com/philipithomas))
- better error printing. [\#186](https://github.com/gobuffalo/buffalo/pull/186) ([markbates](https://github.com/markbates))
- fixes \#183 [\#185](https://github.com/gobuffalo/buffalo/pull/185) ([amedeiros](https://github.com/amedeiros))
- \[feature\] adds .travis.yml when generating a new app [\#184](https://github.com/gobuffalo/buffalo/pull/184) ([paganotoni](https://github.com/paganotoni))
- Allow buffalo build to extract assets and put them into a zip file [\#180](https://github.com/gobuffalo/buffalo/pull/180) ([stanislas-m](https://github.com/stanislas-m))
- set a default session name based on the project when it is created [\#179](https://github.com/gobuffalo/buffalo/pull/179) ([markbates](https://github.com/markbates))
- fixed an issue where flash messages where not clearing properly [\#178](https://github.com/gobuffalo/buffalo/pull/178) ([markbates](https://github.com/markbates))
- Fixes \#173 [\#177](https://github.com/gobuffalo/buffalo/pull/177) ([amedeiros](https://github.com/amedeiros))
- Skipping test generation if test exists on action generation [\#176](https://github.com/gobuffalo/buffalo/pull/176) ([paganotoni](https://github.com/paganotoni))
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
- Implementing the `flash` helper [\#143](https://github.com/gobuffalo/buffalo/pull/143) ([paganotoni](https://github.com/paganotoni))
- Prevent new app creation from outside the Go workspace [\#141](https://github.com/gobuffalo/buffalo/pull/141) ([markbates](https://github.com/markbates))
- UnWrap HttpErrors in pop middleware and return them. [\#138](https://github.com/gobuffalo/buffalo/pull/138) ([lumost](https://github.com/lumost))
- preventing wrapping errors with point in error handler stack [\#137](https://github.com/gobuffalo/buffalo/pull/137) ([lumost](https://github.com/lumost))
- import buffalo when generating a new action [\#136](https://github.com/gobuffalo/buffalo/pull/136) ([lumost](https://github.com/lumost))
- First Attempt at \#105  [\#112](https://github.com/gobuffalo/buffalo/pull/112) ([paganotoni](https://github.com/paganotoni))

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
- Makes our tests run on Go 1.7 and 1.8 [\#131](https://github.com/gobuffalo/buffalo/pull/131) ([paganotoni](https://github.com/paganotoni))
- build\_path does not work for Windows closes \#124 [\#130](https://github.com/gobuffalo/buffalo/pull/130) ([markbates](https://github.com/markbates))
- Edit some typo [\#129](https://github.com/gobuffalo/buffalo/pull/129) ([janczer](https://github.com/janczer))
- Passing some issues from codeclimate [\#122](https://github.com/gobuffalo/buffalo/pull/122) ([paganotoni](https://github.com/paganotoni))
- Provide a mechanism to map status codes to error handles. Closes \#110 [\#121](https://github.com/gobuffalo/buffalo/pull/121) ([markbates](https://github.com/markbates))
- Warn user about SESSION\_SECRET variable? closes \#114 [\#119](https://github.com/gobuffalo/buffalo/pull/119) ([markbates](https://github.com/markbates))
- point people to npm docs if there is an issue running npm [\#118](https://github.com/gobuffalo/buffalo/pull/118) ([markbates](https://github.com/markbates))
- Default css js files and no logo in assets wo webpack [\#117](https://github.com/gobuffalo/buffalo/pull/117) ([fooflare](https://github.com/fooflare))
- Return the RouteInfo when mapping an endpoint. Also make it available in the request context [\#109](https://github.com/gobuffalo/buffalo/pull/109) ([markbates](https://github.com/markbates))
- missing "public/assets" box on new app without webpack closes \#104 [\#107](https://github.com/gobuffalo/buffalo/pull/107) ([markbates](https://github.com/markbates))
- Using gentronics to generate the templates and the actions [\#106](https://github.com/gobuffalo/buffalo/pull/106) ([paganotoni](https://github.com/paganotoni))
- Adds an Actions generator to the cmd package [\#103](https://github.com/gobuffalo/buffalo/pull/103) ([paganotoni](https://github.com/paganotoni))
- Cleaning up some docs [\#102](https://github.com/gobuffalo/buffalo/pull/102) ([CodyOss](https://github.com/CodyOss))
- Avoiding generating log folder when running tests. [\#101](https://github.com/gobuffalo/buffalo/pull/101) ([paganotoni](https://github.com/paganotoni))
- \[grift\] changing the task to be pointing the buffalo repo [\#100](https://github.com/gobuffalo/buffalo/pull/100) ([paganotoni](https://github.com/paganotoni))
- making buffalo call grift for the tasks [\#98](https://github.com/gobuffalo/buffalo/pull/98) ([paganotoni](https://github.com/paganotoni))
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