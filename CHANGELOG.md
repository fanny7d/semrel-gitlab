# CHANGELOG

<!--- next entry here -->

## 0.21.1
2019-10-27

### Fixes

- remove bump inconsistency on pre-releases (dd642bec896bd968debd632c8492f461a5a83f87)
- inconsistent error handling (cefc313c94599966e39f35bef6b6cbdd1a9fa430)

## 0.21.0
2019-08-14

### Features

- add option for listing other changes (2b93f8bf8b51d434666998813389cf9655aadfb9)

## 0.20.4
2019-05-14

### Fixes

- update juranki/go-semrel (2d302e9e1c836af1cb71bf7c0519ef4bb450c0ff)
- handle tag prefix when analyzing repo (065e8cd84670f68750226def3b817a4b6bf2abf5)

## 0.20.3
2019-05-07

### Fixes

- correct url to  uploaded file (e348112c11024e2d338721eb949676b1281fb659)

## 0.20.2
2019-04-12

### Fixes

- file upload test (460da4e3edbfb31e2b5e4ccda59fab3c9e405a05)
- use proxy settings from environment (020b909b3a1580bbcf217b5413d893dc213d9eaa)

## 0.20.1
2019-04-02

### Fixes

- add full url after using add-download, not the relative path returned by the API (d05b116c3f3c118210158dc6e9988102c6f3f20a)

## 0.20.0
2019-03-07

### Features

