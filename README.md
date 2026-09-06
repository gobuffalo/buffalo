<p align="center"><img src="https://raw.githubusercontent.com/gobuffalo/buffalo/master/logo.svg" width="360"></p>

<p align="center">
<a href="https://pkg.go.dev/github.com/gobuffalo/buffalo"><img src="https://pkg.go.dev/badge/github.com/gobuffalo/buffalo" alt="PkgGoDev"></a>
<a href="https://github.com/gobuffalo/buffalo/actions/workflows/standard-go-test.yml"><img src="https://github.com/gobuffalo/buffalo/actions/workflows/standard-go-test.yml/badge.svg"></a>
<a href="https://goreportcard.com/report/github.com/gobuffalo/buffalo"><img src="https://goreportcard.com/badge/github.com/gobuffalo/buffalo" alt="Go Report Card" /></a>
<a href="https://www.codetriage.com/gobuffalo/buffalo"><img src="https://www.codetriage.com/gobuffalo/buffalo/badges/users.svg" alt="Open Source Helpers" /></a>
</p>

# Buffalo

A Go web development eco-system, designed to make your project easier.

Buffalo helps you to generate a web project that already has everything from front-end (JavaScript, SCSS, etc.) to the back-end (database, routing, etc.) already hooked up and ready to run. From there it provides easy APIs to build your web application quickly in Go.

Buffalo **isn't just a framework**; it's a holistic web development environment and project structure that **lets developers get straight to the business** of, well, building their business.

> I :heart: web dev in go again - Brian Ketelsen

## Versions

The current stable version of Buffalo core is v1 (`v1` branch).

Versions (branches):
* `main` is for the current mainstream development.
* `v1` is the current stable release.

## ⚠️ Important

