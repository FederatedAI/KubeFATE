# Contributing

## Welcome
KubeFATE is developed in the open, and is constantly being improved by our **users, contributors, and maintainers**. It is because of you that we can bring great software to the community.

This guide provides information on filing issues and guidelines for open source contributors. **Please leave comments / suggestions if you find something is missing or incorrect.**

When contributing to this repository, please first discuss the change you wish to make via issue, email, or any other method with the owners of this repository before making a change. 

Please note we have a code of conduct, please follow it in all your interactions with the project.

## Contribute Workflow
PR are always welcome, even if they only contain small fixes like typos or a few lines of codes. If there will be a significant effort, please:
1. document it as a proposal;
2. PR it to https://github.com/FederatedAI/KubeFATE/tree/master/proposals;
3. Get a discussion by creating [feature request issues](https://github.com/FederatedAI/KubeFATE/issues/new?assignees=&labels=&template=feature_request.md&title=) and refer to the proposal;
4. Implement the feature when the proposal got 2+ maintainer approved;
5. PR to `KubeFATE` **master** branch.

### Fork and clone
Fork the KubeFATE repository and clone the code to your local workspace. Per [Go's workspace instructions](https://golang.org/doc/code.html#Workspaces), place KubeFATE's code on your `GOPATH`.

Define a local working directory:
```
working_dir=./kubefate
```
Set user to match your github profile name:
```
user={your github profile name}
```

### Branch
Changes should be made on your own fork in a new branch. The branch should be named XXX-description where XXX is the number of the issue. PR should be rebased on top of master without multiple branches mixed into the PR. If your PR do not merge cleanly, use commands listed below to get it up to date.

```
#Suppose `kubefate` is the origin upstream

cd $working_dir
git fetch kubefate
git checkout master
git rebase kubefate/master
```
Branch from the updated master branch:
```
git checkout -b my_feature master
```

### Keep sync with upstream
Once your branch gets out of sync with the kubefate/master branch, use the following commands to update:
```
git checkout my_feature
git fetch -a
git rebase kubfate/master
Please use fetch / rebase (as shown above) instead of git pull. git pull does a merge, which leaves merge commits. These make the commit history messy and violate the principle that commits ought to be individually understandable and useful (see below). You can also consider changing your .git/config file via git config branch.autoSetupRebase always to change the behavior of git pull.
```

### Develop, Build and Test
Write code on the new branch in your fork. The coding style used in KubeFATE is suggested by the Golang community. See the [style doc](https://github.com/golang/go/wiki/CodeReviewComments) for details.

Try to limit column width to 120 characters for both code and markdown documents such as this one.

As we are enforcing standards set by [golint](https://github.com/golang/lint), please always run golint on source code before committing your changes. If it reports an issue, in general, the preferred action is to fix the code to comply with the linter's recommendation
because golint gives suggestions according to the stylistic conventions listed in [Effective Go](https://golang.org/doc/effective_go.html) and the [CodeReviewComments](https://github.com/golang/go/wiki/CodeReviewComments).
```
#Install fgt and golint

go get -u golang.org/x/lint/golint
go get github.com/GeertJohan/fgt

#In the #working_dir/harbor, run

go list ./... | grep -v -E 'vendor|tests' | xargs -L1 fgt golint

```

Unit test cases should be added to cover the new code. Unit test framework for backend services is using [go testing](https://golang.org/doc/code.html#Testing).

Run go test cases:
```
#cd #working_dir/src/[package]
go test -v ./...
```

In both the root folder of the `KubeFATE` project and the `k8s-deploy` project, there is a Makefile. The Makefile in the `KubeFATE` folder includes all commands to create a KubeFATE release that you found in a KubeFATE release, including:
* The docker-compose package;
* The charts;
* The KubeFATE CLI;
* The KubeFATE server image.

The Makefile in `k8s-deploy` includes:
* `go test` command to verify your contributions;
* Quick testing toolkit to deploy KubeFATE to a given Kubernetes;
* Generating new Swag documents if APIs changes
* and other subcommands.

Strongly suggest you use this Makefile as a toolkit for build and testing.

### Update the APIs and related documents
Our RESTful APIs are documented with [Swagger](https://swagger.io/)
If your commit that changes the RESTful APIs, make sure to run `make swag` in `./k8s-deploy/Makefile` to update the Swagger documents.
If your commit that exposes a user-faced function, make sure to add related introductions to documents.

### License
KubeFATE is applying Apache license, please include a short license header at the top of original source documents (code and documentation, but not the LICENSE and NOTICE files). An Apache license header example has list below.
```
/*
 * Copyright 2019-2022 VMware, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
```
_Note: "VMware, Inc." in the above example can be changed to the name of your organization or simplily "The KubeFATE Authors"._

### Commit
As KubeFATE has integrated the [DCO (Developer Certificate of Origin)](https://probot.github.io/apps/dco/) check tool, contributors are required to sign-off that they adhere to those requirements by adding a Signed-off-by line to the commit messages. Git has even provided a -s command line option to append that automatically to your commit messages, please use it when you commit your changes.
```
$ git commit -s -m 'This is my commit message'
```
Commit your changes if they're ready:
```
git add -A
git commit -s #-a
git push --force-with-lease $user my_feature
The commit message should follow the convention on How to Write a Git Commit Message. Be sure to include any related GitHub issue references in the commit message. 
```

### Push and Create PR
When ready for review, push your branch to your fork repository on github.com:
```
git push --force-with-lease $user my_feature
```
Then visit your fork at https://github.com/$user/KubeFATE and click the `Compare & Pull Request` button next to your `my_feature` branch to create a new pull request (PR). The PR should:
1. Ensure all unit test passed;
2. Ensure any install or build dependencies are removed before the end of the layer when doing a build;
3. The tittle of PR should hightlight what it solve briefly. The description of PR should refer to all the issues that it addresses. Ensure to put a reference to issues (such as `Close #xxx` and `Fixed #xxx`)  Please refer to the [PULL_REQUEST_TEMPLATE.md](https://github.com/FederatedAI/KubeFATE/blob/master/PULL_REQUEST_TEMPLATE.md).

Once your pull request has been opened it will be assigned to one or more reviewers. Those reviewers will do a thorough code review, looking for correctness, bugs, opportunities for improvement, documentation and comments, and style.

Commit changes made in response to review comments to the same branch on your fork.

### Automated Testing
Once your pull request has been opened, KubeFATE will run two CI pipelines against it.
1. In the Github action, your source code will be checked via golint, go vet and go race that makes sure the code is readable, safe and correct. Also all of unit tests will be triggered via go test against the pull request. 
	* If any failure in `github action checks`, you need to figure out whether it is introduced by your commits.
	* If the coverage dramatic decline, you need to commit unit test to coverage your code.
2. In the Jenkins CI, the E2E test will be triggered. The pipeline is about to build and install a FATE cluster and run the federated learning workload.

## No-code way contributing
To contribute to the KubeFATE project, there are many ways beyond coding. 

### Submitting an issue
It is a great way to contribute to KubeFATE by reporting an issue. Well-written and complete bug reports are always welcome! Please open an issue on GitHub and follow the template to fill in the required information.

Before opening any issue, please look up the existing issues to avoid submitting a duplication. If you find a match, you can "subscribe" to it to get notified of updates. If you have additional helpful information about the issue, please leave a comment.

When reporting an issue of bugs, always include:
* What deployment mode using? (docker-compose or Kubernetes) This is very important because it is a total difference we implemented between docker-compose and Kubernetes
* Versions. Includes the KubeFATE CLI version, KubeFATE server version, FATE version, etc. And also we need to know the version and type of OS and Kubernetes. 

Because the issues are open to the public, when submitting the log and configuration files, be sure to remove any sensitive information, e.g. user name, password, IP address, and company name. You can replace those parts with "REDACTED" or other strings like "****".

Be sure to include the steps to reproduce the problem if applicable. It can help us understand and fix your issue faster.

Another very appreciated way of contribution is submitting an issue of `feature request`. Make sure you are selecting the right template to draft it. Make sure you describing the problem to solve clearly. And it is also very important to stress the value of the feature requested, it is the important basis of our decision and work arrangement. 

### Documenting (includes Wikis)
Help us update, correct, add documents is greatly appreciated! It will help others better use KubeFATE and build a better community. Some documents to guide a beginner on how to use KubeFATE are also welcome and can be added to Wikis. Contact us if there are some documents or articles you can contribute.

### Advocate and educate

We are welcome any articles or blog posts to advocate the KubeFATE project or educate others on how to use KubeFATE for federated learning lifecycle management. If you can let us know the links to the articles or blog posts, we can gather them to the Wiki and help promote them to the community.