- get gitlab server version (#67) (cbd94d01114bf923df7f97b9a976d34ccbd4d564)
- detect availibility of releases api (#67) (489c6e020b09561a917fc8e4a67edec2ea8d21e5)
- create release using releases api (#67) (4ed263889a2ca922fbfb87fe5ea1b4e9c643d31c)
- add links using releaselinks api (#67) (dddb2ea149ae0ae860e853fd8af11bf4c1bb2eef)

### Fixes

- **doc:** major -> minor (28e7ed97d6c12a2463ae7ed668f5318da0628e8b)
- endpoint determination (774eb01276865124776373ea1ae1edd4ea4411c0)
- add win and mac downloads (b2a8031a73c2c6c41ef1024a45f3c8c90ef88067)
- update go-gitlab and releases mods (64b6851d11a42ced4340053591469aa1ead7dcde)
- upgrade github.com/xanzy/go-gitlab (c8a24c6eae1c8e7a05e176db93389036bf873699)

## 0.19.1
2018-10-25

### Fixes

- update github.com/xanzy/ssh-agent to fix windows build (#63) (2ae290d0a6a098d43729e80021a211e7fecca306)
- retry pipeline creation if reference was not found #65 (09c6a412496959cd6b73d9e37d2eb3354bb257c7)
- fix --create-tag-pipeline option (8cf85eee5b8d72b8a2dc01c3a02b997132a831d8)
- update go-git and go-billy to fix escaping on windows (e771b9ee5b30fc30b5d6fe576e1e1be7e7a90ff2)

## 0.19.0
2018-09-29

### Breaking changes

#### Automatic creation of tag pipeline was removed (8bb0ac4ae9d7e864845fe20379ce982fa12c8222)

Bump commit message is now
configurable, and it might not contain skip ci.

*--create-tag-pipeline* option was added
to *commit-and-tag* for the cases when pipeline
creation is really needed.

### Features

- configuration option for tag-prefix (#61) (ebeceb4b2d8ae5b8b97e4b7a0e8f436652a065ae)
- add bump-commit-tmpl option (#60) (06012561d01d660d773341668f0e7550ad00d9e7)

### Fixes

- typo in tag prefix handling (25d1a03b203e829a4f2de44aa8d5334a936080b3)

## 0.18.0
2018-09-11

### Features

- add configurable pre-release scheme (839cdcc1d9c3c0cff31775738e57343f97f7ce93)

### Fixes

- move case transform out of list parsing (4d72ab3d5184a8c82dd715417d228f6d664ed875)
- update dependencies (b9cd5f345d93a801ea4b818d15e198bf1e4b82b6)
- remove duplicate entries from pre-release notes (e0e132e20120d41479f92765e234fe2fc6474c20)

## 0.17.2
2018-09-05

### Fixes

- handle empty string properly when parsing comma separated list (d2cc321fc8e977a3772975c94cbb6e71ea472c0e)

## 0.17.1
2018-09-04

### Fixes

- fix list parsing for --release-branches (52ab5c8d42cc112cc38d5d0d770114151dc31faa)

## 0.17.0
2018-09-01

### Features

- infer api url from CI_PROJECT_URL if GL_API is not defined (3c3f0ff23ccebfb02dedc0cecdaee33f63e31f87)

### Fixes

- check env vars before starting on commit-and-tag (3ddda5ff04ca6ddb078c02abf1cc3c9b3b75c2b4)
- prevent sigsegv when resp == nil (c14ebf1d4daf5908c8252641f6fb14f38861a4f2)
- remove duplicate error output (364dd3a99f8093f4b338c0a14d54be44fc4de5b7)
- add parameter checks (a36b232a277b3b8746f787f2ba929e804ee94183)
- make changelog update more robust #52 (b20913fd41b73f931a7a87a3921f1582de2c4cfa)

## 0.16.0
2018-08-29

### Breaking changes

#### fix the typo in the prefix of env vars (14907322017a8122e5d2e5191ebe0d7d31cb0c82)

the SGS_* env vars change to GSG_*

### Features

- add cli flag and env var to stop initial development (#41) (db97f95671a57618f0bc596cf220587befc3dfc1)
- add global bump-patch (b98020d8579d53e5db5bb32bb5d889c036e12bcd)
- add release branches option (d18ad51ad82f11c15f7ad92af3cf3eaa5a135ca5)

### Fixes

- add deprecation warning for SGS_* env vars (d14d857831ec8cc0f6c2cae0d78da92e298ff9fb)
- remove unnecessary logging (3fbe70ceeda2aaad4cfa45ee8f272f71f8c15bed)
- update dependences (2b09d0a1c94c7a45384117b23f11e06ec3e5a030)
- don't use build metadata in pre-release version (131e99b5a19e7816f088e335dc93c3a701a731ab)
- use commit timestamp for pre-release (2e47efeebfdf337620a7af8482791c327272f705)
- use utc instead of local time in pre-release version (4c2d852086fcc0ce92044d9b2caebf3e9de9a770)
- fix pipeline creation error when there are no jobs defined #50 (a8c7213a888d4d9752040b4ed1a7014bc5ddf897)

## 0.15.0
2018-08-19

### Breaking changes

#### fix yesterdays brain fart with option names (#47) (79b6a386eb3bf1c4b4b278c81ec568d02eee7a8a)

If someone already started to use v0.14.0,
please upgrade immediately.

This release fixes the names of commit type options.

## 0.14.0
2018-08-18

### Breaking changes

#### remove unused flag --ci-project-url (32e9d3cc41665cbc3f215928957bcb17c212faa1)

the command execution will fail, if --ci-project-url
is specified. The flag was never really used, so removing it will not
change the behaviour of the command.

### Features

- recover from some transient errors and cleanup on failure (00c86dccf4b64c7849a9a78b6e01e13f3677fbbc)
- turn on cpu prof when GO_SEMREL_GITLAB_CPUPROF=filename is set (fd257b92e334595bd128d98c97a2d43ef45009b6)
- add options for customizing commit types (#43, #47) (f202f4fad9b0e3924489197b251f25ccf9bd6b20)

### Fixes

- **intro:** improve wording of semantic-release comparison (bc2b7aaba07331066320f7ba0b2f7a6a995a30f9)
- **download:** add warning about refactor in 1.3.1 (5df15a33f0cc54c13c6d36119074df1a232c4247)
- **help:** Improve help texts. (c68d8d97af30a19cb529955ab4b3ed398e38ee9e)
- fix typo with version numbers (ef0ff3a9dfb93459a042a1516840fc16f353304e)
- implement api test as an action (3c9b70b65c9c14fb9990e39493e6e7fd33655e52)

## 0.13.1
2018-07-15

### Fixes

- **site:** update site structure (834945c3c00384dbb1132a32726f480252363e98)
- **site:** remove tautology (b5fb2ac5f6cd1f7d43fd29bee338a07230d6aad6)
- **site:** correct link for ci example (336851d373475b7e14107ab87bbf8cfe5cc32a4b)
- **site:** fix link to ci script example (ab832b649700cbdc3a230f35bd4d0d60e4f1a8db)
- add minimal workflow coordination (faf0414f37cef963766741e787fac9d046ab27fa)
- improve naming of go-gitlab wrapper (91dd8f8aa87b4a0453dd971ad06bf0913a0594e9)
- implement idemponent undoable actions for tag and get-tag (9bbf0356a676fe88e1090ad64859fbe15ee33c32)
- flatten struct hierachy (77f8cbb67cf6d53b1b1013d50fe34d1e6a42a441)
- **actions:** naming, more actions, linking actions... (6ca579f69635066039acf5abd6a0dd8fefd5dbbc)
- **actions:** switch to using actions api (65d00257d5b1a6eb310ddb7fb0b8b48ade7b40a3)
- remove dead code from gitlabutil (a6423919c01b7fc48ed12c50e65d908e4071e1ae)
- set executable bit for downloaded binary (eeae50a292b4c1f1b2a036642a819625886f5655)
- logic error in commit workflow (724bcead459db6c4ddfd194d496bba69bdf72235)
- print error message in workflow.Apply (eb6398c5920b711031850931d91b16418c5cb0dd)
- return proper error from workflow.Apply (efb1bf152bb43ed26ba5d09edd1ca35882cf04f4)

## 0.13.0
2018-07-06

### Features

- add basic authentication (dcb8d8a1dec2a9ad8fafa125def915461906780f)
- add download link -api (29ae81d5f85326e2f5e5ea43c61dc2c7a9443f93)
- add-download-link command (6ba3030680236110269276232165a602995eb974)

### Fixes

- unify error handling (20c03c9698a86ae49dfdac935f95134f88809771)
- upgrade xanzy/go-gitlab package (4442529ee24b83a161c2c41eaf892fbf9b9ad485)
- update juranki/go-semrel (082a7d36539248d66e861a5d9c48a78bb2db8aff)
- add first tests for gitlab api helpers (bc2213228784fe6ee0e4bfbd207a879c6f7831bc)
- add test stage to ci (acf3abdf6385662ce107f08fe2a4353d1080d264)
- add gitlab test image (ea1e6ed8a1e823568b80592ed24cbfdfa1bf15d1)
- external url of gitlab test image (386f7645571bfdf7695c781d2b83e08c9677ef46)
- gitlab-ci indentation (fa5df0387570afa72fb32029b90eb7ce36364132)
- use modified gitlab image for tests (c37fb4ede6509222a52c0baaae4e2e16aec92321)
- wait for gitlab to respond before running tests (8d698d3e144979eff97dfafa0944b10bd5328a1a)
- ping api instead of ui, before tests (288824693f0ee1bb8471d3030b995e5354d0dc7b)
- try to add more time for gitlab to start (0a287b5dddd9471d343d9c856bac06f849097b76)
- wrap description in quotes (6337e689e8a9744c2943bacbbea7d912822540c2)
- help message of add-download-link (87e76e40aa6dec243b308b8ba8e0f9d410b6702a)

## 0.12.2
2018-06-26

### Fixes

- move markdown rendering to separate module (eaa278fec7619937c30acd746ed29cbe4a33f1ce)
- add helper for creating GL client (a378b87b3c6e3f87667ecd5508e457c6384e61a3)
- add helper for analyzing commits (a5d2cdce75bb89f355502fa7aad2f2ba056a89fc)
- unify mapping of CLI spec and implementation (b0634897f132fd8fc8d0e469326ccd320e511f14)
- unify cli param handling and (b754ae4335b69f2c17f960202575aec1c630af7e)

## 0.12.1
2018-06-24

### Fixes

- change version bump commit message to be neutral about ordering (a41717c69f2fe1f5a484ec4e170f81ed7295553e)

## 0.12.0
2018-06-24

### Breaking changes

#### unify ci env-var override flag naming (69485bc01afd4fff018cd2d2b1de7efc0e36b0a7)

Prefix CI env-var override flags with 'ci-'.

- project-path -> ci-project-path
- project-url -> ci-project-url
- commit-sha -> ci-commit-sha
- commit-tag -> ci-commit-tag
- commit-ref-name -> ci-commit-ref-name

### Features

- **next-version:** flag to return curr-ver if no changes detected (80895de90141d57f34f4320eab8fa54bb253a59e)
- add tag-and-commit #40 (47ac64cb7deaafcc283ac8e89a90e291e89e74ff)

### Fixes

- more accurate help message for changelog command (ccac254e2977f4b07100faa3ca7885aef9a1ae5b)
- fix behaviour of next-version (68811d2778e2dccaa4d96ed88aed9b04b98e4b52)
- upgrade to go 1.10 (1f3a78941c41dafdb2136a68b141dad22617c4f5)
- use stable release of docker image (5b78c97246855ff3c6c893d3469c022bc8802043)
- **site:** add help for tag-and-commit (a9d846ffb33c46e94c98eee0c92d47d1c9438a25)
- **site:** update status chapter (ae63b4df70cb332370f648231b45e987b88d7c65)

## 0.11.9
2018-05-30

### Fixes

- **site:** remove RSS (not used and contains invalid links) (400efde46e91feaa6487cdd9e5d67f582c623254)

## 0.11.8
2018-05-30

### Fixes

- **site:** add description and keywords meta (7421ab5025aef95d46b933937454d2428e28b48f)

## 0.11.7
2018-05-23

### Fixes

- remove unnecessary output (new sha) (851ee83e0f5ab720c32130390b1cfa30dd3f6aac)
- close #36 (758f86910a47d0388112442f0c971ac0cf6a5b5a)
- remove warning about recent refactor (87b7c3572a4e7d02d7183c18578964b0240e5351)

## 0.11.6
2018-05-15

### Fixes

- add comparison to semantic-release (2e9851f1ea05cb253059fe844d79fb836b4dd832)

## 0.11.5
2018-05-15

### Fixes

- update go-semrel to fix performance issue with lot of merges (95c91b25a4d0634f1fb8171ea4217e7bf81128d6)

## 0.11.4
2018-04-30

### Fixes

- increase gitlab api timeout (close #33) (a3d5546f7525873edee649790f420ed3cf430fcc)

## 0.11.3
2018-04-30

### Fixes

- update service desk email (f4caf99edeab0be573b4213fa4068383190cf7b1)

## 0.11.2
2018-04-21

### Fixes

- remove setter for SHA (e2d35528728b423b7186058aaa17060dfb963a49)
- don't create an additional pipeline on tag operation (df5091ebdc5af3b0109c5014f48237c071a17e03)

## 0.11.1
2018-04-14

### Fixes

- fix scope formatting #29 (0b89905fc0091541b0daea2f90bb1e3a41636851)
- add separate template for changelog entry #30 (364adf91a65adfe2061ae534c809794e68d27dbf)
- add download marker to release note #28 (24583c375d07e6a2bed9ef5d428b5e94b832e1ac)
- **site:** update status text (ac5866e2cf78c15b839ea80dcbf79f57f7887cc8)

## 0.11.0
2018-04-12

### Breaking changes

#### Don't search parent directories for git repo (194576ef7e8a187d1de7efc8e704f8c7a9c90f5e)

Commands must be run at root of git repository

#### prepare for using github.com/juranki/semrel (b37f29432a6b64830d6bb032f01417266b123206)

- gitlab specific code to gitlabutil,
- remove local semrel

#### use github/juranki/go-semrel to analyze repository (287bb446308e3d747b1ab08d9ba7e7201e261a75)

- analyzing the whole history for empty changelog is no longer supported
- link to compare has been removed (temporarily)

### Features

- detect and include refactorings (d828691aaacdb533d9335811b3b1debf03ffeac1)

### Fixes

- add-download (1f4e896d6edaec936943036dd9acaddb9cd6bbc6)
- commit (c6c43b65e90c4416dd9f687895a7b9e61e9e7c7e)
- *site:* add warning about the refactor (a7808b2127a5683655cf9cc4f5d6fda61d934c6f)

## 0.10.1

### Fixes

- **site:** make toc easier to read (c9f7d28c577b3febc929531d2ac2e47084bd9888)
- **site:** update status statement (98d5e505311b2de6d4afb03a8b587f4d1eb5ae4a)
- **site:** update feedback section (6ef104274dde9a7a48ae4532bfb26a8bb3d463f5)
- **site:** adjust line-heights (1c09546e9913194d96447d5ca6bafbb075f9cef8)

## 0.10.0

### Breaking changes

#### Remove `dry-run` flag (2b3bf4fad8281c56f9a756292d99aa3a5ddee41f)

`dry-run` flag is removed. It interweaves with the main functionality,
causing complexity that is not worth it at this stage, when the
functionality hasn't yet stabilized, and code needs refactoring.

Separate command(s) will be introduced to check configuration and
output results of change analysis.

### Features

- add `test-git` command #20 (7e5086447f86d9c4dc7415ea0a0d4ccb5a03c2ab)
- add `test-api` command (f53416e94e4adf430b71d92d6ea4bbacb3a751ef)

### Fixes

- Fix breaking change regexp #26 (65c4911147414cf70294b580a478fc6b643e58e9)

## 0.9.4

### Fixes

- **site:** change code font (684134a53188217681e1e376e38681fbb851e341)
- **site:** adjust sizes (92f3ed66f5be7e103a0aaf69ac99b1c1d1bf7cb3)
- improve site style (ec7256b079e43d6937a7fd88242b52d88476236a)

## 0.9.3

### Fixes

- allow capital letters in entry type #25 (56ea5b1950ca9430f396aa2b0f3d0e0bbcccb963)
- **site:** add changelog to sidebar (990ed4fdb178c373611027c93af1ed85180904e6)
- **site:** left align sidebar on mobile (71597e5cf576e476819912d7939daaa44c887d5d)
- **site:** fix issue email link on android (08bb83cfe69b285c0085fb2f3b76dd8fe1eda8f3)
- add descriptions to commands #19 (d692dcd005abd61a461db07292a8514a6f241473)

## 0.9.2

### Fixes

- add feedback section (fdff04b17db343988341c41b8c37d10e1732df4e)

## 0.9.1

### Fixes

- include lightweight tags too (193e7c2ca5f5c35c26f58042d365488515e3392a)
- handle merges more reliably (e60aa27d7e18d0fb6e31b5e2c0a61bf2461e052e)
- show changelog instead of tag page (f16043a0aa5d0a4735508424cdcb8f71ca067412)

## 0.9.0

### Features

- add `-p` to `next-version` (21519ed84d9cedfb227e5164a8fd9f4a4ff01e24)



## 0.8.0

### Breaking changes

#### remove unnecessary flag `note` (715d2321a46099ba34acc888143a9547a769e0a5)

`note` flag is removed from `tag` and `commit-and-tag`

### Fixes

- remove unnecessary empty lines from release note and changelog entry (c1e0205e96a85f116a30a84e4e8be6aa3f312bcf)
- more descriptive error message for `next-version` (fd8aeec33e840fb86bc990e568a0e7fce366af8d)

## 0.7.3





### Fixes

- **site:** add configuration to getting started section (75222e4e14c4707304d18f194d46ff87bcbf43fb)

## 0.7.2





### Fixes

- **site:** remove the extra 'v' from image tag (eb22a5c063609d41164e2ae3418bc72c9fd1d04d)

## 0.7.1





### Fixes

- **site:** add section for docker image (ebade2f5f5dcdac497eef81b2dc5c72e09bf0223)
- remove newline from short-version (50b11062f6418670e13cd552443696a4db1d3061)

## 0.7.0



### Features

- add parameter to skip ssl ca verification (fix #15) (f1cc98496e7c090dd69d82d3e05c19cc4ef5879b)



## 0.6.3





### Fixes

- update website intro (9663c0e47895cc95b777eda01ff3be99fb8d926c)

## 0.6.2





### Fixes

- improve help message for `release add-download` (33a9bd2f2651e4eaec2ca6fb270e06ba8f69b29d)

## 0.6.1





### Fixes

- add hidden command to get version correctly to documentation (268b882d041d543eab0cf888322abd9e2885c76a)

## 0.6.0



### Features

- release also as image (1df22978fe4cd554c81b23e2dcec042e2bd7b285)
- add a command to print the next version (95fb54cb6792e14a5557a6ea79570dfd76cc60ab)
- add param for gitlab api url #11 (ba953748a040aedc1b287370e5a1a81e3a7e2142)

### Fixes

- increase timeout to fix #13 (582c9afd84d68caf2ab4d6647c8e2add1ec4645f)
- **help:** reorder commands in help message (b5fecce6ec82716577ddeb4da8b814d3b065cb31)

## 0.5.0



### Features

- add downloads to release note (7e989bbb3f630b4c5feab10b5cf81edb539c5c76)
- implement commit-and-tag (fb26371862d9bf6a1e22b5ca016a5111264f31ba)
- add entry to changelog (fe6e1cad11a13a62b4bb81cee701ba47b14dafa9)
- create changelog (aeebf18764818af915482826bec05e7c3f424344)

### Fixes

- fix #10 (76fa8415623d4cfd68493fd51e3967b24c99c91f)
- remove todo from usage (9626df2e84edf7988f49fe9b38cfce9f095c6922)

## 0.4.0



### Features

- **cli:** add command line processing (44b444ea2233458a9961e6b730a0ffbf02df7975)

### Fixes

- fix #9 (31892025a41507d385c06a790252a732747ea1c8)
- update comment to match current implementation (b41bc77d672e0f6ceecf9696f817e62ff4558fc6)

## 0.3.0



### Features

- **release note:** add date and changes link #7 (c03dd896213d21093b49416b272e77f1b00a252e)
- **release note:** add scope to release note #7 (a7a515ac3fcf640be3fc876baa9b084cac243f0f)
- **release note:** add link to commit on feat and fix #7 (b516532b31205ac62895fe9000b3ca39a80d903c)
- **release note:** add link to commit on breaking change #7 (c647ceb583afabf451bb0217844e06c593111a0c)

### Fixes

- Fix #8 (8443609a6ca44e4854a0e0da34989bb03c812c21)

## 0.2.0

### Breaking changes

#### Remove `inspect-repo` (315679dd3311452fb2755e13a7cba2747155397b)

Now that `release` works, we can drop `inspect-repo`



### Fixes

- find previous version from commit log (a7be09ea4987e5376c5fb64ac48548841a921f0e)

## 0.1.0

### Breaking changes

#### rename packages (f94feeb34b5fba730aa29ece502c285396dac605)

inspect package renamed to repo

### Features

- create release (b35a245f85d9aa826f8423addb88b47ccf008cc6)
- validate gitlab token (497f84d4f6f809842ec5aaba208471e5a77bc100)
- check that required env vars are available (32bfe931d7cc6dc215716928c552f593264bcdf8)
- add inspect-repo executable (c712ea6fb0022c3e4e5743d6f787cc2f892d6bd7)
- format release note (9ecfd233fcb0e7884c23fa3947471569a89e8112)
- implement version bump (877e5dcce36a3280f235b7ff57a88a1d079a0a06)
- **inspect:** analyze commit message (36e42ab6ea36e23549ff7f86eb76cde78068f70b)

### Fixes

- Fix #6 (251c3bb40fade51d312e899fa86e6dcee9377a2f)
- fix tag analysis (b7ecb4216fd55fd7696c4612147402613a9e9497)
- **inspect:** only analyze commits after previous relese (ddc300cd720426ca80fb0664bbd5f1febc7d1608)
- use TagObjects to inspect tags (2196e1d0621514dced0166bf6b913dc4c281483d)