Buffalo works only with Go [modules](https://blog.golang.org/using-go-modules). `GOPATH` mode is likely to break most of the functionality of the Buffalo eco-system. Please see [this blog post](https://blog.gobuffalo.io/the-road-to-1-0-requiring-modules-5672c6b015e5) for more information.

Also, the Buffalo team actively gives support to the last 2 versions of Go, which at the moment are Go 1.23 and 1.24. While Buffalo `may` work on older versions, we encourage you to upgrade to latest 2 versions of Go for a better development experience.

## Documentation

Please visit [http://gobuffalo.io](http://gobuffalo.io) for the latest documentation, examples, and more.

### Quick Start

- [Installation](https://gobuffalo.io/documentation/getting_started/installation)
- [Create a new project](https://gobuffalo.io/documentation/getting_started/new-project)
- [Tutorials](https://gobuffalo.io/documentation/tutorials/)

## Shoulders of Giants

Buffalo would not be possible if not for all of the great projects it depends on. Please see [SHOULDERS.md](SHOULDERS.md) to see a list of them.

### Templating

[github.com/gobuffalo/plush](https://github.com/gobuffalo/plush) - This templating package was chosen over the standard Go `html/template` package for a variety of reasons. The biggest of which is that it is significantly more flexible and easy to work with.

### Routing

[github.com/gorilla/mux](https://github.com/gorilla/mux) - This router was chosen because of its stability and flexibility. There might be faster routers out there, but this one is definitely the most powerful!

### Models/ORM (Optional)

[github.com/gobuffalo/pop](https://github.com/gobuffalo/pop) - Accessing databases is nothing new in web applications. Pop, and its command line tool, Soda, were chosen because they strike a nice balance between simplifying common tasks, being idiomatic, and giving you the flexibility you need to build your app. Pop and Soda share the same core philosophies as Buffalo, so they were a natural choice.

### Sessions, Cookies, WebSockets, and more

[github.com/gorilla](https://github.com/gorilla) - The Gorilla toolkit is a great set of packages designed to improve upon the standard library for a variety of web-related packages. With these high-quality packages Buffalo can keep its "core" code to a minimum and focus on its goal of gluing them all together to make your life better.

## Benchmarks

Oh, yeah, everyone wants benchmarks! What would a web framework be without its benchmarks? Well, guess what? I'm not giving you any! That's right. This is Go! I assure you that it is plenty fast enough for you. If you want benchmarks you can either a) check out any benchmarks that the [GIANTS](SHOULDERS.md) Buffalo is built upon having published, or b) run your own. I have no interest in playing the benchmark game, and neither should you.

## Contributing

First, thank you so much for wanting to contribute! It means so much that you care enough to want to contribute. We appreciate every PR from the smallest of typos to the be biggest of features.

**Here are the core rules to respect**:

- If you have any question, please consider using the
  [Slack channel](https://gophers.slack.com/messages/buffalo/) (-#buffalo-,
  *#buffalo_fr* or *#buffalo-dev* for contribution related questions) or
  [Stack Overflow](https://stackoverflow.com/questions/tagged/buffalo).
  We use GitHub issues for **bug reports and feature requests only**.
- All contributors of this project are working on their free time: be patient
  and kind. :-
- Consider opening an issue **BEFORE** creating a Pull request (PR): you won't
  lose your time on fixing non-existing bugs, or fixing the wrong bug. Also we
  can help you to produce the best PR!
- Open a PR against the `main` branch if your PR is for mainstream or version
  specific branch e.g. `v1` if your PR is for specific version.
  Note that the valid branch for a new feature request PR should be `main`
  while a PR against a version specific branch are allowed only for bugfixes.

For the full contribution guidelines, please read [CONTRIBUTING](.github/CONTRIBUTING.md).


## 🌐 Web Resources & Interactive Index
- [HEXA PUZZLE MASTER](https://learnquester.pages.dev/hexa-puzzle-master.html)
- [CATEGORY COOKING46](https://iskillquest.pages.dev/category-cooking46.html)
- [MOBILE LEGENDS SLIME 3V3](https://thelearnquesters.pages.dev/mobile-legends-slime-3v3.html)
- [CATEGORY GROW](https://enskillcrafts.pages.dev/category-grow.html)
- [CATEGORY MATH29](https://thequizzone.pages.dev/category-math29.html)
- [SOKOBAN PR](https://thequizzone.pages.dev/sokoban-pr.html)
- [LOL SURPRISE OMG BB DRIVER](https://thequizzone.pages.dev/lol-surprise-omg-bb-driver.html)
- [MATH LAVA TOWER RACE](https://thequizzone.pages.dev/math-lava-tower-race.html)
- [FROM NERD TO SCHOOL POPULAR](https://quizverses-9d2f2.web.app/from-nerd-to-school-popular.html)
- [LINKLINK](https://learnquester.github.io/linklink.html)
- [CATEGORY BATTLE](https://learnquester.github.io/category-battle.html)
- [REAL IMPOSSIBLE SKY TRACKS CAR DRIVING](https://studyplayings.web.app/real-impossible-sky-tracks-car-driving.html)
- [CATEGORY GROW GAMES](https://studyplayings.web.app/category-grow-games.html)
- [DRAW CLIMBER](https://thequizzone.pages.dev/draw-climber.html)
- [BALL AND GIRLFRIEND](https://thequizzone.pages.dev/ball-and-girlfriend.html)
- [OCEAN POP](https://thequizzone.pages.dev/ocean-pop.html)
- [INDEX29](https://themindzone.pages.dev/index29.html)
- [CATEGORY STICKMAN 3](https://thequizzone.pages.dev/category-stickman-3.html)
- [SOLAR SMASH](https://thequizzone.pages.dev/solar-smash.html)
- [CRAFT OF WARS](https://themindplay.pages.dev/craft-of-wars.html)
- [BUS PARKING OUT](https://quizverses.github.io/bus-parking-out.html)
- [PIRATES MATCH THE LOST TREASURE](https://themindplaying.web.app/pirates-match-the-lost-treasure.html)
- [CLUB TYCOON IDLE CLICKER](https://thequizzone.pages.dev/club-tycoon-idle-clicker.html)
- [CATEGORY QUIZ](https://skillplay.github.io/category-quiz.html)
- [CLEAN HOUSE CLEARING TRASH AND DIRT](https://iskillquest.pages.dev/clean-house-clearing-trash-and-dirt.html)
- [CRYPTOGRAPH](https://thequizzone.pages.dev/cryptograph.html)
- [VALENTINE S DAY COUPLE DATE](https://skillplay.github.io/valentine-s-day-couple-date.html)
- [TILE FARM STORY MATCHING GAME](https://studyplayings.web.app/tile-farm-story-matching-game.html)
- [BELL MADNESS](https://studyplaying.github.io/bell-madness.html)
- [GLOVES GROW RUSH](https://studyquests.github.io/gloves-grow-rush.html)
- [MAJESTIC DRAGONS MERGE](https://thequizzone.pages.dev/majestic-dragons-merge.html)
- [MINI GAMES CASUAL COLLECTION](https://studyquests.pages.dev/mini-games-casual-collection.html)
- [CATEGORY STICKMAN 2](https://studyquesthub.web.app/category-stickman-2.html)
- [CODE MAZE](https://thequizzone.pages.dev/code-maze.html)
- [LIPSTICK COLLECTOR RUN](https://studyquests.github.io/lipstick-collector-run.html)
- [HOME BLOCK STORY](https://studyplayings.pages.dev/home-block-story.html)
- [STICK NINJA SURVIVAL](https://studyplaying.github.io/stick-ninja-survival.html)
- [CATEGORY BRAIN261](https://studyplayings.web.app/category-brain261.html)
- [MINI GAMES RELAX COLLECTION 2](https://studyquests.github.io/mini-games-relax-collection-2.html)
- [GUMMY MERGE](https://studyquests.pages.dev/gummy-merge.html)
- [SOLITAIRE MAHJONG FARM 2](https://studyplaying.github.io/solitaire-mahjong-farm-2.html)
- [PIZZA PUZZLE](https://themindplays.pages.dev/pizza-puzzle.html)
- [FLAG PUZZLE JAM COLLECT FLAGS](https://themindplaying.web.app/flag-puzzle-jam-collect-flags.html)
- [ICONIC HALLOWEEN COSTUMES](https://thequizzone.pages.dev/iconic-halloween-costumes.html)
- [ZOMBIE SIEGEIO](https://studyplayings.web.app/zombie-siegeio.html)
- [TERMS](https://studyplayings.pages.dev/terms.html)
- [CATEGORY BASKETBALL 2](https://studyplayings.web.app/category-basketball-2.html)
- [SHAPE TRANSFORM RACE](https://studyplayings.web.app/shape-transform-race.html)
- [GEOMETRY RUSH 4D](https://studyquests.github.io/geometry-rush-4d.html)
- [HELIX CRUSH](https://studyquests.github.io/helix-crush.html)
- [NETQUEL COM](https://thequizzone.pages.dev/netquel-com.html)
- [CATEGORY STICKMAN175](https://thequizzone.pages.dev/category-stickman175.html)
- [ESCAPE STEAL BRAINROT SAHUR HILLS](https://themindplay.github.io/escape-steal-brainrot-sahur-hills.html)
- [CATEGORY SPEED158](https://studyquests.pages.dev/category-speed158.html)
- [LAST WAR SURVIVAL](https://themindplays.pages.dev/last-war-survival.html)
- [CATEGORY BATTLE ROYALE25](https://themindplaying.web.app/category-battle-royale25.html)
- [CATEGORY SPORTS](https://studyplayings.pages.dev/category-sports.html)
- [CATEGORY IO](https://studyplayings.web.app/category-io.html)
- [PRINXY HOUSE OF FASHION](https://thequizzone.pages.dev/prinxy-house-of-fashion.html)
- [POLICE CHASE DRIFTER](https://thequizzone.pages.dev/police-chase-drifter.html)
- [CATEGORY CASUAL](https://studyplayings.web.app/category-casual.html)
- [FPS TOY REALISM](https://studyplayings.web.app/fps-toy-realism.html)
- [SORT BALLS CONES](https://thequizzone.pages.dev/sort-balls-cones.html)
- [BELOTE 3IN1](https://studyplayings.web.app/belote-3in1.html)
- [CATEGORY AGILITY 3](https://studyquesthub.web.app/category-agility-3.html)
- [ANACONDA RUNNER](https://studyquests.github.io/anaconda-runner.html)
- [JUMP MAN](https://studyplayings.pages.dev/jump-man.html)
- [CATEGORY CRAFTING45](https://themindplaying.web.app/category-crafting45.html)
- [LAMPHEAD](https://studyplaying.github.io/lamphead.html)
- [PAWS PALS DINER](https://studyquests.github.io/paws-pals-diner.html)
- [CATEGORY SPACE57](https://studyplaying.github.io/category-space57.html)
- [CATEGORY MMO25](https://themindzone.pages.dev/category-mmo25.html)
- [CATEGORY SPORTS](https://themindplays.pages.dev/category-sports.html)
- [MERMAIDS SPOT THE DIFFERENCES](https://studyplayings.pages.dev/mermaids-spot-the-differences.html)
- [PIMPLE SQUEEZE](https://studyplaying.github.io/pimple-squeeze.html)
- [3 TILES](https://studyplayings.web.app/3-tiles.html)
- [CATEGORY UNBLOCKERS](https://studyquests.github.io/category-unblockers.html)
- [INDEX21](https://studyplayings.web.app/index21.html)
- [ASMR NAIL TREATMENT](https://studyplaying.github.io/asmr-nail-treatment.html)
- [CATEGORY UNBLOCKER](https://studyquests.github.io/category-unblocker.html)
- [INDEX11](https://learnquester.github.io/index11.html)
- [DIRTY MONEY THE RICH GET RICH](https://thequizzone.pages.dev/dirty-money-the-rich-get-rich.html)
- [CATEGORY FOOD](https://skillplay.github.io/category-food.html)
- [CRAZYZOMBIES 3D](https://themindplay.pages.dev/crazyzombies-3d.html)
- [CATEGORY MOBILE2 095](https://themindzone.pages.dev/category-mobile2-095.html)
- [DOGE RUSH DRAW HOME PUZZLE](https://thequizzone.pages.dev/doge-rush-draw-home-puzzle.html)
- [CATEGORY RPG](https://studyplaying.github.io/category-rpg.html)
- [CATEGORY POINT AND CLICK124](https://themindzone.pages.dev/category-point-and-click124.html)
- [PATTERNS](https://studyquests.pages.dev/patterns.html)
- [EMPIRE CITY](https://themindplay.pages.dev/empire-city.html)
- [CATEGORY MONSTER206](https://studyplaying.github.io/category-monster206.html)
- [CATEGORY UNBLOCKED GAMES](https://thequizzone.pages.dev/category-unblocked-games.html)
- [BUS COLLECT](https://themindplaying.web.app/bus-collect.html)
- [SANDSTORM COVERT OPS](https://skillplay.github.io/sandstorm-covert-ops.html)
- [FIDGET TOYS POP IT](https://studyplayings.web.app/fidget-toys-pop-it.html)
- [OBBY VS ZOMBIES](https://studyplayings.web.app/obby-vs-zombies.html)
- [CAPYBARA BLOCK BLAST](https://studyplaying.github.io/capybara-block-blast.html)
- [JUST LUDO](https://studyquests.github.io/just-ludo.html)
- [STICKMAN VS ZOMBIES EPIC FIGHT](https://iskillquest.pages.dev/stickman-vs-zombies-epic-fight.html)
- [CATEGORY STICKMAN175](https://iskillquest.pages.dev/category-stickman175.html)
- [YUMMY TRAILS](https://thequizzone.pages.dev/yummy-trails.html)
- [CATEGORY MMO24](https://studyquests.pages.dev/category-mmo24.html)
- [LABUBU DOLL MUKBANG ASMR UNBLOCKED](https://studyplaying.github.io/labubu-doll-mukbang-asmr-unblocked.html)
- [SUMMER TRIPLE MAHJONG](https://studyplaying.github.io/summer-triple-mahjong.html)
- [METRO ESCAPE](https://studyplaying.github.io/metro-escape.html)
- [CATEGORY FASHION105](https://themindzone.pages.dev/category-fashion105.html)
- [POCKET PARKING](https://iskillquest.pages.dev/pocket-parking.html)
- [SPIDER ROPE HERO CITY FIGHT](https://studyquesthub.web.app/spider-rope-hero-city-fight.html)
- [DRAW BRIDGE CHALLENGE](https://studyquests.pages.dev/draw-bridge-challenge.html)
- [EATING SIMULATOR](https://studyquesthub.web.app/eating-simulator.html)
- [INDEX10](https://themindzone.pages.dev/index10.html)
- [HOLE DEFENSE](https://thequizzone.pages.dev/hole-defense.html)
- [MIRROR SHAPE](https://studyplayings.pages.dev/mirror-shape.html)
- [BUBBLE SHOOTER VINTAGE](https://studyquesthub.web.app/bubble-shooter-vintage.html)
- [ROBLOX CHRISTMAS DRESSUP](https://iskillquest.pages.dev/roblox-christmas-dressup.html)
- [FURRY KUNG FU](https://studyquesthub.web.app/furry-kung-fu.html)
- [100 DOORS PUZZLE BOX](https://iskillquest.pages.dev/100-doors-puzzle-box.html)
- [RESIDENT EVIL PURGE OPERATION](https://thequizzone.pages.dev/resident-evil-purge-operation.html)
- [CATEGORY HORROR](https://themindzone.pages.dev/category-horror.html)
- [KABOOM MINER](https://iskillquest.pages.dev/kaboom-miner.html)
- [INDEX24](https://themindzone.pages.dev/index24.html)
- [CATEGORY MAHJONG CONNECT](https://studyquests.pages.dev/category-mahjong-connect.html)
- [APOCALYPSE SHELTER](https://themindplay.pages.dev/apocalypse-shelter.html)
- [HEXA STACK CHRISTMAS](https://studyplaying.github.io/hexa-stack-christmas.html)
- [ZOMBIE DRIFT 3D](https://studyplayings.web.app/zombie-drift-3d.html)
- [BATTLEDUDES IO](https://thequizzone.pages.dev/battledudes-io.html)
- [CATEGORY STRATEGY](https://studyquests.pages.dev/category-strategy.html)
- [3D ACRYLIC NAIL NAIL ART GAME](https://studyquesthub.web.app/3d-acrylic-nail-nail-art-game.html)
- [CATEGORY JUMPING147](https://themindzone.pages.dev/category-jumping147.html)
- [CATEGORY LOGIC538](https://skillplay.github.io/category-logic538.html)
